package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type MediaService struct {
	Token         string
	PhoneNumberID string
	Client        *http.Client
}

func NewMediaService(token, phoneID string) *MediaService {
	return &MediaService{
		Token:         token,
		PhoneNumberID: phoneID,
		Client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type mediaResponse struct {
	URL      string `json:"url"`
	MimeType string `json:"mime_type"`
}

func (s *MediaService) GetMediaURL(mediaID string) (string, string, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s", mediaID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Authorization", "Bearer "+s.Token)

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 300 {
		return "", "", fmt.Errorf("erro API (%d): %s", resp.StatusCode, string(body))
	}

	var result mediaResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", "", err
	}

	return result.URL, result.MimeType, nil
}

func (s *MediaService) DownloadMedia(mediaURL string) ([]byte, error) {
	req, err := http.NewRequest("GET", mediaURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.Token)

	resp, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("erro download (%d): %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

func (s *MediaService) DownloadByID(mediaID, msgType string) (string, error) {

	url, mime, err := s.GetMediaURL(mediaID)
	if err != nil {
		return "", err
	}

	data, err := s.DownloadMedia(url)
	if err != nil {
		return "", err
	}

	ext := getExtension(mime)

	folder := "downloads/" + msgType
	err = os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return "", err
	}

	fileName := sanitize(fmt.Sprintf("%s_%d%s", mediaID, time.Now().Unix(), ext))
	filePath := fmt.Sprintf("%s/%s", folder, fileName)

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return "", err
	}

	log.Println("[MEDIA] Arquivo salvo:", filePath)

	return filePath, nil
}

func (s *MediaService) SendTextMessage(to, body string) error {

	url := fmt.Sprintf(
		"https://graph.facebook.com/v22.0/%s/messages",
		s.PhoneNumberID,
	)

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "text",
		"text": map[string]string{
			"body": body,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	log.Println("[WHATSAPP RESPONSE]:", string(respBody))

	if resp.StatusCode >= 300 {
		return fmt.Errorf("erro ao enviar mensagem")
	}

	return nil
}

func sanitize(name string) string {
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

func getExtension(mime string) string {
	switch {
	case strings.Contains(mime, "image/jpeg"):
		return ".jpg"
	case strings.Contains(mime, "image/png"):
		return ".png"
	case strings.Contains(mime, "audio/ogg"):
		return ".ogg"
	case strings.Contains(mime, "audio/mpeg"):
		return ".mp3"
	case strings.Contains(mime, "application/pdf"):
		return ".pdf"
	default:
		return ".bin"
	}
}
