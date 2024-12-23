package main

import (
	"QuickPicsAuth/database"
	"QuickPicsAuth/routes"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	_, err := database.ConnectDB()
	if err != nil {
		panic("Could not connect to the database")
	}

	fmt.Println("Connected to the database")

	app := fiber.New()

	// app.Use(csrf.New(csrf.Config{
	// 	KeyLookup:         "header:X-Csrf-Token",
	// 	CookieName:        "_Host-csrf",
	// 	CookieSameSite:    "Lax",
	// 	CookieSecure:      true,
	// 	CookieSessionOnly: true,
	// 	CookieHTTPOnly:    true,
	// 	Expiration:        1 * time.Hour,
	// 	KeyGenerator:      utils.UUIDv4,
	// 	// ErrorHandler:      defaultErrorHandler,
	// 	// Extractor:         CsrfFromHeader(HeaderName),
	// 	Session:           session.New(),
	// 	SessionKey:        "fiber.csrf.token",
	// 	HandlerContextKey: "fiber.csrf.handler",
	// }))

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization,Accept,Origin,Access-Control-Request-Method,Access-Control-Request-Headers,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Access-Control-Allow-Methods,Access-Control-Expose-Headers,Access-Control-Max-Age,Access-Control-Allow-Credentials",
		AllowCredentials: true,
	}))

	// env := godotenv.Load()
	// if env != nil {
	// 	panic("Error loading .env file")
	// }

	routes.SetUpRoutes(app)

	err = app.Listen(":3000")
	if err != nil {
		panic("Could not start the server")
	}
}
