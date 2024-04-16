package realstatehdlr

import (
	"github.com/gin-gonic/gin"
	"github.com/natanchagas/gin-crud/internal/core/domain"
	"github.com/natanchagas/gin-crud/internal/core/ports"
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
		c.JSON(400, err.Error())
		return
	}

	realState, err = h.RealStateService.Create(ctx, realState)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(201, realState)
	return
}

func (h *RealStateHandler) BuildRoutes(router *gin.Engine) {
	realState := router.Group("/realstate/")

	realState.POST("/", h.create)
}
