package controller

import (
	"errors"
	"net/http"
	"strconv"

	"order-management/model"
	"order-management/repository"
	"order-management/service"

	"github.com/gin-gonic/gin"
)

type OrderController struct{ svc *service.OrderService }

func NewOrderController(s *service.OrderService) *OrderController { return &OrderController{svc: s} }
func (c *OrderController) Register(r *gin.RouterGroup) {
	r.POST("/orders", c.Create)
	r.GET("/orders", c.List)
	r.GET("/orders/:id", c.Get)
	r.PUT("/orders/:id/status", c.UpdateStatus)
}
func (c *OrderController) Create(ctx *gin.Context) {
	var req model.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON request"})
		return
	}
	o, err := c.svc.Create(ctx.Request.Context(), req)
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, o)
}
func (c *OrderController) Get(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}
	o, err := c.svc.Get(ctx.Request.Context(), id)
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, o)
}
func (c *OrderController) List(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	orders, err := c.svc.List(ctx.Request.Context(), limit, offset)
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"items": orders, "limit": limit, "offset": offset})
}
func (c *OrderController) UpdateStatus(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}
	var req model.UpdateStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON request"})
		return
	}
	o, err := c.svc.UpdateStatus(ctx.Request.Context(), id, req.Status)
	if err != nil {
		writeError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, o)
}
func parseID(ctx *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return 0, false
	}
	return id, true
}
func writeError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrInvalid):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid order input"})
	case errors.Is(err, repository.ErrNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
	case errors.Is(err, service.ErrProductUnavailable):
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "product service temporarily unavailable"})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
