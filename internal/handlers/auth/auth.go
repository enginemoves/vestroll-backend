package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/codeZe-us/vestroll-backend/internal/config"
)

var cfg = config.Load()
var applecfg = cfg.Apple

// Creates a OAuth instance
type GoogleOAuth struct {
	Config *oauth2.Config
	State  string
}

func NewGoogleOAuth(clientID, clientSecret, redirectURL string) *GoogleOAuth {
	return &GoogleOAuth{
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
		State: generateState(),
	}
}

// Generates a random state string for OAuth2
func generateState() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic("failed to generate random state")
	}
	return base64.URLEncoding.EncodeToString(b)
}

// --- GOOGLE HANDLERS (http.HandlerFunc style) ---
func GoogleLoginHandler(googleOauthConfig *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := googleOauthConfig.AuthCodeURL("state")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

// Google Callback Handler with enhanced error handling
func GoogleCallbackHandler(googleOauthConfig *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "No code in request", http.StatusBadRequest)
			return
		}

		token, err := googleOauthConfig.Exchange(r.Context(), code)
		if err != nil {
			http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
			return
		}

		client := googleOauthConfig.Client(r.Context(), token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			http.Error(w, "Failed to get user info", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var googleUser struct {
			Email string `json:"email"`
			ID    string `json:"id"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
			http.Error(w, "Failed to parse user info", http.StatusInternalServerError)
			return
		}

		// Security best practice: Check user existence and create session securely
		// Check if user exists, if not create a new user
		// Then log the user in
		// ...
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Logged in as %s", googleUser.Email)
	}
}

// --- APPLE HANDLERS ---
