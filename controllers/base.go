package controllers

import (
	"fmt"
	"sync"

	"github.com/yairp7/go-common-lib/logger"
	"github.com/yairp7/papi/common"
)

type Controller interface {
	common.Closer
}

type BaseController struct {
	Name       string
	currentOps sync.WaitGroup
	loggerImpl logger.Logger
	services   []common.Closer
}

func NewBaseController(name string, loggerImpl logger.Logger) BaseController {
	return BaseController{
		Name:       name,
		currentOps: sync.WaitGroup{},
		loggerImpl: loggerImpl,
	}
}

func (c *BaseController) Close() {
	c.loggerImpl.Debug(fmt.Sprintf("%s Shutdown\n", c.Name))
	c.currentOps.Wait()
	for _, closeableService := range c.services {
		closeableService.Close()
	}
}

func (c *BaseController) RegisterOp() {
	c.currentOps.Add(1)
}

func (c *BaseController) UnregisterOp() {
	c.currentOps.Done()
}

func (c *BaseController) RegisterService(services ...common.Closer) {
	for _, service := range services {
		c.services = append(c.services, service)
	}
}
