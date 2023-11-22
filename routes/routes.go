package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kapalfa/go/controllers"
	"github.com/kapalfa/go/middleware"
)

 func Setup(app *fiber.App) {
 	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)
	app.Get("/api/user", middleware.VerifyJWT(), controllers.User) //protected route
	app.Get("/api/logout", controllers.Logout)
	app.Get("/api/refresh", controllers.HandleRefreshToken)
	app.Get("/api/verify", middleware.VerifyJWT(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "verified",
		})
	})
}