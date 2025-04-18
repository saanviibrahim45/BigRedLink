package controllers

import (
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"bigredlink/models"
	"bigredlink/utils"
)


type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	client := c.MustGet("firestoreClient").(*firestore.Client)
	ctx := c.Request.Context()

	docs, err := client.Collection("users").Where("email", "==", req.Email).Limit(1).Documents(ctx).GetAll()
	if err != nil || len(docs) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	var user models.User
	if err := docs[0].DataTo(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	accessToken, err := utils.CreateAccessToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}
	refreshToken, err := utils.CreateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	if _, err := client.Collection("users").Doc(user.ID).Update(ctx, []firestore.Update{{Path: "refreshToken", Value: refreshToken}}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store refresh token"})
		return
	}

	c.SetCookie("refresh_token", refreshToken, int((7*24*time.Hour).Seconds()), "/", "", true, true)

	c.JSON(http.StatusOK, gin.H{"accessToken": accessToken})
}

func Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
		return
	}

	token, err := utils.ParseToken(refreshToken, true)
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}
	userID, ok := claims["sub"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token subject"})
		return
	}

	client := c.MustGet("firestoreClient").(*firestore.Client)
	ctx := c.Request.Context()

	doc, err := client.Collection("users").Doc(userID).Get(ctx)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	var user models.User
	if err := doc.DataTo(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Server error"})
		return
	}

	if user.RefreshToken != refreshToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token mismatch"})
		return
	}

	newAccess, err := utils.CreateAccessToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accessToken": newAccess})
}

func Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil {
		tkn, perr := utils.ParseToken(refreshToken, true)
		if perr == nil && tkn.Valid {
			if claims, ok := tkn.Claims.(jwt.MapClaims); ok {
				if userID, ok := claims["sub"].(string); ok {
					client := c.MustGet("firestoreClient").(*firestore.Client)
					ctx := c.Request.Context()
					client.Collection("users").Doc(userID).Update(ctx, []firestore.Update{{Path: "refreshToken", Value: ""}})
				}
			}
		}
	}

	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	c.Status(http.StatusNoContent)
}
