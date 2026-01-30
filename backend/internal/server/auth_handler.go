package server

import (
	"context"
	"encoding/json"
	"fmt"
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
	pastr := string(hashedPwd)
	user := &models.User{
		Email:        req.Email,
		PasswordHash: &pastr,
	}

	err = database.CreateUser(s.WorkerPool.DB, user)

	otp, err := utils.GenerateOTP(6)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
	}
	subject := "🔐 Your Sentinel Verification Code"
	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Verification Code</title>
    <style>
        body { font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; }
        .container { max-width: 600px; margin: 0 auto; background-color: #ffffff; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        .header { text-align: center; padding-bottom: 20px; border-bottom: 1px solid #eeeeee; }
        .header h1 { color: #333333; margin: 0; font-size: 24px; }
        .content { padding: 20px 0; text-align: center; }
        .otp-box { background-color: #f0f8ff; border: 1px dashed #007bff; color: #007bff; font-size: 32px; font-weight: bold; letter-spacing: 5px; padding: 15px; margin: 20px auto; display: inline-block; border-radius: 5px; }
        .footer { text-align: center; font-size: 12px; color: #888888; margin-top: 20px; border-top: 1px solid #eeeeee; padding-top: 10px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Sentinel 🛡️</h1>
        </div>
        <div class="content">
            <p>Hello,</p>
            <p>Please use the verification code below to complete your registration:</p>

            <div class="otp-box">%s</div>

            <p>This code is valid for <strong>10 minutes</strong>.</p>
            <p>If you didn't request this, you can safely ignore this email.</p>
        </div>
        <div class="footer">
            <p>&copy; 2026 Sentinel Security. All rights reserved.</p>
        </div>
    </div>
</body>
</html>
`

	body := fmt.Sprintf(htmlTemplate, otp)
	OTPSendErr := s.EmailClient.Send(user.Email, subject, body)
	if OTPSendErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to send verification code"})
		return
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

func (s *Server) GoogleLoginHandler(c *gin.Context) {
	url := s.GoogleConfig.AuthCodeURL("state-token")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (s *Server) GoogleCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	token, err := s.GoogleConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Printf("❌ OAuth Exchange Error: %v\n", err) // Log error to console
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to exchange token", "details": err.Error()})
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	// 🔑 Use 'id' to map Google's unique identifier
	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	user, err := database.GetUserByEmail(s.WorkerPool.DB, googleUser.Email)
	if err != nil {
		// New Google users are automatically verified
		user = &models.User{
			Email:      googleUser.Email,
			IsVerified: true,
		}
		if err := database.CreateUser(s.WorkerPool.DB, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
			return
		}
	}

	// Link Google Identity
	database.AddUserIdentity(s.WorkerPool.DB, user.ID, "google", googleUser.ID)

	jwtToken, err := utils.GenerateToken(uint(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Redirect to frontend with token
	frontendURL := "http://localhost:5173/google-callback?token=" + jwtToken
	c.Redirect(http.StatusTemporaryRedirect, frontendURL)
}

type SetPasswordRequest struct {
	Password string `json:"password" binding:"required,min=6"`
}

func (s *Server) SetPasswordHandler(c *gin.Context) {
	// 🛡️ Use "user_id" and (uint) type to match AuthMiddleware
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := int(val.(uint))

	var req SetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	if err := database.UpdateUserPassword(s.WorkerPool.DB, userID, string(hashedPwd)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	// Register the 'local' identity method for this user
	database.AddUserIdentity(s.WorkerPool.DB, userID, "local", "")

	c.JSON(http.StatusOK, gin.H{"message": "Password set successfully"})
}
