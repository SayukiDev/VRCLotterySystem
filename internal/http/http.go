package http

import (
	"github.com/SayukiDev/VRCLotterySystem/config"
	"github.com/SayukiDev/VRCLotterySystem/internal/http/handle"
	"github.com/SayukiDev/VRCLotterySystem/internal/http/middleware"
	"github.com/SayukiDev/VRCLotterySystem/internal/provider"

	"github.com/gin-gonic/gin"
)

const maxBodyBytes = 64 * 1024

type Http struct {
	e     *gin.Engine
	route *Route
	p     *provider.Provider
}

func NewHttp(c *config.Config, p *provider.Provider) *Http {
	e := gin.New()
	e.Use(middleware.Logger())
	e.Use(gin.Recovery())
	e.Use(middleware.BodyLimit(maxBodyBytes))
	return &Http{
		e:     e,
		route: NewRoute(handle.NewHandle(p)),
		p:     p,
	}
}

func (h *Http) Start(addr string) error {
	err := h.route.InjectRoute(h.e)
	if err != nil {
		return err
	}
	err = h.route.InjectAuthedRoute(h.e, h.p.C.Token)
	if err != nil {
		return err
	}
	return h.e.Run(addr)
}
