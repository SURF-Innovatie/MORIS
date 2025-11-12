package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

//go:embed all:../frontend/dist
var frontendFS embed.FS

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost" // Points to host's localhost where docker-compose mapped the port
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "user"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "password"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "your_app_db"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Example Ent client connection (uncomment and adjust when you set up Ent)
	// client, err := ent.Open("postgres", dsn)
	// if err != nil {
	// 	logrus.Fatalf("failed opening connection to postgres: %v", err)
	// }
	// defer client.Close()

	// // Run migration (optional, but good for dev to ensure schema is up-to-date)
	// if err := client.Schema.Create(context.Background()); err != nil {
	// 	logrus.Fatalf("failed creating schema resources: %v", err)
	// }
	logrus.Infof("Attempting to connect to database at %s:%s/%s", dbHost, dbPort, dbName)
	// Add actual database ping/health check here if not using Ent for schema creation
	// For example, using sql.Open and db.Ping()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// --- API Routes ---
	r.Route("/api", func(r chi.Router) {
		// Example API endpoint
		r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello from the Go API!"))
		})
		// Your Ent-generated API handlers would go here, using the ent client
		// e.g., r.Mount("/users", usersHandler(client))
	})

	// --- Serve Frontend Static Files (only for production build) ---
	// In development, Vite dev server handles frontend, so 'dist' won't exist.
	// This makes the Go app still runnable in dev for API-only testing.
	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		logrus.Printf("Frontend 'dist' directory not found at path 'frontend/dist'. This is expected in development when Vite serves the frontend separately. Serving API only. Error: %v", err)
	} else {
		fileServer := http.FileServer(http.FS(distFS))

		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			path := strings.TrimPrefix(r.URL.Path, "/")
			if path == "" {
				path = "index.html"
			}

			// Try to open the requested file
			_, err := distFS.Open(path)
			if err == nil {
				// If file exists, serve it
				fileServer.ServeHTTP(w, r)
				return
			}
			// If file not found (e.g., client-side route), serve index.html for SPA routing
			http.ServeFileFS(w, r, distFS, "index.html")
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logrus.Infof("Go Backend Server starting on http://localhost:%s", port)
	logrus.Fatal(http.ListenAndServe(":"+port, r))
}
