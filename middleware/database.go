package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/Lokeshxs/url-shortener/db"
	"github.com/gin-gonic/gin"
)

func RequireDatabase(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), db.ConnectionTimeout)
	defer cancel()
	if err := db.EnsureReady(ctx); err != nil {
		log.Printf("Database unavailable: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}
	c.Next()
}
