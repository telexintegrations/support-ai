package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/telexintegrations/support-ai/aicom"
	"github.com/telexintegrations/support-ai/internal/repository"
	mongoClient "github.com/telexintegrations/support-ai/internal/repository/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

const port = "8080"

var (
	ErrRouterSetAlready = errors.New("router has aready been set")
)

type Server struct {
	EnvVar    *EnvConfig
	Router    *gin.Engine
	AIService aicom.AIService
	DB        repository.VectorRepo
}

func NewServer(envVar *EnvConfig, db *mongo.Client) *Server {
	// Setup needed services...
	dbService := mongoClient.NewDBService(db)
	aiservice, _ := aicom.NewAIService(envVar.GenaiAPIKey)
	if aiservice == nil || dbService == nil {
		fmt.Println("Unable to instantiate AI client")
	} else {
		fmt.Println("AI & Repository Service Running")
	}
	return &Server{
		EnvVar:    envVar,
		Router:    nil,
		AIService: aiservice,
		DB:        dbService,
	}
}

func (s *Server) SetupRouter() error {
	// Define and setup all routes and middleware here
	if s.Router != nil {
		return ErrRouterSetAlready
	}
	// dbService := mongo.NewDBService(s.DB, s.EnvVar.MONGODATABASE_NAME)
	s.Router = gin.Default()
	r := s.Router

	// Setup cors
	corsConfig := cors.Config{
		AllowOrigins: []string{
			fmt.Sprintf("http://localhost:%s", port),
			"https://telex.im", "https://staging.telex.im",
			"http://telextest.im", "http://staging.telextest.im",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.OPTIONS("/*any", func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})
	// Serve Swagger UI static files
	r.Static("/swaggerui", "./static/swagger")
	// Download the Swagger File
	r.StaticFile("/swagger.yaml", "./static/swagger.yaml")
	r.Use(cors.New(corsConfig))
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "Deployed and Running chatbot AI"})
	})
	r.GET("/test-db", s.FetchEmbeddings)
	r.GET("/upload", s.uploadPage)
	r.POST("/upload", s.UploadFiles)
	r.GET("/integration.json", s.sendIntegrationJson)
	r.GET("/ngrok.json", s.sendNgrokJson)
	r.POST("/target", s.receiveChatQueries)
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
