package server

import (
	"net/http"
	"sentinel/internal/database"
	"sentinel/internal/models"
	"sentinel/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

func (s *Server) RegisterHandler(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPwd),
	}

	err = database.CreateUser(s.WorkerPool.DB, user)

	otp, err := utils.GenerateOTP(6)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
	}

	expireAt := time.Now().Add(10 * time.Minute)

	verification := &models.Verification{
		UserId:   user.ID,
		Code:     otp,
		ExpireAt: expireAt,
	}

	err = database.SaveVerification(s.WorkerPool.DB, verification)
	if err != nil {
		// In a real app, we might want to rollback the user creation here
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save verification code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful. Please verify your email."})
}

func (s *Server) VerifyHandler(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := database.GetUserByEmail(s.WorkerPool.DB, req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	storedAuth, err := database.GetVerification(s.WorkerPool.DB, user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No pending verification found"})
		return
	}

	if storedAuth.Code != req.Code {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid verification code"})
		return
	}

	if time.Now().After(storedAuth.ExpireAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Code has expired"})
		return
	}

	if err := database.MarkUserVerified(s.WorkerPool.DB, user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
		return
	}

	database.DeleteVerification(s.WorkerPool.DB, user.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully!"})
}
