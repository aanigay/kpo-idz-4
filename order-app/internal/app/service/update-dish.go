package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"order-app/internal/app/models"
)

func (s *Service) UpdateDish(c *gin.Context) {
	var dish models.UpdateDishReq
	if err := c.ShouldBindJSON(&dish); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d := models.Dish{
		Id:          dish.Id,
		Name:        dish.Name,
		Description: dish.Description,
		Price:       dish.Price,
		Quantity:    dish.Quantity,
		IsAvailable: dish.IsAvailable,
	}
	query := "UPDATE dish SET name = $1, description = $2, price = $3, quantity = $4, is_available = $5 WHERE id = $6 RETURNING created_at, updated_at"
	err := s.db.QueryRow(query, dish.Name, dish.Description, dish.Price, dish.Quantity, dish.IsAvailable, dish.Id).Scan(&d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(http.StatusOK, d)
}
