package routes

import (
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"

	"bigredlink/models"
)

func ProtectedRoutes(r *gin.RouterGroup) {
	r.GET("/user/me", func(c *gin.Context) {
		userIDIfc, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID missing in context"})
			return
		}

		userID, ok := userIDIfc.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID type"})
			return
		}

		client := c.MustGet("firestoreClient").(*firestore.Client)
		ctx := c.Request.Context()

		doc, err := client.Collection("users").Doc(userID).Get(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
			return
		}

		var user models.User
		if err := doc.DataTo(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user data"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":    user.ID,
			"email": user.Email,
		})
	})
}
