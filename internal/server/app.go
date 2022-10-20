package server

import (
	"log"
	"net/http"

	"github.com/internal/server/nrfiber"
	"github.com/newrelic/go-agent/v3/newrelic"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
)

type App struct {
	*fiber.App
	config Config
}

func (app *App) Start(addr string) error {
	return app.Listen(addr)
}

func New(config ...Config) *App {
	app := &App{
		App: fiber.New(fiber.Config{
			DisableStartupMessage: true,
		}),
		config: Config{
			Recovery:  true,
			Swagger:   true,
			RequestID: true,
			Logger:    true,
		},
	}

	nrapp, err := newrelic.NewApplication(
		newrelic.ConfigAppName("pets_users-api"),
		newrelic.ConfigLicense("2fe9cb4919c003a782ef3028322128dda528NRAL"),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)

	if err != nil {
		log.Fatal(err)
	}

	app.Use(nrfiber.New(nrfiber.Config{
		NewRelicApp: nrapp,
	}))

	if len(config) > 0 {
		app.config = config[0]
	}

	if app.config.Recovery {
		app.Use(recover.New(recover.Config{
			EnableStackTrace: true,
		}))
	}

	if app.config.RequestID {
		app.Use(requestid.New())
	}

	if app.config.Logger {
		app.Use(logger.New(logger.Config{
			Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}\n",
		}))
	}

	if app.config.Swagger {
		app.Add(http.MethodGet, "/swagger/*", swagger.HandlerDefault)
	}

	return app
}

type Config struct {
	Recovery  bool
	Swagger   bool
	RequestID bool
	Logger    bool
}
