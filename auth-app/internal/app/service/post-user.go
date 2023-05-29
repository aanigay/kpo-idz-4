package service

import (
	"auth-app/internal/app/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/mail"
	"strconv"
	"time"
)

func validateEmail(email string) (string, bool) {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", false
	}
	return addr.Address, true
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *Service) CreateUser(c *gin.Context) {
	var u models.CreateUserReq
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	addr, ok := validateEmail(u.Email)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email", "email": addr})
		return
	}
	u.Email = addr
	hashedPassword, err := HashPassword(u.Password)
	if err != nil {
		c.JSON(500, err)
		return
	}

	user := &models.User{
		Username: u.Username,
		Email:    u.Email,
		Password: hashedPassword,
		Role:     u.Role,
	}
	var lastInsertId int
	query := "INSERT INTO \"user\"(username, password_hash, email, role) VALUES ($1, $2, $3, $4) returning id, created_at, updated_at"
	err = s.db.QueryRow(query, user.Username, user.Password, user.Email, user.Role).Scan(&lastInsertId, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		c.JSON(500, err)
		return
	}

	user.Id = int64(lastInsertId)

	res := &models.CreateUserRes{
		Id:       strconv.Itoa(int(user.Id)),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (s *Service) Login(c *gin.Context) {
	var user models.LoginUserReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//u, err := h.Service.Login(c.Request.Context(), &user)

	//u, err := s.Repository.GetUserByEmail(req.Email)
	u := models.User{}
	query := "SELECT id, email, username, password_hash, role, created_at, updated_at FROM \"user\" WHERE email = $1"
	err := s.db.QueryRow(query, user.Email).Scan(&u.Id, &u.Email, &u.Username, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = CheckPassword(user.Password, u.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	expiresAt := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		ID:       u.Id,
		Username: u.Username,
		Role:     u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(u.Id)),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	session := models.Session{
		UserId:       u.Id,
		ExpiresAt:    expiresAt,
		SessionToken: ss,
	}
	var id int64
	query = "INSERT INTO session (user_id, session_token, expires_at) VALUES ($1, $2, $3) RETURNING id"
	err = s.db.QueryRow(query, &session.UserId, &session.SessionToken, &session.ExpiresAt).Scan(&id)
	if err != nil {
		c.JSON(500, err)
		return
	}
	//
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res := models.LoginUserRes{
		Id:       u.Id,
		Username: u.Username,
		Email:    u.Email,
	}

	c.SetCookie("jwt", ss, 60*60*24, "/", "localhost", false, true)
	c.JSON(http.StatusOK, res)
}

type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func (s *Service) GetClaimsFromToken(tokenString string) (*Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		log.Print(token.Claims.(jwt.Claims))
		return nil, err
	}
	log.Print(token.Valid)
	return &claims, nil
}
