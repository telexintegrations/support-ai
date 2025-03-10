package api

import "github.com/gin-gonic/gin"

// Used to upload documents for testing our application
func (s *Server) uploadPage(ctx *gin.Context) {
	ctx.HTML(200, "tmpuploader.html", gin.H{
		"message": "successful",
	})
}
