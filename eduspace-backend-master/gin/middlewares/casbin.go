package middlewares

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CasbinMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的 URL 和方法
		obj := c.Request.URL.Path
		act := c.Request.Method

		// 获取请求的用户（转为str）
		sub := c.GetString("role")

		if sub != "" {
			// 检查用户是否有访问权限
			ok, err := enforcer.Enforce(sub, obj, act)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			if !ok {
				// 用户无访问权限，返回错误
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "没有权限！"})
				return
			}
		}

		// 用户有访问权限，继续处理请求
		c.Next()
	}
}
