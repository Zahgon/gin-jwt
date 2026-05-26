package main

import (
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt/v5"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

const userAdmin = "admin"

var (
	identityKey = "id"
	port        string
)

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

func init() {
	port = os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
}

func main() {
	engine := gin.Default()
	// the jwt middleware
	authMiddleware, err := jwt.New(initParams())
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	// initialize middleware
	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	// register route
	registerRoute(engine, authMiddleware)

	// start http server with proper timeouts
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           engine,
		ReadHeaderTimeout: 5 * time.Second,
	}
	if err = srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func registerRoute(r *gin.Engine, handle *jwt.GinJWTMiddleware) {
	_ = "STUB: not implemented"
	// Public routes
	return
}

// RFC 6749 compliant refresh endpoint

// Protected routes

// Logout with refresh token revocation

func initParams() *jwt.GinJWTMiddleware { _ = "STUB: not implemented"; return nil }

// TokenLookup: "query:token",
// TokenLookup: "cookie:token",

func payloadFunc() func(data any) gojwt.MapClaims { _ = "STUB: not implemented"; return nil }

func identityHandler() func(c *gin.Context) any { _ = "STUB: not implemented"; return nil }

func authenticator() func(c *gin.Context) (any, error) { _ = "STUB: not implemented"; return nil }

func authorizator() func(c *gin.Context, data any) bool { _ = "STUB: not implemented"; return nil }

func unauthorized() func(c *gin.Context, code int, message string) {
	_ = "STUB: not implemented"
	return nil
}

func logoutResponse() func(c *gin.Context) { _ = "STUB: not implemented"; return nil }

// This demonstrates that claims are now accessible during logout

// Show that we can access user information during logout

func handleNoRoute() func(c *gin.Context) { _ = "STUB: not implemented"; return nil }

func helloHandler(c *gin.Context) { _ = "STUB: not implemented"; return }
