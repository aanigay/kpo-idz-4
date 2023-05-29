package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"order-app/internal/app/models"
	"strconv"
)

func (s *Service) GetOrder(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	order := models.Order{}
	query := "SELECT * FROM \"order\" WHERE id = $1"
	err = s.db.QueryRow(query, int64(id)).Scan(&order.Id, &order.UserId, &order.Status, &order.SpecialRequests, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}
