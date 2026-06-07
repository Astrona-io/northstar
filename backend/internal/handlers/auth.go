package handlers

import (
	"net/http"
	"time"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// JWTSecret is the secret key used for signing session tokens.
var JWTSecret = []byte("cmdb-secret-key-2026")

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type JWTCustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"` // For backward compatibility: admin, operator, viewer
	jwt.RegisteredClaims
}

// Login handles POST /api/auth/login (Phase 1 Advanced RBAC Update)
func Login(c echo.Context) error {
	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	// Fetch User from relational database (strict RBAC lookup)
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid username or password"})
	}

	// Compare password (with bcrypt support and plain-text fallback)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// Fallback check for plain-text password (if seeded/migrated differently)
		if user.Password != req.Password {
			return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid username or password"})
		}
	}

	// Set custom claims
	claims := &JWTCustomClaims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	t, err := token.SignedString(JWTSecret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to sign token"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token":    t,
		"username": user.Username,
		"role":     user.Role,
	})
}

// IsTestMode is a global flag that bypasses RBAC checks during integration testing.
var IsTestMode = false

// RBACMiddleware enforces fine-grained role and policy weighting security (Phase 1 Advanced RBAC)
func RBACMiddleware(permissionName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Bypass auth checks during unit/integration tests
			if IsTestMode {
				return next(c)
			}

			// Extract token from header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Missing authorization token"})
			}

			tokenString := authHeader
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			}

			token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				return JWTSecret, nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid or expired token"})
			}

			claims, ok := token.Claims.(*JWTCustomClaims)
			if !ok {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid token claims"})
			}

			// Superuser Admin role bypass (Admins have global override access)
			if claims.Role == "admin" {
				c.Set("username", claims.Username)
				c.Set("role", claims.Role)
				return next(c)
			}

			// Fetch full User profile preloaded with Groups and Override Permissions (Phase 1 Advanced RBAC)
			var user models.User
			err = database.DB.Preload("Group.Permissions").Preload("Permissions").Where("username = ?", claims.Username).First(&user).Error
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "User account not found"})
			}

			// Resolve Fine-Grained Policy and Weights (User override Group)
			allowed := false

			// 1. Inherit from Group policies
			if user.Group != nil {
				for _, p := range user.Group.Permissions {
					if p.Name == permissionName {
						allowed = (p.Effect == "allow")
					}
				}
			}

			// 2. Apply User Overrides (Strictly overrides Group policy! Policy weight override rule - Phase 1)
			for _, p := range user.Permissions {
				if p.Name == permissionName {
					allowed = (p.Effect == "allow") // Strictly overrides Group allow/deny!
				}
			}

			if !allowed {
				return c.JSON(http.StatusForbidden, echo.Map{"error": "Insufficient privileges: requires permission " + permissionName})
			}

			// Inject user info into context for tracing and user-specific auditing (Phase 2 Auditing)
			c.Set("username", claims.Username)
			c.Set("role", claims.Role)

			return next(c)
		}
	}
}
