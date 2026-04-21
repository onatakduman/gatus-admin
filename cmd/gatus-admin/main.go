package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"

	"github.com/onatakduman/gatus-admin/internal/auth"
	"github.com/onatakduman/gatus-admin/internal/brand"
	"github.com/onatakduman/gatus-admin/internal/configio"
	"github.com/onatakduman/gatus-admin/internal/handlers"
	"github.com/onatakduman/gatus-admin/internal/web"
)

func main() {
	migrate := flag.Bool("migrate", false, "Migrate single config.yaml into multi-file config/ directory")
	flag.Parse()

	configDir := envOrDefault("CONFIG_DIR", "/config")
	statusURL := envOrDefault("STATUS_URL", "/")

	if *migrate {
		if err := configio.Migrate(configDir); err != nil {
			log.Fatalf("migration failed: %v", err)
		}
		log.Println("migration complete")
		return
	}

	passwordHash := os.Getenv("ADMIN_PASSWORD_HASH")
	if passwordHash == "" {
		log.Fatal("ADMIN_PASSWORD_HASH is required")
	}
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY is required")
	}
	gatusContainer := envOrDefault("GATUS_CONTAINER", "gatus")
	dev := os.Getenv("DEV") == "1"

	store := auth.NewStore([]byte(secretKey), !dev)
	tpl := web.NewTemplates(dev)

	h := &handlers.Handlers{
		ConfigDir:      configDir,
		GatusContainer: gatusContainer,
		Tpl:            tpl,
		Store:          store,
		PasswordHash:   passwordHash,
		Dev:            dev,
		StatusURL:      statusURL,
		BrandStore:     brand.Store{RootDir: configDir},
	}

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	csrfOpts := []csrf.Option{csrf.Path("/admin")}
	if dev {
		csrfOpts = append(csrfOpts, csrf.Secure(false))
	}
	csrfMw := csrf.Protect([]byte(secretKey), csrfOpts...)

	r.Route("/admin", func(r chi.Router) {
		r.Use(csrfMw)
		r.Get("/static/*", web.StaticHandler())

		r.Get("/login", h.LoginGet)
		r.Post("/login", h.LoginPost)
		r.Post("/logout", h.Logout)

		// Public: logos on the login page need to load before auth.
		r.Get("/assets/{name}", h.AssetGet)

		r.Group(func(r chi.Router) {
			r.Use(h.RequireAuth)

			r.Get("/", h.EndpointsList)

			r.Get("/endpoints", h.EndpointsList)
			r.Post("/endpoints", h.EndpointCreate)
			r.Get("/endpoints/{name}/edit", h.EndpointEdit)
			r.Post("/endpoints/{name}", h.EndpointUpdate)
			r.Post("/endpoints/{name}/delete", h.EndpointDelete)
			r.Get("/endpoints/bulk", h.EndpointsBulkGet)
			r.Post("/endpoints/bulk", h.EndpointsBulkPost)

			r.Get("/alerting", h.AlertingGet)
			r.Post("/alerting", h.AlertingPost)

			r.Get("/ui", h.UIGet)
			r.Post("/ui", h.UIPost)

			r.Get("/announcements", h.AnnouncementsList)
			r.Post("/announcements", h.AnnouncementCreate)
			r.Post("/announcements/{idx}/delete", h.AnnouncementDelete)

			r.Get("/maintenance", h.MaintenanceGet)
			r.Post("/maintenance", h.MaintenancePost)

			r.Get("/security", h.SecurityGet)
			r.Post("/security", h.SecurityPost)

			r.Get("/tunneling", h.TunnelingGet)
			r.Post("/tunneling", h.TunnelingPost)

			r.Get("/external-endpoints", h.ExternalList)
			r.Post("/external-endpoints", h.ExternalCreate)
			r.Post("/external-endpoints/{name}/delete", h.ExternalDelete)

			r.Get("/suites", h.SuitesList)
			r.Post("/suites", h.SuiteCreate)
			r.Post("/suites/{name}/delete", h.SuiteDelete)

			r.Get("/core", h.CoreGet)
			r.Post("/core", h.CorePost)

			r.Get("/advanced", h.AdvancedGet)
			r.Get("/advanced/{file}", h.AdvancedFileGet)
			r.Post("/advanced/{file}", h.AdvancedFilePost)

			r.Get("/brand", h.BrandGet)
			r.Post("/brand", h.BrandPost)
			r.Post("/brand/upload/{slot}", h.BrandUpload)
			r.Post("/brand/reset/{slot}", h.BrandReset)

			r.Post("/apply", h.Apply)
		})
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/admin", http.StatusFound)
	})

	addr := ":8000"
	log.Printf("admin listening on %s (dev=%v configDir=%s)", addr, dev, configDir)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
