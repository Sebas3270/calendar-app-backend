package main

import (
	"os"

	"github.com/Sebas3270/calendar-app-backend/middlewares"
	"github.com/Sebas3270/calendar-app-backend/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		print("Error loading .env file")
	}

	app := fiber.New()

	app.Use(cors.New())

	/*Static files*/
	app.Static("/", "./public")

	/* Api Routes */
	api := app.Group("/api")
	authRoute := api.Group("/auth")
	eventRoute := api.Group("/events")

	/*Auth Routes*/
	authRoute.Post("/new", services.CreateUser)
	authRoute.Post("/login", services.Login)
	authRoute.Post("/renew", middlewares.ValidateJwt, services.RenewToken)

	/*Event routes*/
	eventRoute.Get("/", middlewares.ValidateJwt, services.GetEvents)
	eventRoute.Post("/", middlewares.ValidateJwt, services.CreateEvent)
	eventRoute.Put("/:id", middlewares.ValidateJwt, services.UpdateEvent)
	eventRoute.Delete("/:id", middlewares.ValidateJwt, services.DeleteEvent)

	/*Static files in any other specified route*/
	app.Static("/*", "./public")

	app.Listen(":" + os.Getenv("PORT"))

}
