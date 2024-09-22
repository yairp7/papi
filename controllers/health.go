package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yairp7/go-common-lib/logger"
)

type HealthController struct {
	BaseController
}

func NewHealthController(loggerImpl logger.Logger) *HealthController {
	return &HealthController{
		BaseController: NewBaseController("HealthController", loggerImpl),
	}
}

func (c *HealthController) Status(ctx *gin.Context) {
	ctx.String(http.StatusOK, "OK!")
}
