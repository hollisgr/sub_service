package handler

import (
	"fmt"
	"main/internal/model"
	"main/internal/subscription"
	"main/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	router     *gin.Engine
	subService subscription.SubscriptionInterface
	logger     *logger.Logger
}

func NewHandler(r *gin.Engine, s subscription.SubscriptionInterface, logger *logger.Logger) Handler_interface {
	initSwagger()
	return &Handler{
		router:     r,
		subService: s,
		logger:     logger,
	}
}

func (h *Handler) Register() {
	h.router.Use(CORSMiddleware())
	h.router.POST("/subscriptions", h.Create)
	h.router.GET("/subscriptions/:id", h.Read)
	h.router.PATCH("/subscriptions/:id", h.Update)
	h.router.DELETE("/subscriptions/:id", h.Delete)
	h.router.GET("/subscriptions", h.List)
	h.router.GET("/subscriptions/cost", h.Cost)
	h.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// Create godoc
//
//	@Summary		Create new subscription
//	@Description	Returns a new subscription object.
//	@Tags			Subscription
//	@Accept			json
//	@Produce		json
//	@Param			subscription	body		model.SubRequest	true	"Subscription create data"
//	@Success		200				{object}	handler.RespMsgSuccess
//	@Failure		400				{object}	handler.RespMsgError
//	@Failure		401				{object}	handler.RespMsgError
//	@Router			/subscriptions [post]
func (h *Handler) Create(c *gin.Context) {
	h.logger.Infoln("request to the create handler")
	ctx := c.Request.Context()
	sub := model.Subscription{}
	err := c.ShouldBindBodyWithJSON(&sub)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, fmt.Errorf("reading request body error"))
		return
	}
	id, err := h.subService.Save(ctx, sub)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, fmt.Errorf("creating sub error"))
		return
	}
	h.sendSuccess(c, http.StatusOK, fmt.Sprintf("created new sub with id: %d", id))
}

// Read godoc
//
//	@Summary		Read subscription by ID
//	@Description	Returns a subscription object.
//	@Tags			Subscription
//	@Produce		json
//	@Param			id	path		int	true	"Subscription ID"
//	@Success		200	{object}	handler.RespMsgSuccess
//	@Failure		400	{object}	handler.RespMsgError
//	@Failure		401	{object}	handler.RespMsgError
//	@Router			/subscriptions/{id} [get]
func (h *Handler) Read(c *gin.Context) {
	h.logger.Infoln("request to the read handler")
	subId, err := h.getID(c)
	if err != nil {
		return
	}
	ctx := c.Request.Context()
	sub, err := h.subService.Load(ctx, subId)
	if err != nil {
		h.sendError(c, http.StatusNotFound, fmt.Errorf("read error, sub not found"))
		return
	}
	h.sendSuccess(c, http.StatusOK, sub)
}

// Update godoc
//
//	@Summary		Update subscription by ID
//	@Description	Returns an ID updated subscription.
//	@Tags			Subscription
//	@Accept			json
//	@Produce		json
//	@Param			id				path		int					true	"Subscription ID"
//	@Param			subscription	body		model.SubRequest	true	"Subscription update data"
//	@Success		200				{object}	handler.RespMsgSuccess
//	@Failure		400				{object}	handler.RespMsgError
//	@Failure		401				{object}	handler.RespMsgError
//	@Router			/subscriptions/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	h.logger.Infoln("request to the update handler")
	subId, err := h.getID(c)
	if err != nil {
		return
	}
	ctx := c.Request.Context()
	sub := model.Subscription{}
	err = c.ShouldBindBodyWithJSON(&sub)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, fmt.Errorf("reading request body error"))
		return
	}
	sub.Id = subId
	err = h.subService.Update(ctx, sub)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, fmt.Errorf("update failed, sub not found"))
		return
	}
	h.sendSuccess(c, http.StatusOK, "sub updated")

}

// Delete godoc
//
//	@Summary		Delete subscription by ID
//	@Description	Returns an ID deleted subscription.
//	@Tags			Subscription
//	@Produce		json
//	@Param			id	path		int	true	"Subscription ID"
//	@Success		200	{object}	handler.RespMsgSuccess
//	@Failure		400	{object}	handler.RespMsgError
//	@Failure		401	{object}	handler.RespMsgError
//	@Router			/subscriptions/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	h.logger.Infoln("request to the delete handler")
	subId, err := h.getID(c)
	if err != nil {
		return
	}
	ctx := c.Request.Context()
	err = h.subService.Delete(ctx, subId)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, fmt.Errorf("delete failed, sub not found"))
		return
	}
	h.sendSuccess(c, http.StatusOK, "sub deleted")
}

// List godoc
//
//	@Summary		Read subscription list
//	@Description	Returns a list of subscription objects
//	@Tags			Subscription
//	@Produce		json
//	@Success		200	{object}	handler.RespMsgSuccess
//	@Failure		400	{object}	handler.RespMsgError
//	@Failure		401	{object}	handler.RespMsgError
//	@Router			/subscriptions [get]
func (h *Handler) List(c *gin.Context) {
	h.logger.Infoln("request to the list handler")
	ctx := c.Request.Context()
	subList, err := h.subService.LoadList(ctx)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, fmt.Errorf("sub list is empty"))
		return
	}
	h.sendSuccess(c, http.StatusOK, subList)
}

// Cost godoc
//
//	@Summary		Cost subscription
//	@Description	Returns a cost of subscriptions by user ID, date and service name
//	@Tags			Subscription
//	@Param			user_id			query	string	true	"User ID"
//	@Param			service_name	query	string	true	"Service name"
//	@Param			start			query	string	true	"Start date"
//	@Param			end				query	string	true	"End date"
//	@Produce		json
//	@Success		200	{object}	handler.RespMsgSuccess
//	@Failure		400	{object}	handler.RespMsgError
//	@Failure		401	{object}	handler.RespMsgError
//	@Router			/subscriptions/cost [get]
func (h *Handler) Cost(c *gin.Context) {
	h.logger.Infoln("request to the cost handler")
	serviceName := c.Query("service_name")
	s := c.Query("user_id")
	start := c.Query("start")
	end := c.Query("end")

	if start == "" || end == "" {
		h.sendError(c, http.StatusBadRequest, fmt.Errorf("both dates are required"))
		return
	}

	if serviceName == "" {
		h.sendError(c, http.StatusBadRequest, fmt.Errorf("service name is required"))
		return
	}

	userId, err := uuid.Parse(s)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, fmt.Errorf("wrong uuid"))
		return
	}

	ctx := c.Request.Context()

	data := model.CostRequest{
		UserId:      userId,
		ServiceName: serviceName,
		StartDate:   start,
		EndDate:     end,
	}

	cost, err := h.subService.Cost(ctx, data)

	if err != nil {
		h.sendError(c, http.StatusNotFound, fmt.Errorf("cost request error: %v", err))
		return
	}

	h.sendSuccess(c, http.StatusOK, cost)
}
