package handle

import (
	"time"

	"github.com/SayukiDev/VRCLotterySystem/internal/http/common"
	"github.com/gin-gonic/gin"
)

func (h *Handle) GetBlackList(c *gin.Context) {
	bl := h.p.GetBlacklist()
	c.JSON(200, common.CommonResp{
		Code: 200,
		Msg:  "success",
		Data: bl,
	})
}

type AddBlackListReq struct {
	Id string `form:"id" binding:"required,max=10"`
}

func (h *Handle) AddBlackList(c *gin.Context) {
	req := &AddBlackListReq{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		c.JSON(400, common.CommonResp{
			Code: 400,
			Msg:  "bad request",
			Data: err.Error(),
		})
		return
	}
	err = h.p.AddToBlackList(req.Id)
	if err != nil {
		c.JSON(500, common.CommonResp{
			Code: 500,
			Msg:  "internal server error",
		})
		c.Error(err)
	}
}

func (h *Handle) DeleteBlackList(c *gin.Context) {
	req := &AddBlackListReq{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		c.JSON(400, common.CommonResp{
			Code: 400,
			Msg:  "bad request",
			Data: err.Error(),
		})
		return
	}
	err = h.p.RemoveFromBlackList(req.Id)
	if err != nil {
		c.JSON(500, common.CommonResp{
			Code: 500,
			Msg:  "internal server error",
		})
		c.Error(err)
		return
	}
	c.JSON(200, common.CommonResp{
		Code: 200,
		Msg:  "success",
	})
}

func (h *Handle) GetStaffList(c *gin.Context) {
	c.JSON(200, common.CommonResp{
		Code: 200,
		Msg:  "success",
		Data: h.p.GetStaffList(),
	})
}

type AddStaffReq struct {
	Id string `form:"id" binding:"required,max=10"`
}

func (h *Handle) AddStaff(c *gin.Context) {
	req := &AddStaffReq{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		c.JSON(400, common.CommonResp{
			Code: 400,
			Msg:  "bad request",
			Data: err.Error(),
		})
		c.Error(err)
		return
	}
}

func (h *Handle) DeleteStaff(c *gin.Context) {
	req := &AddStaffReq{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		c.JSON(400, common.CommonResp{
			Code: 400,
			Msg:  "bad request",
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

type SetDrawingReq struct {
	Date time.Time `form:"date" binding:"required"`
	max  int       `form:"max" binding:"required,min=1"`
}

func (h *Handle) SetDrawing(c *gin.Context) {
	req := &SetDrawingReq{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		c.JSON(400, common.CommonResp{
			Code: 400,
			Msg:  "bad request",
			Data: err.Error(),
		})
		return
	}
	err = h.p.SetDrawing(req.max, req.Date)
	if err != nil {
		c.JSON(500, common.CommonResp{
			Code: 500,
			Msg:  "internal server error",
		})
		c.Error(err)
		return
	}
	c.JSON(200, common.CommonResp{
		Code: 200,
		Msg:  "success",
	})
}

func (h *Handle) Drawings(c *gin.Context) {
	us, err := h.p.Drawing()
	if err != nil {
		c.JSON(500, common.CommonResp{
			Code: 500,
			Msg:  "internal server error",
		})
		c.Error(err)
		return
	}
	err = h.p.AddResults(us)
	if err != nil {
		c.JSON(500, common.CommonResp{
			Code: 500,
			Msg:  "internal server error",
		})
		c.Error(err)
		return
	}
	c.JSON(200, common.CommonResp{
		Code: 200,
		Msg:  "success",
		Data: us,
	})
}

func (h *Handle) GetResults(c *gin.Context) {
	c.JSON(200, common.CommonResp{
		Code: 200,
		Msg:  "success",
		Data: h.p.GetResults(),
	})
}

type RemoveResultsReq struct {
	Id string `form:"id" binding:"required,max=10"`
}

func (h *Handle) RemoveResults(c *gin.Context) {
	req := &RemoveResultsReq{}
	err := c.ShouldBindQuery(req)
	if err != nil {
		c.JSON(400, common.CommonResp{
			Code: 400,
			Msg:  "bad request",
			Data: err.Error(),
		})
		return
	}
	err = h.p.DeleteResults(req.Id)
	if err != nil {
		c.JSON(500, common.CommonResp{
			Code: 500,
			Msg:  "internal server error",
		})
		c.Error(err)
		return
	}
	c.JSON(200, common.CommonResp{
		Code: 200,
		Msg:  "success",
	})
}
