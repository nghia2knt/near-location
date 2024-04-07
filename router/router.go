package router

import (
	"near-location/internal/controller"
	"near-location/pkg/util"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func NewRouter(controller *controller.Controller) {
	app := fiber.New(fiber.Config{
		ErrorHandler: util.ErrorHandler,
	})
	app.Use(logger.New())
	app.Get("/locations", controller.GetLocations)
	app.Listen(":3000")
}
