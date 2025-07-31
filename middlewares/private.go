package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/tnqbao/gau-authorization-service/config"
	"github.com/tnqbao/gau-authorization-service/utils"
)

func PrivateMiddleware(config *config.EnvConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		privateKey := c.GetHeader("Private-Key")

		if privateKey == "" {
			utils.JSON400(c, "Private key is required")
			c.Abort()
			return
		}

		if privateKey != config.PrivateKey {
			utils.JSON403(c, "Invalid private key")
			c.Abort()
			return
		}

		c.Next()
	}
}
