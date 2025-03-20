package middleware

import (
	"encoding/json"
	"my-graphql-project/bootstrap"
	"strings"
	"time"

	"github.com/kataras/jwt"
	"github.com/labstack/echo/v4"
)

type ClaimsToken struct {
	Id        string `json:"id,omitempty"`
	Role      string `json:"role,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
}

// var (
// 	JWT_ACCESS_TOKEN  = []byte("superdupersecret")
// 	JWT_REFRESH_TOKEN = []byte("superdupersecretrefresher")
// )

type TokenPair struct {
	AccessToken  json.RawMessage `json:"access_token,omitempty"`
	RefreshToken json.RawMessage `json:"refresh_token,omitempty"`
}

func GenerateJWTToken(id string, role string) (*TokenPair, error) {
	standardClaims := ClaimsToken{
		Id:        id,
		Role:      role,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 300).Unix(),
	}
	encrypt, _, err := jwt.GCM([]byte(bootstrap.GlobalEnv.JWT.AccessToken), nil)
	if err != nil {

		return nil, err
	}
	token, err := jwt.SignEncrypted(jwt.HS256, []byte(bootstrap.GlobalEnv.JWT.AccessToken), encrypt, standardClaims, jwt.MaxAge(time.Hour*24*7))
	if err != nil {

		return nil, err
	}
	reEncrypt, _, _ := jwt.GCM([]byte(bootstrap.GlobalEnv.JWT.RefreshToken), nil)
	refreshToken, err := jwt.SignEncrypted(jwt.HS256, []byte(bootstrap.GlobalEnv.JWT.RefreshToken), reEncrypt, standardClaims, jwt.MaxAge(time.Hour*24*8))
	if err != nil {

		return nil, err
	}
	tokenPairData := jwt.NewTokenPair(token, refreshToken)
	return &TokenPair{
		AccessToken:  BytesQuote(tokenPairData.AccessToken),
		RefreshToken: BytesQuote(tokenPairData.RefreshToken),
	}, nil
}

// AccessTokenMiddleware validates the JWT token using Kataras JWT
func AccessTokenMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(echo.ErrUnauthorized.Code, echo.Map{
				"status": false,
				"result": "Missing or invalid Authorization header",
			})
		}

		// Extract the token
		tokenString := strings.TrimSpace(authHeader[7:])

		// Decrypt and verify the token
		_, decrypt, _ := jwt.GCM([]byte(bootstrap.GlobalEnv.JWT.AccessToken), nil)
		verifiedToken, err := jwt.VerifyEncrypted(jwt.HS256, []byte(bootstrap.GlobalEnv.JWT.AccessToken), decrypt, []byte(tokenString))
		if err != nil {
			return c.JSON(echo.ErrUnauthorized.Code, echo.Map{
				"status": false,
				"result": "Invalid or expired token",
			})
		}

		// Extract claims
		var claims ClaimsToken
		err = verifiedToken.Claims(&claims)
		if err != nil {
			return c.JSON(echo.ErrUnauthorized.Code, echo.Map{
				"status": false,
				"result": "Invalid token claims",
			})
		}

		// Store claims in context for later use
		c.Set("user", claims)

		// Continue to the next handler
		return next(c)
	}
}
func GetOwnerAccessToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Extract the Authorization header
		auth := c.Request().Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return echo.ErrUnauthorized
		}

		// Remove "Bearer " prefix
		jwtFromHeader := strings.TrimSpace(auth[7:])

		// Decrypt and verify the token
		_, decrypt, _ := jwt.GCM([]byte(bootstrap.GlobalEnv.JWT.AccessToken), nil)
		verifiedToken, err := jwt.VerifyEncrypted(jwt.HS256, []byte(bootstrap.GlobalEnv.JWT.AccessToken), decrypt, []byte(jwtFromHeader))
		if err != nil {
			return echo.ErrUnauthorized
		}

		// Extract claims
		var claims ClaimsToken
		err = verifiedToken.Claims(&claims)
		if err != nil {
			return echo.ErrUnauthorized
		}

		// Store claims in context
		c.Set("owner_id", claims.Id)

		// Continue to the next handler
		return next(c)
	}
}

// GetInfoAccessToken extracts JWT claims from the access token
func GetInfoAccessToken(c echo.Context) (*ClaimsToken, error) {
	// Get the Authorization header
	auth := c.Request().Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		return nil, echo.NewHTTPError(401, "Missing or invalid Authorization header") // Use 401 directly
	}

	// Remove "Bearer " prefix
	jwtFromHeader := strings.TrimSpace(auth[7:])

	// Decrypt and verify the token using the global secret
	_, decrypt, _ := jwt.GCM([]byte(bootstrap.GlobalEnv.JWT.AccessToken), nil)
	verifiedToken, err := jwt.VerifyEncrypted(jwt.HS256, []byte(bootstrap.GlobalEnv.JWT.AccessToken), decrypt, []byte(jwtFromHeader))
	if err != nil {
		return nil, err
	}

	// Extract the claims from the verified token
	var claims ClaimsToken
	err = verifiedToken.Claims(&claims)
	if err != nil {
		return nil, err
	}

	// Return the claims
	return &claims, nil
}

// GetOwnerRefresh retrieves the refresh token's claims
func GetOwnerRefresh(c echo.Context) (*ClaimsToken, error) {
	// Get the Authorization header
	auth := c.Request().Header.Get("Authorization")
	if auth == "" {
		return nil, echo.NewHTTPError(401, "Missing Authorization header") // Unauthorized error
	}

	// Remove "Bearer " prefix from the token
	jwtFromHeader := strings.TrimSpace(auth[7:])

	// Decrypt and verify the token using the global secret
	_, decrypt, _ := jwt.GCM([]byte(bootstrap.GlobalEnv.JWT.RefreshToken), nil)
	verifiedToken, err := jwt.VerifyEncrypted(jwt.HS256, []byte(bootstrap.GlobalEnv.JWT.RefreshToken), decrypt, []byte(jwtFromHeader))
	if err != nil {
		return nil, err
	}

	// Extract the claims from the verified token
	var claims ClaimsToken
	err = verifiedToken.Claims(&claims)
	if err != nil {
		return nil, err
	}

	// Return only the Id from the claims
	return &ClaimsToken{
		Id: claims.Id,
	}, nil
}

// GenerateRefreshToken generates a new access and refresh token pair
func GenerateRefreshToken(c echo.Context) (*TokenPair, error) {
	// Get the Authorization header
	auth := c.Request().Header.Get("Authorization")
	if auth == "" {
		return nil, echo.NewHTTPError(401, "Missing Authorization header")
	}

	// Retrieve claims
	claims, err := GetOwnerRefresh(c)
	if err != nil {
		return nil, err
	}

	// Prepare standard claims
	standardClaims := ClaimsToken{
		Id:        claims.Id,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	// Generate the access token
	encrypt, _, _ := jwt.GCM([]byte(bootstrap.GlobalEnv.JWT.AccessToken), nil)
	token, err := jwt.SignEncrypted(jwt.HS256, []byte(bootstrap.GlobalEnv.JWT.AccessToken), encrypt, standardClaims, jwt.MaxAge(time.Second*15))
	if err != nil {
		return nil, err
	}

	// Generate the refresh token
	reEncrypt, _, _ := jwt.GCM([]byte(bootstrap.GlobalEnv.JWT.RefreshToken), nil)
	refreshToken, err := jwt.SignEncrypted(jwt.HS256, []byte(bootstrap.GlobalEnv.JWT.RefreshToken), reEncrypt, standardClaims, jwt.MaxAge(time.Hour*24*8))
	if err != nil {
		return nil, err
	}

	// Create and return token pair
	tokenPairData := jwt.NewTokenPair(token, refreshToken)
	return &TokenPair{
		AccessToken:  tokenPairData.AccessToken,
		RefreshToken: tokenPairData.RefreshToken,
	}, nil
}

func BytesQuote(b []byte) []byte {
	dst := make([]byte, len(b)+2)
	dst[0] = '"'
	copy(dst[1:], b)
	dst[len(dst)-1] = '"'
	return dst
}
