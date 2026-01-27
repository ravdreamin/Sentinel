package server

import (
	"sentinel/internal/email"
	"sentinel/internal/worker"
)

type Server struct {
	WorkerPool  *worker.Pool
	EmailClient *email.Client
}

func NewServer(workerPool *worker.Pool, emailClient *email.Client) *Server {
	return &Server{
		WorkerPool:  workerPool,
		EmailClient: emailClient,
	}
}
