package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Service) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

func (s *Service) GetInfo(c *gin.Context) {
	token, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(401, "NOT AUTHORIZED")
		return
	}
	claims, err := s.GetClaimsFromToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, claims)
}
