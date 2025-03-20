package bootstrap

import (
	"github.com/gofiber/fiber/v2"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Application struct {
	Env   *Env
	Fiber *fiber.App
	Echo  *echo.Echo
	DB    *gorm.DB
}

var GlobalEnv Env

func App() *Application {
	app := &Application{}
	app.Env = NewEnv()
	GlobalEnv = *NewEnv()
	app.Echo = NewEcho()
	app.DB = NewPostgresConnection(app.Env)
	// app.DB = NewPostgresConnecttion(app.Env)
	return app
}
