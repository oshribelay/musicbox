package spotify

import (
	"io"
	"net/http"

	"github.com/oshribelay/musicbox/internal/logger"
)

func GetCurrentUserPlaylists() error {
	client := &http.Client{}	

	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/playlists", nil)

	if err != nil {
		logger.ErrorLogger.Println("error in request:", err)
		return err
	}

	token, _ := getAccessToken()
	req.Header.Add("Authorization", "Bearer " + token)
	
	resp, err := client.Do(req)
	
	if err != nil {
		logger.ErrorLogger.Println("error in response: ", err)
	}

	defer resp.Body.Close()

	if body, err := io.ReadAll(resp.Body); err != nil {
		logger.ErrorLogger.Println("Error reading body:", err)
		return err
	} else {
		logger.InfoLogger.Printf("body %s", body)
		return nil
	}
}