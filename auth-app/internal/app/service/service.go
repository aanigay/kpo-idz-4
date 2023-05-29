package service

import (
	"database/sql"
	"github.com/gin-gonic/gin"
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

func NewService() *Service {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		dbUrl = "postgresql://root:root@localhost:5432/db?sslmode=disable"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
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
	log.Print(s.dbUrl)
	if err != nil {
		return err
	}
	s.db = db
	if err = s.db.Ping(); err != nil {
		return err
	}
	s.engine.Use(CORS())
	s.engine.POST("/signup", s.CreateUser)
	s.engine.POST("/login", s.Login)
	s.engine.GET("/logout", s.Logout)
	s.engine.GET("/info", s.GetInfo)

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
