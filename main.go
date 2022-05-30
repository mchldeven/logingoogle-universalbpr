package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		RedirectURL: "http://localhost:6969/callback",
		ClientID: os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{"https://wwww.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
	
	randomState = "random"
)

func main() {
	http.HandleFunc("/", HandleHome)
	http.HandleFunc("/login", HandleLogin)
	http.HandleFunc("/callback", HandleCallback)
	http.ListenAndServe(":6969", nil)
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	var html = `<html><body><a href="/login">GOOGLE LOGIN UNIVERSALBPR</a></body></html>`
	fmt.Fprint(w, html)
}

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleCallback(w http.ResponseWriter, r *http.Request) {
if r.FormValue("stage") != randomState {
	fmt.Println("stage error")
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return	
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, r.FormValue("code"))
	if err != nil {
		fmt.Println("gagal mendapatkan token: %s/n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get("https://googleapis.com/oauth2/v2/userinfo?access_token="+token.AccessToken)
	if err != nil {
		fmt.Printf("gagal mnendapatkan request get: %s/n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("gagal mnendapatkan request get: %s/n", err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "Response: %s", content)
}



