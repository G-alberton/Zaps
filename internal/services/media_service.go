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
)

type MediaService struct {
	Token         string
	PhoneNumberID string
}

func NewMediaService(token, phoneID string) *MediaService {
	return &MediaService{
		Token:         token,
		PhoneNumberID: phoneID,
	}
}

type mediaResponse struct {
	URL      string `json:"url"`
	mimeType string `json:"mime_type"`
}

func (s *MediaService) GetMediaURL(mediaID string) (string, string, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v19.0/%s", mediaID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Add("Authorization", "Bearer "+s.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()

	var result mediaResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", "", err
	}

	return result.URL, result.mimeType, nil
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
		return ""
	}
}

func (s *MediaService) DownloadMedia(mediaURL, filePath string) error {

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
	if err != nil {
		return err
	}

	return nil
}

func saveToFile(path string, data []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

func (s *MediaService) DownloadByID(mediaID, msgType string) (string, error) {

	url, mime, err := s.GetMediaURL(mediaID)
	if err != nil {
		return "", err
	}

	ext := getExtension(mime)

	folder := "downloads/" + msgType
	err = os.MkdirAll(folder, os.ModePerm)
	if err != nil {
		return "", err
	}

	filePath := fmt.Sprintf("%s/%s%s", folder, mediaID, ext)

	err = s.DownloadMedia(url, filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func (s *MediaService) SendTextMessage(to, body string) error {
	url := fmt.Sprintf(
		"https://graph.facebook.com/v22.0/%s/messages",
		s.PhoneNumberID,
	)

	data := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "text",
		"text": map[string]string{
			"body": body,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+s.Token)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	log.Println("Resposta WhatsApp:", string(respBody))

	if resp.StatusCode >= 300 {
		return fmt.Errorf("erro ao enviar mensagem")
	}

	return nil
}
