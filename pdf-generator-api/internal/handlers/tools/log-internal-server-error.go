package tools

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogInternalServerError(c *gin.Context, handler string, method string, err error) {
	slog.Error(
		handler,
		slog.String("method", method),
		slog.String("error", err.Error()),
	)

	c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
}
