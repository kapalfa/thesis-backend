package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kapalfa/go/config"
	"github.com/kapalfa/go/controllers/authControllers"
	"github.com/kapalfa/go/controllers/filesCRUD"
	"github.com/kapalfa/go/controllers/githubControllers"
	"github.com/kapalfa/go/controllers/invitationsCRUD"
	"github.com/kapalfa/go/controllers/projectsCRUD"
	"github.com/kapalfa/go/middleware"
)

// app.Get("/api/getFiles/:id", middleware.VerifyJWT(), middleware.VerifyAccess(), filesCRUD.GetFiles)
func Setup(r *mux.Router) {
	// file routes
	r.HandleFunc("/getFile/{filepath:.*}", filesCRUD.GetFile)
	r.HandleFunc("/getFiles/{id}", filesCRUD.GetFiles)
	r.HandleFunc("/saveFile/{filepath:.*}", filesCRUD.SaveFile)
	r.HandleFunc("/upload/{filepath:.*}", filesCRUD.UploadFile)
	r.HandleFunc("/uploadFolder/{filepath:.*}", filesCRUD.UploadFolder)
	r.HandleFunc("/createFile/{filepath:.*}", filesCRUD.CreateFile)
	r.HandleFunc("/createFolder/{folderpath:.*}", filesCRUD.CreateFolder)
	r.HandleFunc("/deleteFile/{filepath:.*}", filesCRUD.DeleteFile)
	r.HandleFunc("/deleteFolder/{folderpath:.*}", filesCRUD.DeleteFolder)
	// user routes
	r.HandleFunc("/login", authControllers.Login)
	r.HandleFunc("/logout", authControllers.Logout)
	r.HandleFunc("/register", authControllers.Register)
	r.HandleFunc("/refresh", authControllers.HandleRefreshToken)
	// project routes
	r.HandleFunc("/getProject/{id}", projectsCRUD.GetProject)
	r.HandleFunc("/getProjects/{userid}", projectsCRUD.GetProjects)
	r.HandleFunc("/createProject", projectsCRUD.CreateProject)
	r.HandleFunc("/deleteProject/{id}", projectsCRUD.DeleteProject)
	r.HandleFunc("/searchProjects/{projectName}", projectsCRUD.SearchProjects)
	r.HandleFunc("/getPublicProjects", projectsCRUD.GetPublicProjects)
	r.HandleFunc("/copyProject", projectsCRUD.CopyProject)
	r.HandleFunc("/getCollaborators/{id}/{userid}", projectsCRUD.GetCollaborators)
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
	r.HandleFunc("/createInvitation", invitationsCRUD.CreateInvitation)
	r.HandleFunc("/getInvitations", invitationsCRUD.GetInvitations)
	r.HandleFunc("/handleInvitation/{projectid}", invitationsCRUD.HandleInvitation)

	r.HandleFunc("/confirmEmail", authControllers.VerifyMail)
	r.HandleFunc("/forgotPassword", authControllers.ForgotPassword)
	r.HandleFunc("/setNewPassword", authControllers.SetNewPassword)
}
