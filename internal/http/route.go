package http

import (
	"github.com/SayukiDev/VRCLotterySystem/internal/http/handle"
	"github.com/SayukiDev/VRCLotterySystem/internal/http/middleware"
	"github.com/gin-gonic/gin"
)

type Route struct {
	h *handle.Handle
}

func NewRoute(h *handle.Handle) *Route {
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

func (r *Route) InjectAuthedRoute(e *gin.Engine, token string) error {
	authed := e.Group("/api/authed")
	authed.Use(middleware.TokenAuth(token))

	// Blacklist
	authed.GET("getBlackList", r.h.GetBlackList)
	authed.POST("addBlackList", r.h.AddBlackList)
	authed.POST("deleteBlackList", r.h.DeleteBlackList)

	// Staff list
	authed.GET("getStaffList", r.h.GetStaffList)
	authed.POST("addStaff", r.h.AddStaff)
	authed.POST("deleteStaff", r.h.DeleteStaff)

	// Drawing
	authed.POST("setDrawing", r.h.SetDrawing)
	authed.GET("getDrawing", r.h.GetDrawing)
	authed.POST("drawing", r.h.Drawings)

	// Results
	authed.GET("getResults", r.h.GetResults)
	authed.POST("removeResults", r.h.RemoveResults)
	return nil
}
