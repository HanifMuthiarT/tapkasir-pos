package middleware

import (
	"net/http"
	"strings"

	"tapkasir-pos/pkg/security"

	"github.com/gin-gonic/gin"
)

// RequireRole memverifikasi token JWT dan memastikan role pengguna berhak mengakses endpoint sensitif
func RequireRole(jwtSecret string, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Akses ditolak: Header Authorization diperlukan (Bearer <token>)",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Format otorisasi tidak valid. Gunakan format: Bearer <token>",
			})
			c.Abort()
			return
		}

		claims, err := security.ValidateToken(parts[1], jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token tidak valid atau sudah kedaluwarsa. Silakan masukkan PIN kembali.",
			})
			c.Abort()
			return
		}

		if len(allowedRoles) > 0 {
			roleAllowed := false
			for _, allowed := range allowedRoles {
				if claims.Role == allowed {
					roleAllowed = true
					break
				}
			}

			if !roleAllowed {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "Akses ditolak: Anda tidak memiliki izin untuk fitur sensitif ini",
				})
				c.Abort()
				return
			}
		}

		c.Set("user_role", claims.Role)
		c.Next()
	}
}