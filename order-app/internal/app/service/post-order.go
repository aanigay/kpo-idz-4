package service

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"order-app/internal/app/models"
)

func (s *Service) Order(c *gin.Context) {
	var r struct {
		SpecialRequests string  `json:"special_requests"`
		DishIds         []int64 `json:"dish_ids"`
		Quantities      []int64 `json:"quantities"`
	}
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(r.DishIds) != len(r.Quantities) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids and quantities have different sizes"})
		return
	}
	req := models.CreateOrderReq{
		UserId:          c.MustGet("user_id").(int64),
		SpecialRequests: r.SpecialRequests,
		DishIds:         r.DishIds,
		Quantities:      r.Quantities,
	}
	order, err := s.CreateOrder(&req)
	if err != nil {
		log.Print(order)
		c.JSON(500, err)
	}
	for index, id := range req.DishIds {
		price, err := s.GetPrice(id)
		if err != nil {
			c.JSON(500, err)
		}

		_, err = s.CreateOrderDish(&models.CreateOrderDishReq{
			OrderId:  order.Id,
			DishId:   id,
			Quantity: req.Quantities[index],
			Price:    price,
		})
		if err != nil {
			c.JSON(500, err)
		}
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go s.SetStatus(&models.UpdateStatusReq{
		Id:     order.Id,
		Status: "done",
	})
	c.JSON(http.StatusOK, order)
}

func (s *Service) CreateOrder(req *models.CreateOrderReq) (*models.Order, error) {
	order := models.Order{
		Status:          "in_progress",
		UserId:          req.UserId,
		SpecialRequests: req.SpecialRequests,
	}
	query := "INSERT INTO \"order\" (user_id, status, special_requests) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at"
	err := s.db.QueryRow(query, order.UserId, order.Status, order.SpecialRequests).Scan(&order.Id, &order.CreatedAt, &order.UpdatedAt)
	return &order, err
}
func (s *Service) GetPrice(id int64) (float64, error) {
	var price float64
	query := "SELECT price FROM dish WHERE id = $1"
	err := s.db.QueryRow(query, id).Scan(&price)
	return price, err
}

func (s *Service) SetStatus(req *models.UpdateStatusReq) {
	log.Printf("try update where id = %d status = %s", req.Id, req.Status)
	//r.db.QueryRowContext(ctx, "UPDATE \"order\" SET status = $1 WHERE id = $2", req.Status, req.Id)
	db, err := sql.Open("postgres", "postgresql://root:root@localhost:5432/db?sslmode=disable")
	if err != nil {
		log.Print(err)
		return
	}
	db.QueryRow("UPDATE \"order\" SET status = $1 WHERE id = $2", req.Status, req.Id)

}

func (s *Service) CreateOrderDish(req *models.CreateOrderDishReq) (*models.OrderDish, error) {
	orderDish := models.OrderDish{
		OrderId:  req.OrderId,
		DishId:   req.DishId,
		Quantity: req.Quantity,
		Price:    req.Price,
	}
	query := "INSERT INTO \"order_dish\" (order_id, dish_id, quantity, price) VALUES ($1, $2, $3, $4) RETURNING id"
	err := s.db.QueryRow(query, orderDish.OrderId, orderDish.DishId, orderDish.Quantity, orderDish.Price).Scan(&orderDish.Id)
	return &orderDish, err
}
