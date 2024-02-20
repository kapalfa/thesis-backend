package invitationsCRUD

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
)

func CreateInvitation(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Email string `json:"email"`
		Id    int    `json:"id"`
	}
	var input Req
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		log.Println(err)
		return
	}

	user := models.User{}
	database.DB.First(&user, "email = ?", input.Email)

	if user.Id == 0 {
		response := map[string]string{"message": "User not found"}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}
	projid := uint(input.Id)
	access := &models.Access{
		UserId:    user.Id,
		ProjectId: projid,
		Status:    "pending",
	}

	if err := database.DB.Create(access).Error; err != nil {
		if err.Error() == "UNIQUE constraint failed: accesses.user_id, accesses.project_id" {
			response := map[string]string{"message": "User already has access to project"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}
		response := map[string]string{"message": "Error creating access"}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}
}
