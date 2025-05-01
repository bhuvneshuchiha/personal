package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load
	if err != nil {
		log.Fatal("Error loading in the env file")
	}
}

func handleOAuthLogin(w http.ResponseWriter, r *http.Request){
	service := r.URL.Query().Get("service")
	var authURL string

	if service == "microsoft" {
		authURL = fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize?client_id=%s&response_type=code&redirect_uri=%s&response_mode=query&scope=offline_access%%20Chat.Read%%20Mail.Read&state=random_string",
			os.Getenv("MICROSOFT_TENANT_ID"), os.Getenv("MICROSOFT_CLIENT_ID"), os.Getenv("MICROSOFT_REDIRECT_URI"))

	} else if service == "gitlab" {
		authURL = fmt.Sprintf("https://gitlab.com/oauth/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=api%%20read_user",
			os.Getenv("GITLAB_CLIENT_ID"), os.Getenv("GITLAB_REDIRECT_URI"))
	} else {
		http.Error(w, "Invalid Service", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, authURL, http.StatusFound)
}


func handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Query().Get("service")
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Error receiving the code", http.StatusBadRequest)
		return
	}

	var tokenEndpoint string
	var clientID, clientSecret, redirectURI string

	if service == "microsoft" {
		tokenEndpoint = "https://login.microsoftonline.com/" + os.Getenv("MICROSOFT_TENANT_ID") + "/oauth2/v2.0/token"
		clientID = os.Getenv("MICROSOFT_CLIENT_ID")
		clientSecret = os.Getenv("MICROSOFT_CLIENT_SECRET")
		redirectURI = os.Getenv("MICROSOFT_REDIRECT_URI")
	} else if service == "gitlab" {
		tokenEndpoint = "https://gitlab.com/oauth/token"
		clientID = os.Getenv("GITLAB_CLIENT_ID")
		clientSecret = os.Getenv("GITLAB_CLIENT_SECRET")
		redirectURI = os.Getenv("GITLAB_REDIRECT_URI")
	}else {
		http.Error(w, "Invalid service", http.StatusBadRequest)
		return
	}

	token, err := exchangeCodeForToken(tokenEndpoint, clientSecret, clientID, redirectURI, code)
	if err != nil {
		http.Error(w, "Failed to get the token", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w,"OAuth token successfully received -> %s", token.AccessToken)
}



