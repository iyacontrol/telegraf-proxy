package api

import (
	"net"

	"github.com/iyacontrol/telegraf-proxy/config"
	"github.com/iyacontrol/telegraf-proxy/discovery"
	baa "gopkg.in/baa.v1"
)

func InitAPI(stop chan struct{}, reg *discovery.Registry) {

	go func() {
		app := baa.New()
		// init middleware
		initMiddleware(app)
		// init router
		initRouter(app)
		// register router
		register(app)

		addr := net.JoinHostPort(config.Cfg.HTTP.Address, config.Cfg.HTTP.Port)
		app.Run(addr)
	}()
}
