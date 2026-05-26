package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt/v5"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

const identityKey = "id"

var (
	port              string
	googleOauthConfig *oauth2.Config
	githubOauthConfig *oauth2.Config
	// Store OAuth state tokens to prevent CSRF attacks
	oauthStateStore = make(map[string]time.Time)
)

// User represents the user information from OAuth provider
type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Provider  string `json:"provider"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// GoogleUserInfo represents Google user information
type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

// GitHubUserInfo represents GitHub user information
type GitHubUserInfo struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Google OAuth2 Configuration
	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  fmt.Sprintf("http://localhost:%s/auth/google/callback", port),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	// GitHub OAuth2 Configuration
	githubOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  fmt.Sprintf("http://localhost:%s/auth/github/callback", port),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	// Clean up expired state tokens every 10 minutes
	go cleanupExpiredStates()
}

func main() {
	// Validate OAuth configuration
	if googleOauthConfig.ClientID == "" && githubOauthConfig.ClientID == "" {
		log.Println(
			"Warning: No OAuth providers configured. Set GOOGLE_CLIENT_ID/SECRET or GITHUB_CLIENT_ID/SECRET",
		)
	}

	engine := gin.Default()

	// Initialize JWT middleware
	authMiddleware, err := jwt.New(initJWTParams())
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	if err := authMiddleware.MiddlewareInit(); err != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + err.Error())
	}

	// Register routes
	registerRoute(engine, authMiddleware)

	// Start HTTP server
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("Server starting on http://localhost:%s", port)
	log.Printf("Demo Page: http://localhost:%s/demo", port)
	log.Println("\nOAuth Login URLs:")
	if googleOauthConfig.ClientID != "" {
		log.Printf("  Google: http://localhost:%s/auth/google/login", port)
	}
	if githubOauthConfig.ClientID != "" {
		log.Printf("  GitHub: http://localhost:%s/auth/github/login", port)
	}
	if googleOauthConfig.ClientID == "" && githubOauthConfig.ClientID == "" {
		log.Println("  (No OAuth providers configured)")
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func registerRoute(r *gin.Engine, handle *jwt.GinJWTMiddleware) {
	_ = "STUB: not implemented"
	// Enable CORS for development
	return
}

// Public routes

// OAuth login initiation

// OAuth callbacks

// JWT token refresh

// Protected routes

func initJWTParams() *jwt.GinJWTMiddleware { _ = "STUB: not implemented"; return nil }

// Enable built-in cookie support for better security

// Set to true in production with HTTPS

// Send Authorization header in LoginResponse

// Custom LoginResponse to handle OAuth redirect vs regular JSON response

// OAuth: redirect to demo page (token already set in cookie by gin-jwt)

// Regular login: return JSON response

func getSecretKey() string { _ = "STUB: not implemented"; return "" }

func payloadFunc() func(data any) gojwt.MapClaims { _ = "STUB: not implemented"; return nil }

func identityHandler() func(c *gin.Context) any { _ = "STUB: not implemented"; return nil }

func authenticator() func(c *gin.Context) (any, error) { _ = "STUB: not implemented"; return nil }

// This is not used for OAuth flow, but required by the middleware
// OAuth authentication happens in the callback handlers

func authorizer() func(c *gin.Context, data any) bool { _ = "STUB: not implemented"; return nil }

// All authenticated OAuth users are authorized

func unauthorized() func(c *gin.Context, code int, message string) {
	_ = "STUB: not implemented"
	return nil
}

func logoutResponse() func(c *gin.Context) { _ = "STUB: not implemented"; return nil }

// handleOAuthSuccess handles the successful OAuth authentication
// and uses gin-jwt's built-in features (SendCookie, LoginResponse, etc.)
func handleOAuthSuccess(
	c *gin.Context,
	authMiddleware *jwt.GinJWTMiddleware,
	user *User,
	provider string,
) error {
	_ = "STUB: not implemented"
	// Set user identity in context (for middleware callbacks)
	return nil
}

// Generate JWT token

// Set cookies (both access token and refresh token)

// Let gin-jwt handle everything (cookies, headers, response) via LoginResponse
// The middleware will automatically:
// - Set httpOnly cookie (if SendCookie is enabled)
// - Set Authorization header (if SendAuthorization is enabled)
// - Call LoginResponse callback (defined in initJWTParams)

// OAuth handlers
func handleGoogleLogin(c *gin.Context) { _ = "STUB: not implemented"; return }

func handleGitHubLogin(c *gin.Context) { _ = "STUB: not implemented"; return }

func handleGoogleCallback(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	_ = "STUB: not implemented"
	return *new(gin.HandlerFunc)
}

// Validate state token (CSRF protection)

// Exchange authorization code for access token

// Get user info from Google

// Create user object

// Handle OAuth success with proper JWT middleware integration

func handleGitHubCallback(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	_ = "STUB: not implemented"
	return *new(gin.HandlerFunc)
}

// Validate state token (CSRF protection)

// Exchange authorization code for access token

// Get user info from GitHub

// Get user email if not public

// Create user object

// Handle OAuth success with proper JWT middleware integration

// Helper functions
func generateStateToken() string { _ = "STUB: not implemented"; return "" }

func validateStateToken(state string) bool { _ = "STUB: not implemented"; return false }

func cleanupExpiredStates() { _ = "STUB: not implemented"; return }

func getUserEmail(client *http.Client) string { _ = "STUB: not implemented"; return "" }

func getStringFromClaims(claims gojwt.MapClaims, key string) string {
	_ = "STUB: not implemented"
	return ""
}

// CORS middleware for development
func corsMiddleware() gin.HandlerFunc { _ = "STUB: not implemented"; return *new(gin.HandlerFunc) }

// Route handlers
func indexHandler(c *gin.Context) { _ = "STUB: not implemented"; return }

func profileHandler(c *gin.Context) { _ = "STUB: not implemented"; return }

func handleNoRoute() func(c *gin.Context) { _ = "STUB: not implemented"; return nil }
