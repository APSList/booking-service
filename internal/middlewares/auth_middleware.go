package middlewares

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	jwksURL string
}

func NewAuthMiddleware() AuthMiddleware {
	return AuthMiddleware{
		jwksURL: "https://frauwrkbphmjngymcdyk.supabase.co/auth/v1/.well-known/jwks.json",
	}
}

func (m AuthMiddleware) Handler() gin.HandlerFunc {
	// Initialize the key function (it handles caching the public key for you)
	k, err := keyfunc.NewDefault([]string{m.jwksURL})
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch JWKS: %v", err))
	}

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, k.Keyfunc,
			jwt.WithValidMethods([]string{"ES256"}),
			jwt.WithLeeway(5*time.Minute),
		)

		if err != nil {
			fmt.Println("--- JWT DEBUG START ---")
			fmt.Printf("Error: %v\n", err)
			fmt.Printf("Token String: %s\n", tokenString)
			fmt.Println("--- JWT DEBUG END ---")

			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid Token Signature",
				"details": err.Error(), // Send the real error back to Postman for a moment
			})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userMeta := claims["user_metadata"].(map[string]interface{})

		// Pass the data to your controller
		c.Set("organization_id", int64(userMeta["organization_id"].(float64)))
		c.Set("role", userMeta["role"].(string))

		c.Next()
	}
}
