package routes

import (
	"QuickPicsAuth/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App) {
	app.Get("/", controllers.Hello)
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Post("/api/forgot-password", controllers.ForgotPassword)
	app.Get("/api/user", controllers.User)
	app.Post("/api/logout", controllers.Logout)
}
