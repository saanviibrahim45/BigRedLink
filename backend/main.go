package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"

	"bigredlink/routes"
	"bigredlink/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, relying on system environment")
	}

	projectID := os.Getenv("GCP_PROJECT_ID")
	credPath := os.Getenv("GCP_FIRESTORE_CREDENTIALS")
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID, option.WithCredentialsFile(credPath))
	if err != nil {
		log.Fatalf("Failed to create Firestore client: %v", err)
	}
	defer client.Close()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("firestoreClient", client)
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		client := c.MustGet("firestoreClient").(*firestore.Client)
		ctx := c.Request.Context()

		_, err := client.Collection("ping").Doc("hello").Set(ctx, map[string]string{"msg": "world"})
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		snap, err := client.Collection("ping").Doc("hello").Get(ctx)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, snap.Data())
	})

	auth := r.Group("/api/auth")
	routes.AuthRoutes(auth)

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	routes.ProtectedRoutes(protected)

	if err := r.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
