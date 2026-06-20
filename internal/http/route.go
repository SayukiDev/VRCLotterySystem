package http

import "github.com/gin-gonic/gin"

type Route struct {
	h *Handle
}

func NewRoute(h *Handle) *Route {
	return &Route{
		h: h,
	}
}

func (r *Route) InjectRoute(e *gin.Engine) error {
	api := e.Group("/api")
	api.GET("getSiteData", r.h.GetSiteData)
	api.GET("getForm", r.h.GetForm)
	api.POST("submitForm", r.h.SubmitForm)
	api.GET("getAllowList", r.h.GetAllowList)
	api.GET("isActive", r.h.IsActive)
	return nil
}
