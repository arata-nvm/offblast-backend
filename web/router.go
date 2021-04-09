package web

import (
	"log"

	"github.com/arata-nvm/offblast/config"
	"github.com/arata-nvm/offblast/web/handler"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewServer() *echo.Echo {
	initSentry()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.Use(sentryecho.New(sentryecho.Options{}))

	v1 := e.Group("/api/v1")
	v1.POST("/detect", handler.Detect)
	v1.GET("/random", handler.Random)

	return e
}

func initSentry() {
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: config.SentryDsn(),
	}); err != nil {
		log.Fatalln(err)
	}
}
