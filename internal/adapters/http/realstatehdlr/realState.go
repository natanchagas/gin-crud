package realstatehdlr

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/natanchagas/gin-crud/internal/core/domain"
	"github.com/natanchagas/gin-crud/internal/core/ports"
	"github.com/natanchagas/gin-crud/internal/pkg/customerrors"
)

type RealStateHandler struct {
	RealStateService ports.RealStateService
}

func NewRealStateHandler(service ports.RealStateService) *RealStateHandler {
	return &RealStateHandler{
		RealStateService: service,
	}
}

func (h *RealStateHandler) create(c *gin.Context) {
	ctx := c.Request.Context()

	var realState domain.RealState

	err := c.BindJSON(&realState)
	if err != nil {
		c.JSON(400, customerrors.BadRequest)
		return
	}

	realState, err = h.RealStateService.Create(ctx, realState)
	if err != nil {
		if cerr, ok := err.(customerrors.Error); ok {
			c.JSON(cerr.StatusCode, cerr)
			return
		}

		c.JSON(500, customerrors.Unexpected)
		return
	}

	c.JSON(201, realState)
	return
}

func (h *RealStateHandler) get(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	rid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, customerrors.BadRequest)
		return
	}

	realstate, err := h.RealStateService.Get(ctx, rid)
	if err != nil {
		if cerr, ok := err.(customerrors.Error); ok {
			c.JSON(cerr.StatusCode, cerr)
			return
		}

		c.JSON(500, customerrors.Unexpected)
		return
	}

	c.JSON(200, realstate)
	return
}

func (h *RealStateHandler) update(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	rid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, customerrors.BadRequest)
		return
	}

	var realState domain.RealState

	err = c.BindJSON(&realState)
	if err != nil {
		c.JSON(400, customerrors.BadRequest)
		return
	}

	realState, err = h.RealStateService.Update(ctx, realState, rid)
	if err != nil {
		if cerr, ok := err.(customerrors.Error); ok {
			c.JSON(cerr.StatusCode, cerr)
			return
		}

		c.JSON(500, customerrors.Unexpected)
		return
	}

	c.JSON(200, realState)
	return

}

func (h *RealStateHandler) delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	rid, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(400, customerrors.BadRequest)
		return
	}

	err = h.RealStateService.Delete(ctx, rid)
	if err != nil {
		if cerr, ok := err.(customerrors.Error); ok {
			c.JSON(cerr.StatusCode, cerr)
			return
		}

		c.JSON(500, customerrors.Unexpected)
		return
	}

	c.JSON(204, nil)
	return
}

func (h *RealStateHandler) BuildRoutes(router *gin.Engine) {
	realState := router.Group("/realstate/")

	realState.POST("/", h.create)
	realState.GET("/:id", h.get)
	realState.PUT("/:id", h.update)
	realState.DELETE("/:id", h.delete)
}
