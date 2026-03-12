package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

var (
	ErrInvalidToken      = errors.New("invalid google token")
	ErrTokenVerification = errors.New("failed to verify google token")
)

const googleTokenInfoURL = "https://oauth2.googleapis.com/tokeninfo?id_token=%s"

// GoogleUserInfo represents user information from Google
type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
}

// VerifyGoogleIDToken verifies a Google ID token and returns user information
func VerifyGoogleIDToken(ctx context.Context, idToken string) (*GoogleUserInfo, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make request to Google's tokeninfo endpoint
	url := fmt.Sprintf(googleTokenInfoURL, idToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTokenVerification, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTokenVerification, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrInvalidToken
	}

	// Parse response
	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrTokenVerification, err)
	}

	// Validate required fields
	if userInfo.Sub == "" || userInfo.Email == "" {
		return nil, ErrInvalidToken
	}

	return &userInfo, nil
}
