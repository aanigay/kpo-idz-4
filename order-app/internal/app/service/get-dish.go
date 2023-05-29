package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"order-app/internal/app/models"
	"strconv"
)

func (s *Service) GetDish(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	d := models.Dish{}
	query := "SELECT id, name, description, price, quantity, is_available, created_at, updated_at FROM \"dish\" WHERE id = $1"
	err = s.db.QueryRow(query, id).Scan(&d.Id, &d.Name, &d.Description, &d.Price, &d.Quantity, &d.IsAvailable, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		c.JSON(500, err)
		return
	}
	c.JSON(http.StatusOK, d)
}

func (s *Service) GetAll(c *gin.Context) {
	dishes := make([]models.Dish, 0)
	query := "SELECT * FROM dish"
	rows, err := s.db.Query(query)
	if err != nil {
		c.JSON(500, err)
	}
	for rows.Next() {
		var dish models.Dish
		err = rows.Scan(&dish.Id, &dish.Name, &dish.Description, &dish.Price, &dish.Quantity, &dish.IsAvailable, &dish.CreatedAt, &dish.UpdatedAt)
		if err != nil {
			c.JSON(500, err)
		}
		dishes = append(dishes, dish)
	}
	c.JSON(200, dishes)
}
