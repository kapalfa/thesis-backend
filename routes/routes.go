package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kapalfa/go/controllers"
	"github.com/kapalfa/go/middleware"
	"github.com/kapalfa/go/controllers/projectsCRUD"
	"github.com/kapalfa/go/controllers/filesCRUD"
	"github.com/kapalfa/go/controllers/authControllers"
	"github.com/kapalfa/go/controllers/uploads"
)

 func Setup(app *fiber.App) {
 	app.Post("/api/register", authControllers.Register)
	app.Post("/api/login", authControllers.Login)
	app.Get("/api/user", middleware.VerifyJWT(), controllers.User) //protected route
	app.Get("/api/logout", authControllers.Logout)
	app.Get("/api/refresh", authControllers.HandleRefreshToken)
	app.Get("/api/verify", middleware.VerifyJWT(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"message": "verified",
		})
	})
	app.Post("/api/upload/*", uploads.FileUpload)
	app.Post("/api/uploadFolder", uploads.FolderUpload)

	app.Post("/api/createProject", projectsCRUD.CreateProject)
	app.Get("/api/getProject/:id", projectsCRUD.GetProject)
	app.Get("/api/getProjects/:userid", projectsCRUD.GetProjects)
	app.Get("/api/searchProjects/:projectName", projectsCRUD.SearchProjects)

	app.Get("/api/getFiles/:id", middleware.VerifyJWT(), middleware.VerifyAccess(), filesCRUD.GetFiles)
	app.Get("/api/getFile/*", filesCRUD.GetFile)
	app.Post("/api/saveFile/*", filesCRUD.SaveFile)
}