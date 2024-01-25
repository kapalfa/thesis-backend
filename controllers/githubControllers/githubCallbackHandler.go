package githubControllers

import (
	"io"
	"log"
	"fmt"
	"net/http"
	"net/url"
)
func GithubCallbackHandler(w http.ResponseWriter, r *http.Request) {
// 	//state := r.FormValue("state")
// 	//if state != oauthStateString {
// 	//	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
// 	//	return
// 	//}
	v := url.Values{}
   	code := r.URL.Query().Get("code")
	client_id := GetGithubClientID()
	client_secret := GetGithubClientSecret()
	url := "https://github.com/login/oauth/access_token"
	v.Set("code", code)
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
	fmt.Println(string(body))

	w.Write(body)
}
// //	_, err := oauthConf.Exchange(oauth2.NoContext, code)
// 	// if err != nil {
// 	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
// 	// 	return
// 	// }
// 	//oauthClient := oauthConf.Client(oauth2.NoContext, token)
// 	//client := github.NewClient(oauthClient)
// 	//user, _, err := client.Users.Get("")
// 	//if err != nil {
// 	//	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
// 	//	return
// 	//}
// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)	
