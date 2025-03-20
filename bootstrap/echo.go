package bootstrap

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewEcho() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())   // Logging requests
	e.Use(middleware.Recover())  // Recover from panics
	e.Use(middleware.CORS())     // Enable CORS
	e.Use(middleware.Secure())   // Security headers (Helmet equivalent)
	e.Use(middleware.BodyLimit("10M")) // Set body size limit

	return e
}
