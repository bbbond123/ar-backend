package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"gorm.io/gorm"
)

type Server struct {
	port   int
	gormDB *gorm.DB
}

func NewServer(gormDB *gorm.DB) *http.Server {
	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil || port == 0 {
		port = 3000 // 默认端口3000
	}

	NewServer := &Server{
		port:   port,
		gormDB: gormDB,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
