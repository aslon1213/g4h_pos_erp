package app

import (
	"aslon1213/magazin_pos/pkg/configs"
	"aslon1213/magazin_pos/platform/cache"
	"aslon1213/magazin_pos/platform/database"
	"aslon1213/magazin_pos/platform/logger"
	"log"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	Logger *zerolog.Logger
	Redis  *redis.Client
	DB     *mongo.Client
	Config *configs.Config
	Router *fiber.App
}

func NewFiberApp() *fiber.App {
	app := fiber.New()

	app.Use(logger.CustomZerologMiddleware)
	return app
}

func NewApp() *App {
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	return &App{
		Logger: logger.SetupLogger(),
		Redis:  cache.NewCache(),
		DB:     database.NewDB(),
		Router: NewFiberApp(),
		Config: config,
	}
}

func (a *App) Run() {
	a.Router.Listen(a.Config.Server.Port)
}
