package service

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/lib/pq"
	"log"
	"os"
)

const (
	secretKey = "secret"
)

type Service struct {
	engine *gin.Engine
	dbUrl  string
	port   string
	db     *sql.DB
}

type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func NewService() *Service {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		dbUrl = "postgresql://root:root@localhost:5432/db?sslmode=disable"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8004"
	}

	return &Service{
		engine: gin.Default(),
		dbUrl:  dbUrl,
		port:   port,
	}
}

func (s *Service) Run() error {

	db, err := sql.Open("postgres",
		s.dbUrl)
	if err != nil {
		return err
	}
	s.db = db
	if err = s.db.Ping(); err != nil {
		return err
	}
	s.engine.Use(CORS())
	s.engine.GET("/get_dish", s.GetDish)
	s.engine.GET("/get_dishes", s.GetAll)

	s.engine.Use(s.UseId())
	s.engine.POST("/create_order", s.Order)
	s.engine.GET("/get_order", s.GetOrder)
	s.engine.Use(s.CheckRole("manager"))
	s.engine.POST("/create_dish", s.CreateDish)
	s.engine.PUT("/update_dish", s.UpdateDish)
	s.engine.DELETE("/delete_dish", s.DeleteDish)

	err = s.engine.Run(":" + s.port)
	if err != nil {
		return err
	}
	return nil
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func (s *Service) GetClaimsFromToken(tokenString string) (*Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Print(err)
		log.Print(token.Claims)
		return nil, err
	}
	log.Print(token.Valid)
	return &claims, nil
}

func (s *Service) CheckRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("jwt")
		if err != nil {
			c.JSON(401, err)
			c.Abort()
			return
		}
		claims, err := s.GetClaimsFromToken(token)
		if err != nil {
			c.JSON(500, err)
			c.Abort()
			return
		}
		if claims.Role != role {
			c.JSON(401, err)
			c.Abort()
			return
		}
		c.Next()
	}
}

func (s *Service) UseId() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("jwt")
		if err != nil {
			c.JSON(401, err)
			c.Abort()
			return
		}

		claims, err := s.GetClaimsFromToken(token)
		if err != nil {
			c.JSON(500, err)
			c.Abort()
			return
		}
		c.Set("user_id", claims.ID)
		c.Next()
	}
}
