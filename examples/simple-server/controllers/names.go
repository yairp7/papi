package controllers

import (
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yairp7/go-common-lib/logger"
	"github.com/yairp7/papi/controllers"
)

type NamesController struct {
	controllers.BaseController
	names []string
}

func NewNamesController(
	loggerImpl logger.Logger,
	names []string,
) *NamesController {
	return &NamesController{
		BaseController: controllers.NewBaseController("NamesController", loggerImpl),
		names:          names,
	}
}

func (c *NamesController) Name(ctx *gin.Context) {
	randName := c.names[rand.Intn(len(c.names))]
	ctx.String(http.StatusOK, randName)
}
