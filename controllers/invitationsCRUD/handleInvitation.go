package invitationsCRUD

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
)

type Req struct {
	Response string `json:"response"`
}
func HandleInvitation(w http.ResponseWriter, r *http.Request) {
	var req Req 
	vars := mux.Vars(r)
	projectid := vars["projectid"]
	userid := vars["userid"]

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println(err)
		return
	}

	if req.Response == "yes" {
		database.DB.Model(&models.Access{}).Where("user_id = ? AND project_id = ?", userid, projectid).Update("status", "accepted")
	} else {
		database.DB.Model(&models.Access{}).Where("user_id = ? AND project_id = ?", userid, projectid).Delete(&models.Access{})
	}
}