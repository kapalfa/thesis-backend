package routes

import (
 	"github.com/gorilla/mux"
// 	"github.com/kapalfa/go/middleware"
	"github.com/kapalfa/go/config"
    "github.com/kapalfa/go/controllers/projectsCRUD"
 	"github.com/kapalfa/go/controllers/filesCRUD"
 	"github.com/kapalfa/go/controllers/authControllers"
 	"github.com/kapalfa/go/controllers/githubControllers"
)

// 	app.Get("/api/user", middleware.VerifyJWT(), controllers.User) //protected route
// 	app.Get("/api/verify", middleware.VerifyJWT(), func(c *fiber.Ctx) error {
// 		return c.JSON(fiber.Map{
// 			"success": true,
// 			"message": "verified",
// 		})
// 	})
// 	app.Post("/api/createProject", projectsCRUD.CreateProject)
// 	app.Get("/api/searchProjects/:projectName", projectsCRUD.SearchProjects)
// 	app.Get("/api/getFiles/:id", middleware.VerifyJWT(), middleware.VerifyAccess(), filesCRUD.GetFiles)

func Setup(r *mux.Router) {
	// file routes
	r.HandleFunc("/api/getFile/{filepath:.*}", filesCRUD.GetFile)
	r.HandleFunc("/api/getFiles/{id}", filesCRUD.GetFiles)
	r.HandleFunc("/api/saveFile/{filepath:.*}", filesCRUD.SaveFile)
	r.HandleFunc("/api/upload/{filepath:.*}", filesCRUD.UploadFile)
	r.HandleFunc("/api/uploadFolder/{filepath:.*}", filesCRUD.UploadFolder)
	r.HandleFunc("/api/createFile/{filepath:.*}", filesCRUD.CreateFile)
	r.HandleFunc("/api/deleteFile/{filepath:.*}", filesCRUD.DeleteFile)
	// user routes
	r.HandleFunc("/api/login", authControllers.Login)
	r.HandleFunc("/api/logout", authControllers.Logout)
	r.HandleFunc("/api/register", authControllers.Register)
	r.HandleFunc("/api/refresh", authControllers.HandleRefreshToken)
	r.HandleFunc("/api/githubLogin", authControllers.LoginGithub)
	// project routes
	r.HandleFunc("/api/getProject/{id}", projectsCRUD.GetProject)
	r.HandleFunc("/api/getProjects/{userid}", projectsCRUD.GetProjects)
	r.HandleFunc("/api/createProject", projectsCRUD.CreateProject)
	r.HandleFunc("/api/deleteProject/{id}", projectsCRUD.DeleteProject)
	r.HandleFunc("/api/searchProjects/{projectName}", projectsCRUD.SearchProjects)
	r.HandleFunc("/api/getPublicProjects", projectsCRUD.GetPublicProjects)
	r.HandleFunc("/api/copyProject", projectsCRUD.CopyProject)
	// github routes
	r.HandleFunc("/github/login", githubControllers.GithubLoginHandler)
	r.HandleFunc("/github/callback", githubControllers.GithubCallbackHandler)
	// socket routes
	r.HandleFunc("/sockets", config.HandleWebsocketConnection)

}