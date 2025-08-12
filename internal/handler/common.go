package handler

import (
	"fmt"
	"main/docs"
	"main/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RespMsgError struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"error text"`
}

type RespMsgSuccess struct {
	Success bool `json:"success" example:"true"`
	Message any  `json:"message"`
}

func initSwagger() {
	cfg := config.GetConfig()
	docs.SwaggerInfo.Title = "Subscription API server"
	docs.SwaggerInfo.Description = "This is a sample CRUDL subscription server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = cfg.Listen.Addr
	docs.SwaggerInfo.BasePath = "/"
}

func (h *Handler) sendError(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, RespMsgError{
		Success: false,
		Message: fmt.Sprint(err),
	})
}

func (h *Handler) sendSuccess(c *gin.Context, code int, msg any) {
	h.logger.Infoln("request completed successfully")
	c.AbortWithStatusJSON(code, RespMsgSuccess{
		Success: true,
		Message: msg,
	})
}

func (h *Handler) getID(c *gin.Context) (id int, err error) {
	s := c.Params.ByName("id")
	subId := 0
	count, err := fmt.Sscanf(s, "%d", &subId)
	if count == 0 || err != nil {
		h.sendError(c, http.StatusNotFound, fmt.Errorf("incorrect sub id"))
		return subId, err
	}
	return subId, nil
}
