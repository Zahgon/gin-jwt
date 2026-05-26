package jwt

import (
	"context"
	"crypto/rsa"
	"errors"
	"net/http"
	"time"

	"github.com/appleboy/gin-jwt/v3/core"
	"github.com/appleboy/gin-jwt/v3/store"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	tokenContextKey         = "JWT_TOKEN"
	algRS256                = "RS256"
	keyCode                 = "code"
	keyMessage              = "message"
	defaultRealm            = "gin jwt"
	defaultCookieName       = "jwt"
	defaultRefreshTokenName = "refresh_token"
	claimExp                = "exp"
	tokenLookupCookie       = "cookie"
)

// GinJWTMiddleware provides a Json-Web-Token authentication implementation. On failure, a 401 HTTP response
// is returned. On success, the wrapped middleware is called, and the userID is made available as
// c.Get("userID").(string).
// Users can get a token by posting a json request to LoginHandler. The token then needs to be passed in
// the Authentication header. Example: Authorization:Bearer XXX_TOKEN_XXX
type GinJWTMiddleware struct {
	// Realm name to display to the user. Required.
	Realm string

	// signing algorithm - possible values are HS256, HS384, HS512, RS256, RS384 or RS512
	// Optional, default is HS256.
	SigningAlgorithm string

	// Secret key used for signing. Required.
	Key []byte

	// Callback to retrieve key used for signing. Setting KeyFunc will bypass
	// all other key settings
	KeyFunc func(token *jwt.Token) (any, error)

	// Duration that a jwt token is valid. Optional, defaults to one hour.
	Timeout time.Duration
	// Callback function that will override the default timeout duration.
	TimeoutFunc func(data any) time.Duration

	// This field allows clients to refresh their token until MaxRefresh has passed.
	// Note that clients can refresh their token in the last moment of MaxRefresh.
	// This means that the maximum validity timespan for a token is TokenTime + MaxRefresh.
	// Optional, defaults to 0 meaning not refreshable.
	MaxRefresh time.Duration

	// Callback function that should perform the authentication of the user based on login info.
	// Must return user data as user identifier, it will be stored in Claim Array. Required.
	// Check error (e) to determine the appropriate error message.
	Authenticator func(c *gin.Context) (any, error)

	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	Authorizer func(c *gin.Context, data any) bool

	// Callback function that will be called during login.
	// Using this function it is possible to add additional payload data to the webtoken.
	// The data is then made available during requests via c.Get("JWT_PAYLOAD").
	// Note that the payload is not encrypted.
	// The attributes mentioned on jwt.io can't be used as keys for the map.
	// Optional, by default no additional data will be set.
	PayloadFunc func(data any) jwt.MapClaims

	// User can define own Unauthorized func.
	Unauthorized func(c *gin.Context, code int, message string)

	// User can define own LoginResponse func.
	LoginResponse func(c *gin.Context, token *core.Token)

	// User can define own LogoutResponse func.
	LogoutResponse func(c *gin.Context)

	// User can define own RefreshResponse func.
	RefreshResponse func(c *gin.Context, token *core.Token)

	// Set the identity handler function
	IdentityHandler func(*gin.Context) any

	// Set the identity key
	IdentityKey string

	// TokenLookup is a string in the form of "<source>:<name>" that is used
	// to extract token from the request.
	// Optional. Default value "header:Authorization".
	// Possible values:
	// - "header:<name>"
	// - "query:<name>"
	// - "cookie:<name>"
	TokenLookup string

	// TokenHeadName is a string in the header. Default value is "Bearer"
	TokenHeadName string

	// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
	TimeFunc func() time.Time

	// HTTP Status messages for when something in the JWT middleware fails.
	// Check error (e) to determine the appropriate error message.
	HTTPStatusMessageFunc func(c *gin.Context, e error) string

	// Private key file for asymmetric algorithms
	PrivKeyFile string

	// Private Key bytes for asymmetric algorithms
	//
	// Note: PrivKeyFile takes precedence over PrivKeyBytes if both are set
	PrivKeyBytes []byte

	// Public key file for asymmetric algorithms
	PubKeyFile string

	// Private key passphrase
	PrivateKeyPassphrase string

	// Public key bytes for asymmetric algorithms.
	//
	// Note: PubKeyFile takes precedence over PubKeyBytes if both are set
	PubKeyBytes []byte

	// Private key
	privKey *rsa.PrivateKey

	// Public key
	pubKey *rsa.PublicKey

	// Optionally return the token as a cookie
	SendCookie bool

	// Duration that a cookie is valid. Optional, by default equals to Timeout value.
	CookieMaxAge time.Duration

	// Allow insecure cookies for development over http
	SecureCookie bool

	// Allow cookies to be accessed client side for development
	CookieHTTPOnly bool

	// Allow cookie domain change for development
	CookieDomain string

	// SendAuthorization allow return authorization header for every request
	SendAuthorization bool

	// Disable abort() of context.
	DisabledAbort bool

	// CookieName allow cookie name change for development
	CookieName string

	// RefreshTokenCookieName allow refresh token cookie name change for development
	RefreshTokenCookieName string

	// CookieSameSite allow use http.SameSite cookie param
	CookieSameSite http.SameSite

	// ParseOptions allow to modify jwt's parser methods.
	// WithTimeFunc is always added to ensure the TimeFunc is propagated to the validator
	ParseOptions []jwt.ParserOption

	// Default value is "exp"
	ExpField string

	// RefreshTokenTimeout specifies how long refresh tokens are valid
	// Defaults to 30 days if not set
	RefreshTokenTimeout time.Duration

	// RefreshTokenStore interface for storing and retrieving refresh tokens
	// If nil, an in-memory store will be used
	RefreshTokenStore core.TokenStore

	// RefreshTokenLength specifies the byte length of refresh tokens (default: 32)
	RefreshTokenLength int

	// UseRedisStore indicates whether to use Redis store instead of in-memory store
	// When true, will attempt to connect to Redis using RedisConfig
	UseRedisStore bool

	// RedisConfig configuration for Redis store when UseRedisStore is true
	// If nil when UseRedisStore is true, will use default Redis configuration
	RedisConfig *store.RedisConfig

	// inMemoryStore internal fallback refresh token store
	inMemoryStore *store.InMemoryRefreshTokenStore
}

var (
	// ErrMissingSecretKey indicates Secret key is required
	ErrMissingSecretKey = errors.New("secret key is required")

	// ErrForbidden when HTTP status 403 is given
	ErrForbidden = errors.New("you don't have permission to access this resource")

	// ErrMissingAuthenticatorFunc indicates Authenticator is required
	ErrMissingAuthenticatorFunc = errors.New("ginJWTMiddleware.Authenticator func is undefined")

	// ErrMissingLoginValues indicates a user tried to authenticate without username or password
	ErrMissingLoginValues = errors.New("missing Username or Password")

	// ErrFailedAuthentication indicates authentication failed, could be faulty username or password
	ErrFailedAuthentication = errors.New("incorrect Username or Password")

	// ErrFailedTokenCreation indicates JWT Token failed to create, reason unknown
	ErrFailedTokenCreation = errors.New("failed to create JWT Token")

	// ErrExpiredToken indicates JWT token has expired. Can't refresh.
	ErrExpiredToken = errors.New(
		"token is expired",
	) // in practice, this is generated from the jwt library not by us

	// ErrEmptyAuthHeader can be thrown if authing with a HTTP header, the Auth header needs to be set
	ErrEmptyAuthHeader = errors.New("auth header is empty")

	// ErrMissingExpField missing exp field in token
	ErrMissingExpField = errors.New("missing exp field")

	// ErrWrongFormatOfExp field must be float64 format
	ErrWrongFormatOfExp = errors.New("exp must be float64 format")

	// ErrInvalidAuthHeader indicates auth header is invalid, could for example have the wrong Realm name
	ErrInvalidAuthHeader = errors.New("auth header is invalid")

	// ErrEmptyQueryToken can be thrown if authing with URL Query, the query token variable is empty
	ErrEmptyQueryToken = errors.New("query token is empty")

	// ErrEmptyCookieToken can be thrown if authing with a cookie, the token cookie is empty
	ErrEmptyCookieToken = errors.New("cookie token is empty")

	// ErrEmptyParamToken can be thrown if authing with parameter in path, the parameter in path is empty
	ErrEmptyParamToken = errors.New("parameter token is empty")

	// ErrInvalidSigningAlgorithm indicates signing algorithm is invalid, needs to be HS256, HS384, HS512, RS256, RS384 or RS512
	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")

	// ErrNoPrivKeyFile indicates that the given private key is unreadable
	ErrNoPrivKeyFile = errors.New("private key file unreadable")

	// ErrNoPubKeyFile indicates that the given public key is unreadable
	ErrNoPubKeyFile = errors.New("public key file unreadable")

	// ErrInvalidPrivKey indicates that the given private key is invalid
	ErrInvalidPrivKey = errors.New("private key invalid")

	// ErrInvalidPubKey indicates the the given public key is invalid
	ErrInvalidPubKey = errors.New("public key invalid")

	// IdentityKey default identity key
	IdentityKey = "identity"

	// ErrMissingRefreshToken indicates the refresh token parameter is missing
	ErrMissingRefreshToken = errors.New("missing refresh_token parameter")

	// ErrInvalidRefreshToken indicates the refresh token is invalid or expired
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")

	// ErrRefreshTokenNotFound indicates the refresh token was not found in storage
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
)

// New creates and initializes a new GinJWTMiddleware instance
func New(m *GinJWTMiddleware) (*GinJWTMiddleware, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

func (mw *GinJWTMiddleware) readKeys() error { _ = "STUB: not implemented"; return nil }

func (mw *GinJWTMiddleware) privateKey() error { _ = "STUB: not implemented"; return nil }

// Log detailed error for debugging but don't expose to client

// Clear passphrase from memory immediately after use

func (mw *GinJWTMiddleware) publicKey() error { _ = "STUB: not implemented"; return nil }

// Log detailed error for debugging but don't expose to client

func (mw *GinJWTMiddleware) usingPublicKeyAlgo() bool { _ = "STUB: not implemented"; return false }

// MiddlewareInit initializes JWT middleware configuration with default values
func (mw *GinJWTMiddleware) MiddlewareInit() error { _ = "STUB: not implemented"; return nil }

// Initialize refresh token settings (RFC 6749 compliant by default)

// 30 days default

// 256 bits default

// Initialize in-memory store first (will be used as fallback)

// Try to initialize Redis store if enabled

// If Redis initialization didn't set a store, use in-memory

// bypass other key settings if KeyFunc is set

// MiddlewareFunc makes GinJWTMiddleware implement the Middleware interface.
func (mw *GinJWTMiddleware) MiddlewareFunc() gin.HandlerFunc {
	_ = "STUB: not implemented"
	return *new(gin.HandlerFunc)
}

func (mw *GinJWTMiddleware) middlewareImpl(c *gin.Context) { _ = "STUB: not implemented"; return }

// For backwards compatibility since technically exp is not required in the spec but has been in gin-jwt

// GetClaimsFromJWT get claims from JWT token
func (mw *GinJWTMiddleware) GetClaimsFromJWT(c *gin.Context) (jwt.MapClaims, error) {
	_ = "STUB: not implemented"
	return *new(jwt.MapClaims), nil
}

// Return the claims directly without unnecessary copying

// LoginHandler can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *GinJWTMiddleware) LoginHandler(c *gin.Context) { _ = "STUB: not implemented"; return }

// Generate complete token pair

// Set cookies

func (mw *GinJWTMiddleware) extractRefreshToken(c *gin.Context) string {
	_ = "STUB: not implemented"
	// Try to get refresh token from cookie first (most common for browser-based apps)
	return ""
}

// SECURITY: Query parameters are NOT supported for refresh tokens to prevent
// token leakage through server logs, proxy logs, browser history, and Referer headers.
// Only secure methods are supported: httpOnly cookies, request body (form/JSON).

// Check Content-Type to determine how to parse body
// This prevents consuming the body twice

// Try POST form

// Try JSON body

// LogoutHandler can be used by clients to remove the jwt cookie and revoke refresh token
func (mw *GinJWTMiddleware) LogoutHandler(c *gin.Context) {
	_ = "STUB: not implemented"
	// Extract JWT claims to make them available in LogoutResponse
	// This allows developers to access user information during logout
	return
}

// Handle refresh token revocation (RFC 6749 compliant)

// delete auth cookies

// Delete access token cookie

// Delete refresh token cookie

// Must match the secure flag used when setting the cookie

func (mw *GinJWTMiddleware) signedString(token *jwt.Token) (string, error) {
	_ = "STUB: not implemented"
	return "", nil
}

// generateRefreshToken creates a cryptographically secure random refresh token
func (mw *GinJWTMiddleware) generateRefreshToken() (string, error) {
	_ = "STUB: not implemented"
	return "", nil
}

// refreshTokenExpiry returns the effective expiry time for a refresh token,
// capped to min(RefreshTokenTimeout, MaxRefresh) when MaxRefresh is set.
func (mw *GinJWTMiddleware) refreshTokenExpiry(now time.Time) time.Time {
	_ = "STUB: not implemented"
	return *new(time.Time)
}

// storeRefreshToken stores a refresh token with user data.
func (mw *GinJWTMiddleware) storeRefreshToken(
	ctx context.Context,
	token string,
	userData any,
) error {
	_ = "STUB: not implemented"
	return nil
}

// validateRefreshToken validates a refresh token and returns associated user data
func (mw *GinJWTMiddleware) validateRefreshToken(ctx context.Context, token string) (any, error) {
	_ = "STUB: not implemented"
	return *new(any), nil
}

// revokeRefreshToken removes a refresh token from storage
func (mw *GinJWTMiddleware) revokeRefreshToken(ctx context.Context, token string) error {
	_ = "STUB: not implemented"
	return nil
}

// RefreshHandler can be used to refresh a token using RFC 6749 compliant refresh tokens.
// This handler expects a refresh_token parameter and returns a new access token and refresh token.
// Reply will be of the form {"access_token": "TOKEN", "refresh_token": "REFRESH_TOKEN"}.
func (mw *GinJWTMiddleware) RefreshHandler(c *gin.Context) {
	_ = "STUB: not implemented"
	// Extract refresh token from request
	return
}

// Validate refresh token

// Generate new token pair and revoke old refresh token

// Set cookies

// CheckIfTokenExpire check if token expire
func (mw *GinJWTMiddleware) CheckIfTokenExpire(c *gin.Context) (jwt.MapClaims, error) {
	_ = "STUB: not implemented"
	return *new(jwt.MapClaims), nil
}

// If we receive an error, and the error is anything other than a single
// ValidationErrorExpired, we want to return the error.
// If the error is just ValidationErrorExpired, we want to continue, as we can still
// refresh the token if it's within the MaxRefresh time.
// (see https://github.com/appleboy/gin-jwt/issues/176)

// generateAccessToken method that clients can use to get a jwt token.
func (mw *GinJWTMiddleware) generateAccessToken(data any) (string, time.Time, error) {
	_ = "STUB: not implemented"
	// 1. Validate signing algorithm
	return "", *new(time.Time), nil
}

// 2. Define framework-controlled claims that PayloadFunc cannot overwrite
// Only claims that the framework calculates/manages internally are reserved.
// Standard JWT claims (sub, iss, aud, nbf, iat, jti) are allowed to be set by users
// via PayloadFunc to comply with RFC 7519 best practices.

// Framework calculates expiration time
// Framework uses this for refresh mechanism

// 3. Safely add custom payload, avoiding framework-controlled field overwrites

// 4. Calculate expiration time using original data instead of claims

// 5. Set required system claims

// 6. Sign the token

// TokenGenerator generates a complete token pair (access + refresh) with RFC 6749 compliance
func (mw *GinJWTMiddleware) TokenGenerator(ctx context.Context, data any) (*core.Token, error) {
	_ = "STUB: not implemented"
	// Generate access token
	return nil, nil
}

// Generate refresh token

// Store refresh token

// TokenGeneratorWithRevocation generates a new token pair and revokes the old refresh token
func (mw *GinJWTMiddleware) TokenGeneratorWithRevocation(
	ctx context.Context,
	data any,
	oldRefreshToken string,
) (*core.Token, error) {
	_ = "STUB: not implemented"
	// Generate new token pair
	return nil, nil
}

// Revoke old refresh token, ignore if token already doesn't exist

func (mw *GinJWTMiddleware) jwtFromHeader(c *gin.Context, key string) (string, error) {
	_ = "STUB: not implemented"
	return "", nil
}

func (mw *GinJWTMiddleware) jwtFromQuery(c *gin.Context, key string) (string, error) {
	_ = "STUB: not implemented"
	return "", nil
}

func (mw *GinJWTMiddleware) jwtFromCookie(c *gin.Context, key string) (string, error) {
	_ = "STUB: not implemented"
	return "", nil
}

func (mw *GinJWTMiddleware) jwtFromParam(c *gin.Context, key string) (string, error) {
	_ = "STUB: not implemented"
	return "", nil
}

func (mw *GinJWTMiddleware) jwtFromForm(c *gin.Context, key string) (string, error) {
	_ = "STUB: not implemented"
	return "", nil
}

// ParseToken parse jwt token from gin context
func (mw *GinJWTMiddleware) ParseToken(c *gin.Context) (*jwt.Token, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

// save token string if valid

// save token string if valid

// ParseTokenString parse jwt token string
func (mw *GinJWTMiddleware) ParseTokenString(token string) (*jwt.Token, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

// unauthorized handles unauthorized requests by setting the WWW-Authenticate header
// and calling the user-defined Unauthorized callback.
//
// According to RFC 6750 (OAuth 2.0 Bearer Token Usage) and RFC 7235 (HTTP Authentication),
// a 401 Unauthorized response must include a WWW-Authenticate header with the Bearer scheme.
// This ensures compatibility with standard HTTP clients and authentication frameworks.
//
// See:
//   - https://tools.ietf.org/html/rfc6750
//   - https://tools.ietf.org/html/rfc7235
//   - https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401
func (mw *GinJWTMiddleware) unauthorized(c *gin.Context, code int, message string) {
	_ = "STUB: not implemented"
	return
}

// ExtractClaims help to extract the JWT claims
func ExtractClaims(c *gin.Context) jwt.MapClaims {
	_ = "STUB: not implemented"
	return *new(jwt.MapClaims)
}

// ExtractClaimsFromToken help to extract the JWT claims from token
func ExtractClaimsFromToken(token *jwt.Token) jwt.MapClaims {
	_ = "STUB: not implemented"
	return *new(jwt.MapClaims)
}

// Return the claims directly without unnecessary copying

// GetToken help to get the JWT token string
func GetToken(c *gin.Context) string { _ = "STUB: not implemented"; return "" }

// SetCookie help to set the token in the cookie
func (mw *GinJWTMiddleware) SetCookie(c *gin.Context, token string) {
	_ = "STUB: not implemented"
	// set cookie
	return
}

// SetRefreshTokenCookie help to set the refresh token in the cookie
func (mw *GinJWTMiddleware) SetRefreshTokenCookie(c *gin.Context, refreshToken string) {
	_ = "STUB: not implemented"
	return
}

// round up sub-second positive durations

// Always secure for refresh tokens (HTTPS only)
// Always httpOnly for security

// handleTokenError handles different types of JWT token validation errors
func (mw *GinJWTMiddleware) handleTokenError(c *gin.Context, err error) {
	_ = "STUB: not implemented"
	return
}

// generateTokenResponse creates a RFC 6749 compliant token response with refresh token
func (mw *GinJWTMiddleware) generateTokenResponse(_ *gin.Context, token *core.Token) gin.H {
	_ = "STUB: not implemented"
	return *new(gin.H)
}

// Include refresh token if present

// ClearSensitiveData clears sensitive data from memory
func (mw *GinJWTMiddleware) ClearSensitiveData() {
	_ = "STUB: not implemented"
	// Clear symmetric key
	return
}

// Clear private key bytes

// Clear public key bytes

// Clear passphrase

// Convert to []byte to clear, then back to string

// Note: RSA keys (mw.privKey, mw.pubKey) are harder to clear completely
// due to Go's garbage collector, but setting to nil helps

// Clear refresh token store if using in-memory store
