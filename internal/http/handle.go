package http

import (
	"fmt"
	"strings"

	"github.com/SayukiDev/VRCLotterySystem/internal/data"
	"github.com/SayukiDev/VRCLotterySystem/internal/provider"

	"github.com/gin-gonic/gin"
)

type Handle struct {
	p *provider.Provider
}

func NewHandle(p *provider.Provider) *Handle {
	return &Handle{p: p}
}

func (h *Handle) GetTerms(c *gin.Context) {
	c.JSON(200, CommonResp{
		Code: 200,
		Msg:  "success",
		Data: h.p.C.Terms,
	})
}

func (h *Handle) GetForm(c *gin.Context) {
	c.JSON(200, CommonResp{
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
		c.JSON(400, CommonResp{
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
		case "inputs":
			_, ok = req.Input.Content[i]
		case "options":
			_, ok = req.Input.Selected[i]
		}
		if !ok {
			c.JSON(400, CommonResp{
				Code: 400,
				Msg:  "bad request",
				Data: fmt.Sprintf("required field(%s) not filled", v.Title),
			})
			return
		}
	}
	err = h.p.AddInput(req.Input)
	if err != nil {
		c.JSON(500, CommonResp{
			Code: 500,
			Msg:  "internal server error",
			Data: err.Error(),
		})
		c.Error(err)
		return
	}
	c.JSON(200, CommonResp{
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
