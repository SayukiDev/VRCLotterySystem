package handle

import (
	"fmt"
	"strings"
	"time"

	"github.com/SayukiDev/VRCLotterySystem/internal/data"
	"github.com/SayukiDev/VRCLotterySystem/internal/http/common"
	"github.com/gin-gonic/gin"
)

func (h *Handle) GetSiteData(c *gin.Context) {
	c.JSON(200, common.CommonResp{
		Code: 200,
		Msg:  "success",
		Data: h.p.C.SiteData,
	})
}

type GetFormReq struct {
	Id string `form:"id" binding:"required,max=10"`
}

func (h *Handle) GetForm(c *gin.Context) {
	req := &GetFormReq{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		c.JSON(400, common.CommonResp{
			Code: 400,
			Msg:  "bad request",
			Data: err.Error(),
		})
		return
	}
	ok := false
	h.p.Data.Read(func(d *data.Content) {
		if d.Id == req.Id {
			ok = true
		}
	})
	if !ok {
		c.JSON(404, common.CommonResp{
			Code: 404,
			Msg:  "not found",
		})
		return
	}
	c.JSON(200, common.CommonResp{
		Code: 200,
		Msg:  "success",
		Data: h.p.C.Form,
	})
}

type SubmitFormReq struct {
	Input data.Input `json:"inputs" binding:"required"`
}

func (h *Handle) SubmitForm(c *gin.Context) {
	req := &SubmitFormReq{}
	err := c.ShouldBindJSON(req)
	if err != nil {
		c.JSON(400, common.CommonResp{
			Code: 400,
			Msg:  "bad request",
			Data: err.Error(),
		})
		return
	}
	// Check required fields
	for i, v := range h.p.C.Form {
		if !v.Required {
			continue
		}
		var ok bool
		switch v.Type {
		case "input":
			_, ok = req.Input.Content[i]
		case "options":
			_, ok = req.Input.Selected[i]
		}
		if !ok {
			c.JSON(400, common.CommonResp{
				Code: 400,
				Msg:  "bad request",
				Data: fmt.Sprintf("required field(%s) not filled", v.Title),
			})
			return
		}
	}
	err = h.p.AddInput(req.Input)
	if err != nil {
		c.JSON(500, common.CommonResp{
			Code: 500,
			Msg:  "internal server error",
			Data: err.Error(),
		})
		c.Error(err)
		return
	}
	c.JSON(200, common.CommonResp{
		Code: 200,
		Msg:  "success",
	})
}

func (h *Handle) GetAllowList(c *gin.Context) {
	r := h.p.GetResults()
	r = append(r, h.p.GetStaffList()...)
	resp := strings.Join(r, "\n")
	c.String(200, resp)
}

func (h *Handle) IsActive(c *gin.Context) {
	ok := false
	h.p.Data.Read(func(d *data.Content) {
		if d.Date.After(time.Now()) {
			ok = true
		}
	})
	c.JSON(200, common.CommonResp{
		Code: 200,
		Msg:  "success",
		Data: ok,
	})
}
