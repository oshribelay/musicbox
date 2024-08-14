package spotify

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/oshribelay/musicbox/internal/logger"
)

var (
	accessToken 	string
	tokenExpiration int
	authCode		string
)

type SpotifyTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type SpotifyAuthCodeResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope		 string `json:"scope"`
}

func RequestAccessTokenFromSpotify() error {
	var tokenResponse SpotifyTokenResponse

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", os.Getenv("CLIENT_ID"))
	data.Set("client_secret", os.Getenv("CLIENT_SECRET"))
	
	resp, err := http.Post(
		"https://accounts.spotify.com/api/token", 
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		logger.ErrorLogger.Println("error in fetching access token:", err)
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println("Error reading body:", err)
		return err
	}

	
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return err
	}

	setAccessToken(tokenResponse.AccessToken, tokenResponse.ExpiresIn)

	fmt.Printf("access token: %s \n expires in: %d", tokenResponse.AccessToken, tokenResponse.ExpiresIn)
	return nil
}

func setAccessToken(token string, expiresIn int) {
	accessToken = token
	tokenExpiration = expiresIn
}

func getAccessToken() (string, error) {
	if tokenExpiration <= 0 {
		if err := RequestAccessTokenFromSpotify(); err != nil {
			return accessToken, err
		}
	}

	return accessToken, nil
}

func generateRandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
    seed := rand.NewSource(time.Now().UnixNano())
    random := rand.New(seed)

    result := make([]byte, length)
    for i := range result {
        result[i] = charset[random.Intn(len(charset))]
    }
    return string(result)
}

func GenerateAuthUrl() (string, error) {
	state := generateRandomString(16)
	scope := "user-read-private user-read-email"
	redirectUri := "http://localhost:8080/spotify/callback"

	base, err := url.Parse("https://accounts.spotify.com/authorize")

	if err != nil {
        return "", err 
    }

	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", os.Getenv("CLIENT_ID"))
	params.Add("scope", scope)
	params.Add("redirect_uri", redirectUri)
	params.Add("state", state)

	base.RawQuery = params.Encode()

	return base.String(), nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	url, err := GenerateAuthUrl()
	if err != nil {
		logger.ErrorLogger.Println("error in generating auth url", err)
	}

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func RequestAuthCode(code string) string {
	client := &http.Client{}
	var spotifyAuthCodeResponse SpotifyAuthCodeResponse
	
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", "http://localhost:8080/spotify/callback")
	
	req, _ := http.NewRequest(
		"POST",
	 	"https://accounts.spotify.com/api/token", 
		strings.NewReader(data.Encode()),
	)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic " + b64.StdEncoding.EncodeToString([]byte(os.Getenv("CLIENT_ID") + ":" + os.Getenv("CLIENT_SECRET"))))
	
	resp, err := client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		logger.ErrorLogger.Println("error in response: ", err)
	}

	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)

	logger.InfoLogger.Println("code: " + code)

	if err != nil {
		fmt.Println("Error reading body:", err)
	}

	if err := json.Unmarshal(body, &spotifyAuthCodeResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
	}

	authCode = spotifyAuthCodeResponse.AccessToken

	logger.InfoLogger.Println("Auth code inside: " + authCode)

	return authCode

}

func GetAuthCode() string {
	return authCode
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	logger.InfoLogger.Println("Auth code: " + RequestAuthCode(code))
}