package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type MediaService struct {
	Token string
}

func NewMediaService(token string) *MediaService {
	return &MediaService{Token: token}
}

func (s *MediaService) GetMediaURL(mediaID string) (string, error) {
	url := fmt.Sprintf("http://graph.facebook.com/v19.0/%s", mediaID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+s.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		URL string `json:"url"`
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.URL, nil
}

func (s *MediaService) DowloadMedia(mediaURL, filePath string) error {

	req, err := http.NewRequest("GET", mediaURL, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+s.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func (s *MediaService) DowloadByID(mediaID, filePath string) error {
	url, err := s.GetMediaURL(mediaID)
	if err != nil {
		return err
	}

	return s.DowloadMedia(url, filePath)
}
