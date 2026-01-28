package server

import (
	"os"
	"sentinel/internal/email"
	"sentinel/internal/worker"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Server struct {
	WorkerPool   *worker.Pool
	EmailClient  *email.Client
	GoogleConfig *oauth2.Config
}

func NewServer(workerPool *worker.Pool, emailClient *email.Client) *Server {
	return &Server{
		WorkerPool:  workerPool,
		EmailClient: emailClient,
		GoogleConfig: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Endpoint:     google.Endpoint,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		},
	}
}
