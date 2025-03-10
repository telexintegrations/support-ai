package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const port = "8080"

var (
	ErrRouterSetAlready = errors.New("router has aready been set")
)

type Server struct {
	EnvVar *EnvConfig
	Router *gin.Engine
}

func NewServer(envVar *EnvConfig) *Server {
	return &Server{
		EnvVar: envVar,
		Router: nil,
	}
}

func (s *Server) SetupRouter() error {
	// Define and setup all routes and middleware here
	if s.Router != nil {
		return ErrRouterSetAlready
	}

	s.Router = gin.Default()
	r := s.Router

	// Setup cors
	corsConfig := cors.Config{
		AllowOrigins:     []string{fmt.Sprintf("http://localhost:%s", port), "https://telex.im"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(corsConfig))
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Deployed and Running chatbot AI"})
	})

	r.GET("/upload", s.uploadPage)
	// r.GET("/integration",
	return nil
}

func (s *Server) StartServr(addr string) error {
	s.SetupRouter()
	if err := s.Router.Run(addr); err != nil {
		return err
	}
	return nil
}
