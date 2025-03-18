package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"

	"EDU_TH_2_backend/gin/utils"
)

func JWTAuthMiddleware(tokenUtil utils.TokenUtil) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Authorization header is missing"})
			return
		}

		// 从 Authorization header 中解析出 token
		tokenStrings := strings.Split(authHeader, "Bearer ")
		if len(tokenStrings) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header format"})
			return
		}

		tokenString := strings.TrimSpace(tokenStrings[1])
		token, err := tokenUtil.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token claims"})
			return
		}

		// 将用户ID放入Context中，以便后续使用
		userId := claims["id"].(float64)
		role := claims["role"].(string)
		c.Set("userId", int(userId))      // 将用户ID放入Context中，以便后续使用
		c.Set("role", role)               // 将用户Role放入Context中，以便后续使用
		c.Set("tokenString", tokenString) // 将tokenString放入Context中，以便后续使用
		c.Next()
	}
}
