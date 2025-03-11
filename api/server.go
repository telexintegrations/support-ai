package api

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/telexintegrations/support-ai/aicom"
)

var (
	ErrRouterSetAlready = errors.New("router has aready been set")
)

type Server struct {
	EnvVar *EnvConfig
	Router *gin.Engine
	AIService aicom.AIService
}

func NewServer(envVar *EnvConfig) *Server {
	aiservice, _ := aicom.NewAIService(envVar.GenaiAPIKey)
	if aiservice == nil{
		fmt.Println("Unable to instantiate AI client")
	}else{
		fmt.Println("AI client initiated")
	}
	return &Server{
		EnvVar: envVar,
		Router: nil,
		AIService: aiservice,
	}
}

func (s *Server) SetupRouter() error {
	// Define and setup all routes and middleware here
	if s.Router != nil {
		return ErrRouterSetAlready
	}

	s.Router = gin.Default()

	r := s.Router
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Deployed and Running chatbot AI"})
	})

	r.GET("/upload", s.uploadPage)
	r.GET("/support-response", s.RaggedResponse)
	r.GET("/basic-response", s.BasicResponse)
	// r.GET("/integration",
	return nil
}

func (s *Server) StartServer(addr string) error {
	s.SetupRouter()
	if err := s.Router.Run(addr); err != nil {
		return err
	}
	return nil
}
