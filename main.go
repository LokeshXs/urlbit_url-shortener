package main

import (
	"log"
	"os"
	"time"

	"github.com/Lokeshxs/url-shortener/db"
	"github.com/Lokeshxs/url-shortener/routes"
	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Loading the ENV variabels

	err := godotenv.Load()

	if err != nil {
		panic("Could not load the ENV variables!")
	}

	// Intializing the server
	server := gin.Default()

	// Connecting to Postgres DB
	if err := db.InitDB(); err != nil {
		log.Printf("Database unavailable; starting HTTP server anyway: %v", err)
	}
	if db.DB != nil {
		defer db.DB.Close()
	}

	// Configuring CORS
	server.Use(func(c *gin.Context) {
		if c.Request.URL.Path == "/webhook/signup" {
			cors.New(cors.Config{
				AllowOrigins: []string{"*"},
				AllowMethods: []string{"POST", "OPTIONS"},
				AllowHeaders: []string{"*"},
			})(c)
		} else {
			cors.New(cors.Config{
				AllowOrigins:     []string{"http://localhost:3000", "https://urlbit.space"},
				AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			})(c)
		}
	})

	// Set Clerk secret key once
	clerk.SetKey(os.Getenv("CLERK_SECRET_KEY"))

	// Calling a function to handle incoming requests
	routes.RoutingHandler(server)

	// Start Server
	server.Run(":3000")

}
