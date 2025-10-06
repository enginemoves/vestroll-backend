package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuth struct {
	Config *oauth2.Config
}

type AppleOAuth struct {
	Config *oauth2.Config
}

type AppleClaims struct {
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Sub           string `json:"sub"`
	jwt.RegisteredClaims
}

type OAuthRequest struct {
	Code string `json:"code"`
}

type OAuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Provider string `json:"provider"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
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
	}
}

func NewAppleOAuth(clientID, clientSecret, redirectURL string) *AppleOAuth {
	return &AppleOAuth{
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"email", "name"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://appleid.apple.com/auth/authorize",
				TokenURL: "https://appleid.apple.com/auth/token",
			},
		},
	}
}

func parseAppleIDToken(idToken string) (*AppleClaims, error) { // I couldn't do it as verified
	token, _, err := jwt.NewParser().ParseUnverified(idToken, &AppleClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AppleClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims type")
	}

	return claims, nil
}

// --- GOOGLE HANDLERS ---
func (g *GoogleOAuth) GoogleHandleLogin(c *gin.Context) {
	url := g.Config.AuthCodeURL("state")
	c.JSON(http.StatusOK, gin.H{"auth_url": url})
}

func (g *GoogleOAuth) GoogleHandleCallback(c *gin.Context) {
	var req OAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Details: err.Error(),
		})
		return
	}

	if req.Code == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Code is required",
		})
		return
	}

	token, err := g.Config.Exchange(c.Request.Context(), req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to exchange token",
			Details: err.Error(),
		})
		return
	}

	client := g.Config.Client(c.Request.Context(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get user info",
			Details: err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		Email string `json:"email"`
		ID    string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to parse user info",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, OAuthResponse{
		Token: token.AccessToken,
		User: User{
			ID:       googleUser.ID,
			Email:    googleUser.Email,
			Provider: "google",
		},
	})
}

// --- APPLE HANDLERS ---
func (a *AppleOAuth) AppleHandleLogin(c *gin.Context) {
	url := a.Config.AuthCodeURL("state")
	c.JSON(http.StatusOK, gin.H{"auth_url": url})
}

func (a *AppleOAuth) AppleHandleCallback(c *gin.Context) {
	var req OAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Details: err.Error(),
		})
		return
	}

	if req.Code == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Code is required",
		})
		return
	}

	token, err := a.Config.Exchange(c.Request.Context(), req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to exchange token",
			Details: err.Error(),
		})
		return
	}

	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "No id_token in response",
		})
		return
	}

	claims, err := parseAppleIDToken(idToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to parse id_token",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, OAuthResponse{
		Token: token.AccessToken,
		User: User{
			ID:       claims.Sub,
			Email:    claims.Email,
			Provider: "apple",
		},
	})
}

