package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"order-app/internal/app/models"
)

func (s *Service) CreateDish(c *gin.Context) {
	var r models.CreateDishReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d := models.Dish{
		Name:        r.Name,
		Description: r.Description,
		Price:       r.Price,
		Quantity:    r.Quantity,
		IsAvailable: r.IsAvailable,
	}
	var lastInsertId int
	query := "INSERT INTO \"dish\"(name, description, price, quantity, is_available) VALUES ($1, $2, $3, $4, $5) returning id, created_at, updated_at"
	err := s.db.QueryRow(query, r.Name, r.Description, r.Price, r.Quantity, r.IsAvailable).Scan(&lastInsertId, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		c.JSON(500, err)
		return
	}

	d.Id = int64(lastInsertId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, d)
}
