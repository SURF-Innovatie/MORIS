package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/domain/commands"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/projection"
	"github.com/SURF-Innovatie/MORIS/internal/platform/eventstore"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var embeddedFrontend embed.FS

func distFS() (fs.FS, bool) {
	// 1) Prefer on-disk dist during dev
	wd, _ := os.Getwd()
	disk := filepath.Join(wd, "../frontend/dist")
	if st, err := os.Stat(disk); err == nil && st.IsDir() {
		return os.DirFS(disk), true
	}
	// 2) Fall back to embedded files for prod builds
	sub, err := fs.Sub(embeddedFrontend, "../frontend/dist")
	if err == nil {
		return sub, false
	}
	return nil, false
}

func dbDSN() string {
	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	user := getenv("DB_USER", "moris")
	pass := getenv("DB_PASSWORD", "moris")
	name := getenv("DB_NAME", "moris")
	// lib/pq style DSN
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name)
}

func openAndPingPostgres(ctx context.Context) (*sql.DB, error) {
	dsn := dbDSN()
	logrus.Infof("Connecting to Postgres: %s/%s@%s:%s",
		getenv("DB_NAME", "moris"),
		getenv("DB_USER", "moris"),
		getenv("DB_HOST", "localhost"),
		getenv("DB_PORT", "5432"),
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	// Short ping with timeout
	pctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := db.PingContext(pctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	// DB
	client, err := ent.Open("postgres", dbDSN())
	if err != nil {
		logrus.Fatalf("ent open: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Run automatic schema migration (creates tables)
	if err := client.Schema.Create(ctx); err != nil {
		logrus.Fatalf("ent migrate: %v", err)
	}

	// Create event store
	store := eventstore.NewEntStore(client)
	// Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API
	r.Route("/api", func(r chi.Router) {
		r.Get("/hello", func(w http.ResponseWriter, _ *http.Request) {
			_, err := w.Write([]byte("Hello from the Go API"))
			if err != nil {
				logrus.Fatal(err)
			}
		})
		r.Get("/demo", func(w http.ResponseWriter, _ *http.Request) {
			id := uuid.New()

			start := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)
			end := time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC)
			org := entities.Organisation{Name: "Org"}
			people := []entities.Person{{Name: "Ada"}}

			ev, err := commands.StartProject(id, "Alpha", "First", start, end, people, org)
			if err != nil {
				logrus.Fatal(err)
			}

			if err := store.Append(context.Background(), id, 0, []events.Event{ev}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			loaded, ver, err := store.Load(context.Background(), id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			cur := projection.Reduce(id, loaded)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(fmt.Sprintf(`{"id":"%s", "ver":%d,"title":"%s"}`, id, ver, cur.Title)))
		})
	})

	// Frontend
	if fsys, fromDisk := distFS(); fsys != nil {
		fileServer := http.FileServer(http.FS(fsys))
		if fromDisk {
			logrus.Info("Serving frontend from disk: ../frontend/dist")
		} else {
			logrus.Info("Serving frontend from embedded files")
		}
		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimPrefix(r.URL.Path, "/")
			if path == "" {
				path = "index.html"
			}
			// Try file; else serve index.html for SPA routes
			if _, err := fs.Stat(fsys, path); err == nil {
				fileServer.ServeHTTP(w, r)
				return
			}
			http.ServeFileFS(w, r, fsys, "index.html")
		})
	} else {
		logrus.Warn("No frontend dist found or embedded. Only API will be served.")
	}

	port := getenv("PORT", "8080")
	logrus.Infof("Server listening on http://localhost:%s", port)
	logrus.Fatal(http.ListenAndServe(":"+port, r))
}
