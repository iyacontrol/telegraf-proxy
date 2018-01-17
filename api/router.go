package api

import (
	"github.com/go-baa/router/regtree"
	baa "gopkg.in/baa.v1"
)

func initRouter(b *baa.Baa) {
	// router
	b.SetDI("router", regtree.New(b))
}

func register(b *baa.Baa) {
	b.SetAutoHead(true)
	b.SetAutoTrailingSlash(true)

	b.Get("/", func(c *baa.Context) {
		c.JSON(200, map[string]interface{}{"code": 0, "message": ""})
	})

	// 应用状态监测
	b.Get("/status", func(c *baa.Context) {
		c.Text(200, []byte("success"))
	})
}
