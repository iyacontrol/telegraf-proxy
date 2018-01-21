package api

import (
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/iyacontrol/telegraf-proxy/config"
	"github.com/iyacontrol/telegraf-proxy/discovery"
)

func InitApi(reg *discovery.Register) {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome telegraf-proxy"))
	})

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
