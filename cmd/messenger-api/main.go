package main

import (
	"net/http"
	"time"

	"github.com/asciishell/avito-backend/internal/postgresqldb"
	"github.com/asciishell/avito-backend/pkg/environment"
	"github.com/asciishell/avito-backend/pkg/log"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type config struct {
	DB          postgresqldb.DBCredential
	HTTPAddress string
	HTTPTimeout time.Duration
	MaxRequests int
	PrintConfig bool
}

func loadConfig() config {
	cfg := config{}
	cfg.DB.URL = environment.GetStr("DB_URL", "")
	cfg.DB.Debug = environment.GetBool("DB_DEBUG", false)
	cfg.DB.Migrate = environment.GetBool("DB_MIGRATE", false)
	cfg.MaxRequests = environment.GetInt("MAX_REQUESTS", 100)
	cfg.HTTPAddress = environment.GetStr("ADDRESS", ":9000")
	cfg.HTTPTimeout = environment.GetDuration("HTTP_TIMEOUT", 5*time.Second)
	cfg.PrintConfig = environment.GetBool("PRINT_CONFIG", false)
	if cfg.PrintConfig {
		log.New().Infof("%+v", cfg)
	}
	return cfg
}
func main() {
	cfg := loadConfig()
	db, err := postgresqldb.NewPostgresStorage(cfg.DB)
	if err != nil {
		log.New().Fatalf("can't use database:%s", err)
	}
	defer func() {
		_ = db.DB.Close()
	}()
	logger := log.New()
	handler := NewHandler(logger, db)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Throttle(cfg.MaxRequests))
	r.Use(middleware.Timeout(cfg.HTTPTimeout))

	r.Route("/", func(r chi.Router) {
		r.Post("/users/add", handler.CreateUser)
		r.Route("/chats", func(r chi.Router) {
			r.Post("/get", handler.GetChats)
			r.Post("/add", handler.CreateChat)
		})
		r.Route("/messages", func(r chi.Router) {
			r.Post("/get", handler.GetMessages)
			r.Post("/add", handler.CreateMessage)
		})
	})
	if err := http.ListenAndServe(cfg.HTTPAddress, r); err != nil {
		logger.Fatalf("server error:%s", err)
	}
}
