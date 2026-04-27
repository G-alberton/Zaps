package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type MediaService struct {
	Token         string
	PhoneNumberID string
	Client        *http.Client
}

func (s *MediaService) sendRequest(ctx context.Context, url string, payload interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println("[WA RESPONSE]:", string(body))

	if resp.StatusCode >= 300 {
		return fmt.Errorf("erro (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func NewMediaService() *MediaService {

	token := os.Getenv("WHATSAPP_TOKEN")
	phoneID := os.Getenv("PHONE_NUMBER_ID")

	if token == "" || phoneID == "" {
		log.Fatal("WHATSAPP_TOKEN ou PHONE_NUMBER_ID não definido")
	}

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

func (s *MediaService) GetMediaURL(ctx context.Context, mediaID string) (string, string, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s", mediaID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Authorization", "Bearer "+s.Token)

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	if resp.StatusCode >= 300 {
		return "", "", fmt.Errorf("erro API (%d): %s", resp.StatusCode, string(body))
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		return "", "", err
	}

	if errData, ok := raw["error"]; ok {
		return "", "", fmt.Errorf("erro API: %v", errData)
	}

	mediaURL, _ := raw["url"].(string)
	mime, _ := raw["mime_type"].(string)

	if mediaURL == "" {
		return "", "", fmt.Errorf("url vazia")
	}

	return mediaURL, mime, nil
}

func (s *MediaService) DownloadMedia(ctx context.Context, mediaURL, filePath string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", mediaURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+s.Token)

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro download (%d): %s", resp.StatusCode, string(body))
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func (s *MediaService) DownloadByID(ctx context.Context, mediaID, msgType string) (string, error) {
	url, mime, err := s.GetMediaURL(ctx, mediaID)
	if err != nil {
		return "", err
	}

	ext := getExtension(mime)

	folder := "downloads/" + msgType
	if err := os.MkdirAll(folder, 0755); err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("%s_%d%s", sanitize(mediaID), time.Now().Unix(), ext)
	filePath := fmt.Sprintf("%s/%s", folder, fileName)

	if err := s.DownloadMedia(ctx, url, filePath); err != nil {
		return "", err
	}

	log.Println("[MEDIA] Salvo em:", filePath)
	return filePath, nil
}

// arrumando aqui
func (s *MediaService) SendTextMessage(ctx context.Context, to, bodyText string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages", s.PhoneNumberID)

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "text",
		"text": map[string]string{
			"body": bodyText,
		},
	}

	return s.sendRequest(ctx, url, payload)
}

func (s *MediaService) SendImageByID(ctx context.Context, to, mediaID, caption string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages", s.PhoneNumberID)

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "image",
		"image": map[string]string{
			"id":      mediaID,
			"caption": caption,
		},
	}

	return s.sendRequest(ctx, url, payload)
}

func (s *MediaService) SendDocumentByID(ctx context.Context, to, mediaID, caption, filename string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages", s.PhoneNumberID)

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "document",
		"document": map[string]string{
			"id":       mediaID,
			"caption":  caption,
			"filename": filename,
		},
	}

	return s.sendRequest(ctx, url, payload)
}

func (s *MediaService) UploadMedia(ctx context.Context, filePath string) (string, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/media", s.PhoneNumberID)

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}

	if _, err = io.Copy(part, file); err != nil {
		return "", err
	}

	if err = writer.WriteField("messaging_product", "whatsapp"); err != nil {
		return "", err
	}

	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", url, &b)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+s.Token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	log.Println("[UPLOAD MEDIA]:", string(body))

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("erro upload (%d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		ID string `json:"id"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.ID == "" {
		return "", fmt.Errorf("mediaID vazio")
	}

	return result.ID, nil
}

func (s *MediaService) SendAudioByID(ctx context.Context, to, mediaID string) error {
	url := fmt.Sprintf("https://graph.facebook.com/v22.0/%s/messages", s.PhoneNumberID)

	payload := map[string]interface{}{
		"messaging_product": "whatsapp",
		"to":                to,
		"type":              "audio",
		"audio": map[string]string{
			"id": mediaID,
		},
	}

	return s.sendRequest(ctx, url, payload)
}

func sanitize(name string) string {
	replacer := strings.NewReplacer("/", "_", "\\", "_", " ", "_", ":", "_", "*", "_")
	return replacer.Replace(name)
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
