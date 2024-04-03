package authControllers

import (
	"encoding/json"
	"net/http"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/config"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"log"
	"net/url"
)

func LoginGithub(w http.ResponseWriter, r *http.Request) {
	githubToken := r.Context().Value("githubResponse").([]byte)
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		log.Printf("Error creating request to fetch user email: %v", err)
		return 
	}
	values, err := url.ParseQuery(string(githubToken))
	if err != nil {
		log.Printf("Error parsing github token: %v", err)
		return
	}
	access := values.Get("access_token")
	if access == "" {
		log.Printf("Error getting access token")
		return
	}
	req.Header.Set("Authorization", "Bearer " + string(access))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error making request to fetch user email: %v", err)
		return
	}
	defer resp.Body.Close()
	
	type Emails struct {
		Email string `json:"email"`
		Verified bool `json:"verified"`
		Primary bool `json:"primary"`
		Visibility string `json:"visibility"`
	}

	var emails []Emails
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		log.Printf("Error decoding response body: %v", err)
		return
	}
	var primaryEmail string
	for _, email := range emails {
		if email.Primary && email.Verified {
			primaryEmail = email.Email
			break
		}
	}
	if primaryEmail == "" {
		log.Printf("No primary email found")
		return
	}	
	userModel, err := getUserByEmail(primaryEmail)
	if userModel == nil {
	 	response := map[string]interface{}{
	 		"status": false,
	 		"message": "User not found",
	 	}
	 	json.NewEncoder(w).Encode(response)
	 	return
	} 

	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["id"] = userModel.Id
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token, err := accessToken.SignedString([]byte(config.Config("ACCESS_TOKEN_SECRET")))
	if err != nil {
	 	http.Error(w, err.Error(), http.StatusInternalServerError)
	 	return
	}

	resfreshToken := jwt.New(jwt.SigningMethodHS256)
	claims = resfreshToken.Claims.(jwt.MapClaims)
	claims["id"] = userModel.Id
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	rt, err := resfreshToken.SignedString([]byte(config.Config("REFRESH_TOKEN_SECRET")))
	if err != nil {
	 	http.Error(w, err.Error(), http.StatusInternalServerError)
	 	return
	}

    userModel.RefreshToken = rt
	userModel.GithubToken = string(access)
	database.DB.Save(&userModel)

	cookie := http.Cookie{
		Name: "jwt",
	 	Value: rt,
	 	Expires: time.Now().Add(time.Hour),
	 	HttpOnly: true,
	 	Path: "/",
	 	SameSite: http.SameSiteNoneMode,
	 	Secure: true,
	}

	http.SetCookie(w, &cookie)
	response := map[string]interface{}{
	 	"status": true,
	 	"message": "Logged in",
	 	"access_token": token,
	 	"cookie": cookie,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
	 	http.Error(w, err.Error(), http.StatusInternalServerError)
	 	return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
