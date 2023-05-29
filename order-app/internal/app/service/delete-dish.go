package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"order-app/internal/app/models"
	"strconv"
)

func (s *Service) DeleteDish(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	d := models.Dish{
		Id: int64(id),
	}
	query := "DELETE FROM dish WHERE id = $1 RETURNING name, description, price, quantity, is_available, created_at, updated_at"
	err = s.db.QueryRow(query, id).Scan(&d.Name, &d.Description, &d.Price, &d.Quantity, &d.IsAvailable, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(http.StatusOK, d)
}
