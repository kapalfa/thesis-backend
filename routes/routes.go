package routes

import (
 	"github.com/gorilla/mux"
 	"github.com/kapalfa/go/middleware"
	"github.com/kapalfa/go/config"
    "github.com/kapalfa/go/controllers/projectsCRUD"
 	"github.com/kapalfa/go/controllers/filesCRUD"
 	"github.com/kapalfa/go/controllers/authControllers"
 	"github.com/kapalfa/go/controllers/githubControllers"
	"github.com/kapalfa/go/controllers/invitationsCRUD"
	"net/http"
)
// 	app.Get("/api/user", middleware.VerifyJWT(), controllers.User) //protected route
// 	app.Get("/api/verify", middleware.VerifyJWT(), func(c *fiber.Ctx) error {
// 		return c.JSON(fiber.Map{
// 			"success": true,
// 			"message": "verified",
// 		})
// 	})// 	app.Get("/api/getFiles/:id", middleware.VerifyJWT(), middleware.VerifyAccess(), filesCRUD.GetFiles)
func Setup(r *mux.Router) {
	// file routes
	r.HandleFunc("/api/getFile/{filepath:.*}", filesCRUD.GetFile)
	r.HandleFunc("/api/getFiles/{id}", filesCRUD.GetFiles)
	r.HandleFunc("/api/saveFile/{filepath:.*}", filesCRUD.SaveFile)
	r.HandleFunc("/api/upload/{filepath:.*}", filesCRUD.UploadFile)
	r.HandleFunc("/api/uploadFolder/{filepath:.*}", filesCRUD.UploadFolder)
	r.HandleFunc("/api/createFile/{filepath:.*}", filesCRUD.CreateFile)
	r.HandleFunc("/api/createFolder/{folderpath:.*}", filesCRUD.CreateFolder)
	r.HandleFunc("/api/deleteFile/{filepath:.*}", filesCRUD.DeleteFile)
	r.HandleFunc("/api/deleteFolder/{folderpath:.*}", filesCRUD.DeleteFolder)
	// user routes
	r.HandleFunc("/api/login", authControllers.Login)
	r.HandleFunc("/api/logout", authControllers.Logout)
	r.HandleFunc("/api/register", authControllers.Register)
	r.HandleFunc("/api/refresh", authControllers.HandleRefreshToken)
	// project routes
	r.HandleFunc("/api/getProject/{id}", projectsCRUD.GetProject)
	r.HandleFunc("/api/getProjects/{userid}", projectsCRUD.GetProjects)
	r.HandleFunc("/api/createProject", projectsCRUD.CreateProject)
	r.HandleFunc("/api/deleteProject/{id}", projectsCRUD.DeleteProject)
	r.HandleFunc("/api/searchProjects/{projectName}", projectsCRUD.SearchProjects)
	r.HandleFunc("/api/getPublicProjects", projectsCRUD.GetPublicProjects)
	r.HandleFunc("/api/copyProject", projectsCRUD.CopyProject)
	r.HandleFunc("/api/getCollaborators/{id}/{userid}", projectsCRUD.GetCollaborators)
	// github routes
	r.HandleFunc("/github/login", githubControllers.GithubLoginHandler)
	// r.HandleFunc("/github/callback", authControllers.LoginGithub).Use(middleware.GithubCallbackHandler)
	r.Handle("/github/callback", middleware.GithubCallbackHandler(http.HandlerFunc(authControllers.LoginGithub)))
	r.HandleFunc("/github/initRepo", githubControllers.InitGitRepo)
	r.HandleFunc("/github/commitRepo", githubControllers.CommitGitRepo)
	r.HandleFunc("/github/downloadRepo/{userid}", githubControllers.DownloadGitRepo)
	// socket routes
	r.HandleFunc("/sockets/projId={projId:[0-9]+}&userId={userId:[0-9]+}", config.HandleWebsocketConnection)
	// invitation routes
	r.HandleFunc("/api/createInvitation", invitationsCRUD.CreateInvitation)
	r.HandleFunc("/api/getInvitations/{userid}", invitationsCRUD.GetInvitations)
	r.HandleFunc("/api/handleInvitation/{projectid}/{userid}", invitationsCRUD.HandleInvitation)
}