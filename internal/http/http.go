package http

import (
	"github.com/SayukiDev/VRCLotterySystem/config"
	"github.com/SayukiDev/VRCLotterySystem/internal/provider"

	"github.com/gin-gonic/gin"
)

type Http struct {
	e     *gin.Engine
	route *Route
}

func NewHttp(c *config.Config, p *provider.Provider) *Http {
	e := gin.New()
	e.Use(Logger())
	e.Use(gin.Recovery())
	return &Http{
		e:     e,
		route: NewRoute(NewHandle(p)),
	}
}

func (h *Http) Start(addr string) error {
	err := h.route.InjectRoute(h.e)
	if err != nil {
		return err
	}
	return h.e.Run(addr)
}
