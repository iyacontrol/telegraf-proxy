package api

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/iyacontrol/telegraf-proxy/aggregate"
	"github.com/iyacontrol/telegraf-proxy/api/handlers"
	"github.com/iyacontrol/telegraf-proxy/config"
	"github.com/iyacontrol/telegraf-proxy/discovery"
)

func InitAPI(reg *discovery.Center) {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome telegraf-proxy"))
	})

	aggregator := &aggregate.Aggregator{HTTP: &http.Client{Timeout: time.Duration(config.Cfg.Aggregate.Timeout) * time.Millisecond}}
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		aggregator.Aggregate(reg, w)
	})

	r.Post("/reload", handlers.Reload)

	addr := net.JoinHostPort("", config.Cfg.HTTP.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("Httpserver: ListenAndServe() error: %s", err)
		}
	}()
}
