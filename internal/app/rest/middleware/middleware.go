package middleware

import (
	"net/http"
	"prodigo/internal/app/rest/casbin"
	"prodigo/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

type Middleware struct {
	maker    jwt.TokenMaker
	enforcer casbin.Enforcer
}

func New(maker jwt.TokenMaker, enforcer casbin.Enforcer) *Middleware {
	return &Middleware{maker: maker, enforcer: enforcer}
}

func (m *Middleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		const (
			authScheme = "Bearer"
			authHeader = "Authorization"
		)

		auth := c.GetHeader(authHeader)
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing auth header"})
			return
		}

		scheme, token, ok := strings.Cut(auth, " ")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
			return
		}

		if scheme != authScheme {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth scheme"})
			return
		}

		claims, err := m.maker.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if len(claims.Audience) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token audience"})
			return
		}

		sub := claims.Audience[0]
		obj := c.Request.URL.Path
		act := c.Request.Method

		ok, err = m.enforcer.Enforce(sub, obj, act)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}
