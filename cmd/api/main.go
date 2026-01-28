package main

import (
	"fmt"
	"log"
	"os"
	"sentinel/internal/database"
	"sentinel/internal/email"
	"sentinel/internal/server"

	"sentinel/internal/worker"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

}
func main() {

	dbPool, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer dbPool.Close()
	fmt.Println("🚀 Sentinel Database Connection Established")

	workerPool := worker.New(dbPool, 3)
	workerPool.Run()

	fmt.Println("⚡ Creating jobs in DB and sending to workers...")

	emailClient := email.NewClient(os.Getenv("EMAIL_APIKEY"))
	srv := server.NewServer(workerPool, emailClient)

	r := gin.Default()

	r.POST("/upload", srv.UploadHandler)
	r.POST("/register", srv.RegisterHandler)
	r.POST("/verify", srv.VerifyHandler)
	r.POST("/login", srv.LoginHandler)
	r.GET("/auth/google/login", srv.GoogleLoginHandler)
	r.GET("/auth/google/callback", srv.GoogleCallbackHandler)
	protected := r.Group("/api")
	protected.Use(srv.AuthMiddleware())
	{

		protected.GET("/profile", srv.ProfileHandler)
		protected.POST("/set-password", srv.SetPasswordHandler)

	}

	fmt.Println("Listening at 8081:")
	if err := r.Run(":8081"); err != nil {
		log.Fatal("Server failed to start: ", err)
	}

}
