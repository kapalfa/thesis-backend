package middleware 

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"context"
	"github.com/kapalfa/go/controllers/githubControllers"
)

func GithubCallbackHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type Code struct {
			Code string `json:"code"`
		}
		var code Code 
		err := json.NewDecoder(r.Body).Decode(&code)
		if err != nil {
			log.Printf("Error decoding response body: %v", err)
			return 
		}
		v := url.Values{}
		client_id := githubControllers.GetGithubClientID()
		client_secret := githubControllers.GetGithubClientSecret()
		url := "https://github.com/login/oauth/access_token"
		v.Set("code", code.Code)
		v.Set("client_id", client_id)
		v.Set("client_secret", client_secret)

		resp, err := http.PostForm(url, v)
		if err != nil {
			log.Printf("Error making request to github: %v", err)
			http.Error(w, "Error making request to github", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading body: %v", err)
			http.Error(w, "Error reading body", http.StatusInternalServerError)
			return
		}
		ctx := context.WithValue(r.Context(), "githubResponse", body)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
	//w.Write(body)
}