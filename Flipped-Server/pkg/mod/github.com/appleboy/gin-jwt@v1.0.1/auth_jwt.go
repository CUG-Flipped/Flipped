package jwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gopkg.in/dgrijalva/jwt-go.v2"
	"log"
	"net/http"
	"strings"
	"time"
)

// GinJWTMiddleware provides a Json-Web-Token authentication implementation. On failure, a 401 HTTP response
// is returned. On success, the wrapped middleware is called, and the userId is made available as
// c.Get("userId").(string).
// Users can get a token by posting a json request to LoginHandler. The token then needs to be passed in
// the Authentication header. Example: Authorization:Bearer XXX_TOKEN_XXX#!/usr/bin/env
type GinJWTMiddleware struct {
	// Realm name to display to the user. Required.
	Realm string

	// signing algorithm - possible values are HS256, HS384, HS512
	// Optional, default is HS256.
	SigningAlgorithm string

	// Secret key used for signing. Required.
	Key []byte

	// Duration that a jwt token is valid. Optional, defaults to one hour.
	Timeout time.Duration

	// This field allows clients to refresh their token until MaxRefresh has passed.
	// Note that clients can refresh their token in the last moment of MaxRefresh.
	// This means that the maximum validity timespan for a token is MaxRefresh + Timeout.
	// Optional, defaults to 0 meaning not refreshable.
	MaxRefresh time.Duration

	// Callback function that should perform the authentication of the user based on userId and
	// password. Must return true on success, false on failure. Required.
	// Option return user id, if so, user id will be stored in Claim Array.
	Authenticator func(userId string, password string, c *gin.Context) (string, bool)

	// Callback function that should perform the authorization of the authenticated user. Called
	// only after an authentication success. Must return true on success, false on failure.
	// Optional, default to success.
	Authorizator func(userId string, c *gin.Context) bool

	// Callback function that will be called during login.
	// Using this function it is possible to add additional payload data to the webtoken.
	// The data is then made available during requests via c.Get("JWT_PAYLOAD").
	// Note that the payload is not encrypted.
	// The attributes mentioned on jwt.io can't be used as keys for the map.
	// Optional, by default no additional data will be set.
	PayloadFunc func(userId string) map[string]interface{}

	// User can define own Unauthorized func.
	Unauthorized func(*gin.Context, int, string)
}

// Login form structure.
type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// MiddlewareInit initialize jwt configs.
func (mw *GinJWTMiddleware) MiddlewareInit() error {
	if mw.Realm == "" {
		return errors.New("Realm is required")
	}

	if mw.Authenticator == nil {
		return errors.New("Authenticator is required")
	}

	if mw.Key == nil {
		return errors.New("Key is required")
	}

	if mw.SigningAlgorithm == "" {
		mw.SigningAlgorithm = "HS256"
	}

	if mw.Timeout == 0 {
		mw.Timeout = time.Hour
	}

	if mw.Authorizator == nil {
		mw.Authorizator = func(userId string, c *gin.Context) bool {
			return true
		}
	}

	if mw.Unauthorized == nil {
		mw.Unauthorized = func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		}
	}

	return nil
}

// MiddlewareFunc makes GinJWTMiddleware implement the Middleware interface.
func (mw *GinJWTMiddleware) MiddlewareFunc() gin.HandlerFunc {
	if err := mw.MiddlewareInit(); err != nil {
		log.Fatal(err.Error())
	}

	return func(c *gin.Context) {
		mw.middlewareImpl(c)
		return
	}
}

func (mw *GinJWTMiddleware) middlewareImpl(c *gin.Context) {
	token, err := mw.parseToken(c)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, err.Error())
		return
	}

	id := token.Claims["id"].(string)
	c.Set("JWT_PAYLOAD", token.Claims)
	c.Set("userID", id)

	if !mw.Authorizator(id, c) {
		mw.unauthorized(c, http.StatusForbidden, "You don't have permission to access.")
		return
	}

	c.Next()
}

// LoginHandler can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *GinJWTMiddleware) LoginHandler(c *gin.Context) {

	// Initial middleware default setting.
	mw.MiddlewareInit()

	var loginVals Login

	if c.BindJSON(&loginVals) != nil {
		mw.unauthorized(c, http.StatusBadRequest, "Missing Username or Password")
		return
	}

	userId, ok := mw.Authenticator(loginVals.Username, loginVals.Password, c)

	if !ok {
		mw.unauthorized(c, http.StatusUnauthorized, "Incorrect Username / Password")
		return
	}

	// Create the token
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(loginVals.Username) {
			token.Claims[key] = value
		}
	}

	if userId == "" {
		userId = loginVals.Username
	}

	expire := time.Now().Add(mw.Timeout)
	token.Claims["id"] = userId
	token.Claims["exp"] = expire.Unix()
	token.Claims["orig_iat"] = time.Now().Unix()

	tokenString, err := token.SignedString(mw.Key)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, "Create JWT Token faild")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339),
	})
}

// RefreshHandler can be used to refresh a token. The token still needs to be valid on refresh.
// Shall be put under an endpoint that is using the GinJWTMiddleware.
// Reply will be of the form {"token": "TOKEN"}.
func (mw *GinJWTMiddleware) RefreshHandler(c *gin.Context) {
	token, _ := mw.parseToken(c)

	origIat := int64(token.Claims["orig_iat"].(float64))

	if origIat < time.Now().Add(-mw.MaxRefresh).Unix() {
		mw.unauthorized(c, http.StatusUnauthorized, "Token is expired.")
		return
	}

	// Create the token
	newToken := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))

	for key := range token.Claims {
		newToken.Claims[key] = token.Claims[key]
	}

	expire := time.Now().Add(mw.Timeout)
	newToken.Claims["id"] = token.Claims["id"]
	newToken.Claims["exp"] = expire.Unix()
	newToken.Claims["orig_iat"] = origIat

	tokenString, err := newToken.SignedString(mw.Key)

	if err != nil {
		mw.unauthorized(c, http.StatusUnauthorized, "Create JWT Token faild")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  tokenString,
		"expire": expire.Format(time.RFC3339),
	})
}

// ExtractClaims help to extract the JWT claims
func ExtractClaims(c *gin.Context) map[string]interface{} {

	if _, exists := c.Get("JWT_PAYLOAD"); !exists {
		emptyClaims := make(map[string]interface{})
		return emptyClaims
	}

	jwtClaims, _ := c.Get("JWT_PAYLOAD")

	return jwtClaims.(map[string]interface{})
}

// TokenGenerator handler that clients can use to get a jwt token.
func (mw *GinJWTMiddleware) TokenGenerator(userID string) string {
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(userID) {
			token.Claims[key] = value
		}
	}

	token.Claims["id"] = userID
	token.Claims["exp"] = time.Now().Add(mw.Timeout).Unix()
	token.Claims["orig_iat"] = time.Now().Unix()

	tokenString, _ := token.SignedString(mw.Key)

	return tokenString
}

func (mw *GinJWTMiddleware) parseToken(c *gin.Context) (*jwt.Token, error) {
	authHeader := c.Request.Header.Get("Authorization")

	if authHeader == "" {
		return nil, errors.New("Auth header empty")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return nil, errors.New("Invalid auth header")
	}

	return jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(mw.SigningAlgorithm) != token.Method {
			return nil, errors.New("Invalid signing algorithm")
		}

		return mw.Key, nil
	})
}

func (mw *GinJWTMiddleware) unauthorized(c *gin.Context, code int, message string) {
	c.Header("WWW-Authenticate", "JWT realm="+mw.Realm)
	c.Abort()

	mw.Unauthorized(c, code, message)

	return
}
