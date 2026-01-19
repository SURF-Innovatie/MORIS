package api

import (
	"context"
	"net/http"
	"os"
	"time"

	_ "github.com/SURF-Innovatie/MORIS/api/swag-docs"
	coreauth "github.com/SURF-Innovatie/MORIS/internal/app/auth"
	"github.com/SURF-Innovatie/MORIS/internal/app/errorlog"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/cachewarmup"
	adapterhandler "github.com/SURF-Innovatie/MORIS/internal/handler/adapter"
	analyticshandler "github.com/SURF-Innovatie/MORIS/internal/handler/analytics"
	apikeyhandler "github.com/SURF-Innovatie/MORIS/internal/handler/apikey"
	authhandler "github.com/SURF-Innovatie/MORIS/internal/handler/auth"
	budgethandler "github.com/SURF-Innovatie/MORIS/internal/handler/budget"
	crossrefhandler "github.com/SURF-Innovatie/MORIS/internal/handler/crossref"
	eventHandler "github.com/SURF-Innovatie/MORIS/internal/handler/event"
	eventpolicyhandler "github.com/SURF-Innovatie/MORIS/internal/handler/eventpolicy"
	authmiddleware "github.com/SURF-Innovatie/MORIS/internal/handler/middleware"
	notificationhandler "github.com/SURF-Innovatie/MORIS/internal/handler/notification"
	odatahandler "github.com/SURF-Innovatie/MORIS/internal/handler/odata"
	orcidhandler "github.com/SURF-Innovatie/MORIS/internal/handler/orcid"
	organisationhandler "github.com/SURF-Innovatie/MORIS/internal/handler/organisation"
	personhandler "github.com/SURF-Innovatie/MORIS/internal/handler/person"
	producthandler "github.com/SURF-Innovatie/MORIS/internal/handler/product"
	projecthandler "github.com/SURF-Innovatie/MORIS/internal/handler/project"
	commandHandler "github.com/SURF-Innovatie/MORIS/internal/handler/project/command"
	systemhandler "github.com/SURF-Innovatie/MORIS/internal/handler/system"
	userhandler "github.com/SURF-Innovatie/MORIS/internal/handler/user"
	zenodohandler "github.com/SURF-Innovatie/MORIS/internal/handler/zenodo"
	"github.com/SURF-Innovatie/MORIS/internal/infra/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/do/v2"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// SetupRouter initializes the HTTP router with all handlers and middleware
func SetupRouter(injector do.Injector) *chi.Mux {
	// Get Services from DI Container
	authSvc := do.MustInvoke[coreauth.Service](injector)
	errorLogSvc := do.MustInvoke[errorlog.Service](injector)
	warmup := do.MustInvoke[cachewarmup.Service](injector)

	// Get HTTP Handlers from DI Container
	personHandler := do.MustInvoke[*personhandler.Handler](injector)
	userHandler := do.MustInvoke[*userhandler.Handler](injector)
	apiKeyHandler := do.MustInvoke[*apikeyhandler.Handler](injector)
	authHandler := do.MustInvoke[*authhandler.Handler](injector)
	orcidHandler := do.MustInvoke[*orcidhandler.Handler](injector)
	zenodoHandler := do.MustInvoke[*zenodohandler.Handler](injector)
	crossrefHandler := do.MustInvoke[*crossrefhandler.Handler](injector)
	organisationHandler := do.MustInvoke[*organisationhandler.Handler](injector)
	rbacHandler := do.MustInvoke[*organisationhandler.RBACHandler](injector)
	productHandler := do.MustInvoke[*producthandler.Handler](injector)
	notificationHandler := do.MustInvoke[*notificationhandler.Handler](injector)
	evtHandler := do.MustInvoke[*eventHandler.Handler](injector)
	eventPolicyHandler := do.MustInvoke[*eventpolicyhandler.Handler](injector)
	projHandler := do.MustInvoke[*projecthandler.Handler](injector)
	projCmdHandler := do.MustInvoke[*commandHandler.Handler](injector)
	systemHandler := do.MustInvoke[*systemhandler.Handler](injector)
	adapterHandler := do.MustInvoke[*adapterhandler.Handler](injector)
	budgetHandler := do.MustInvoke[*budgethandler.Handler](injector)
	analyticsHandler := do.MustInvoke[*analyticshandler.Handler](injector)
	oDataHandler := do.MustInvoke[*odatahandler.Handler](injector)

	// Setup Router
	r := chi.NewRouter()

	// Default zerolog logger
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Basic CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   env.Global.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route(env.Global.APIBasePath, func(r chi.Router) {
		r.Use(authmiddleware.ErrorLoggingMiddleware(errorLogSvc))
		authhandler.MountRoutes(r, authSvc, authHandler)
		systemhandler.MountRoutes(r, systemHandler)

		r.Group(func(r chi.Router) {
			r.Use(authmiddleware.AuthMiddleware(authSvc))
			r.Route("/projects", func(r chi.Router) {
				projecthandler.MountProjectRoutes(r, projHandler)
				r.Get("/{id}/roles", projHandler.ListAvailableRoles)
				commandHandler.MountProjectCommandRouter(r, projCmdHandler)
				// Event policy routes for projects
				eventPolicyHandler.RegisterProjectRoutes(r)

				// Budget routes for projects
				budgetHandler.RegisterProjectRoutes(r)
			})
			organisationhandler.MountOrganisationRoutes(r, organisationHandler, rbacHandler)
			eventHandler.MountEventRoutes(r, evtHandler)
			personhandler.MountPersonRoutes(r, personHandler)
			orcidhandler.MountRoutes(r, orcidHandler)
			zenodohandler.MountRoutes(r, zenodoHandler)
			producthandler.MountProductRoutes(r, productHandler)
			notificationhandler.MountNotificationRoutes(r, notificationHandler)
			userhandler.MountUserRoutes(r, userHandler)
			apiKeyHandler.RegisterRoutes(r)
			crossrefhandler.MountCrossrefRoutes(r, crossrefHandler)
			adapterhandler.MountRoutes(r, adapterHandler)

			// Event Policies routes (standalone and org-scoped)
			eventPolicyHandler.RegisterRoutes(r)
			r.Route("/organisations/{id}/policies", func(r chi.Router) {
				r.Get("/", eventPolicyHandler.ListForOrgNode)
				r.Post("/", eventPolicyHandler.CreateForOrgNode)
			})

			budgetHandler.RegisterRoutes(r)
			analyticsHandler.RegisterRoutes(r)
			oDataHandler.RegisterRoutes(r)
		})
	})

	// Warmup cache in background
	go func() {
		cached, err := warmup.WarmupProjects(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("Failed to warmup cache")
		} else {
			log.Info().Msgf("Warmed up cache for %d projects", cached)
		}
	}()

	// Serve the generated swagger JSON and assets and the Swagger UI at /swagger/
	r.Get("/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "api/swag-docs/swagger.json")
	})
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("swagger.json"), // the url pointing to API definition
	))

	return r
}
