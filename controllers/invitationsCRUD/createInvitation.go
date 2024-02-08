package invitationsCRUD 

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
)

func CreateInvitation(w http.ResponseWriter, r* http.Request) {
	type Req struct	{
		Email string `json:"email"`
		Id int `json:"id"`
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
		log.Println("User not found")
		return 
	}
	//i, err := strconv.ParseUint(input.ProjectId, 10, 32)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	projid := uint(input.Id)
	access := &models.Access{
		UserId:    user.Id,
		ProjectId: projid,
		Status:    "pending",
	}

	if err := database.DB.Create(access).Error; err != nil {
		log.Println("Could not create access on project")
		return
	}
}