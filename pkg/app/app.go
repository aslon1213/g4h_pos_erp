package app

import (
	"aslon1213/magazin_pos/pkg/configs"
	"aslon1213/magazin_pos/platform/cache"
	"aslon1213/magazin_pos/platform/database"
	"aslon1213/magazin_pos/platform/logger"
	"log"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	fiberSwagger "github.com/swaggo/fiber-swagger"

	_ "aslon1213/magazin_pos/docs"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	Logger *zerolog.Logger
	Cache  *cache.Cache
	DB     *mongo.Client
	Config *configs.Config
	Router *fiber.App
}

func NewFiberApp() *fiber.App {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.CustomZerologMiddleware)
	app.Get("/docs/*", fiberSwagger.WrapHandler)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/docs/index.html")
	})
	return app
}

func NewApp() *App {
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	return &App{
		Logger: logger.SetupLogger(),
		Cache:  cache.New(),
		DB:     database.NewDB(),
		Router: NewFiberApp(),
		Config: config,
	}
}

func (a *App) Run() {
	controllers := NewControllers(a.DB.Database(a.Config.DB.Database), a.Cache)
	SetupRoutes(a.Router, controllers)
	a.Router.Listen(a.Config.Server.Port)
}
