func SendMedia(
	mediaService *services.MediaService,
	messageService *services.MessageService,
	conversationService *services.ConversationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "erro ao ler arquivo", http.StatusBadRequest)
			return
		}
		defer file.Close()

		if header.Size == 0 {
			http.Error(w, "arquivo vazio", http.StatusBadRequest)
			return
		}

		buffer := make([]byte, 512)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			http.Error(w, "erro ao ler arquivo", 500)
			return
		}

		if n == 0 {
			http.Error(w, "arquivo inválido", 400)
			return
		}

		realType := http.DetectContentType(buffer[:n])

		if !strings.HasPrefix(realType, "image/") &&
			!strings.HasPrefix(realType, "audio/") &&
			realType != "application/pdf" {

			http.Error(w, "tipo de arquivo não permitido", 400)
			return
		}

		if strings.HasPrefix(realType, "application/") && realType != "application/pdf" {
			http.Error(w, "somente PDF permitido", 400)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, "erro ao reposicionar arquivo", 500)
			return
		}

		if !strings.HasPrefix(realType, "image/") &&
			!strings.HasPrefix(realType, "audio/") &&
			!strings.HasPrefix(realType, "application/") {
			http.Error(w, "tipo de arquivo não permitido", 400)
			return
		}

		to := strings.TrimSpace(r.FormValue("to"))
		caption := r.FormValue("caption")

		if to == "" {
			http.Error(w, "destinatário obrigatório", http.StatusBadRequest)
			return
		}

		filename := fmt.Sprintf("%d_%s",
			time.Now().Unix(),
			strings.ReplaceAll(filepath.Base(header.Filename), " ", "_"),
		)

		var folder string
		if strings.HasPrefix(realType, "image/") {
			folder = "images"
		} else if strings.HasPrefix(realType, "audio/") {
			folder = "audio"
		} else {
			folder = "files"
		}

		fullPath := filepath.Join("uploads", folder, filename)

		if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
			http.Error(w, "erro ao criar diretório", 500)
			return
		}

		out, err := os.Create(fullPath)
		if err != nil {
			http.Error(w, "erro ao salvar arquivo", 500)
			return
		}
		defer out.Close()

		if _, err = io.Copy(out, file); err != nil {
			http.Error(w, "erro ao salvar arquivo", 500)
			return
		}

		publicURL := fmt.Sprintf("http://localhost:8080/uploads/%s/%s", folder, filename)

		var mediaID string

		if strings.HasPrefix(realType, "image/") {

			err = mediaService.SendImageByURL(ctx, to, publicURL, caption)
			mediaID = publicURL

		} else if strings.HasPrefix(realType, "audio/") {

			mediaIDUpload, errUpload := mediaService.UploadMedia(ctx, fullPath)
			if errUpload != nil {
				http.Error(w, errUpload.Error(), 500)
				return
			}

			err = mediaService.SendAudioByID(ctx, to, mediaIDUpload)
			mediaID = mediaIDUpload
		} else if strings.HasPrefix(realType, "application/") {

			mediaIDUpload, errUpload := mediaService.UploadMedia(ctx, fullPath)
			if errUpload != nil {
				http.Error(w, errUpload.Error(), 500)
				return
			}

			err = mediaService.SendDocumentByID(
				ctx,
				to,
				mediaIDUpload,
				caption,
				header.Filename,
			)

			mediaID = mediaIDUpload
		}

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		msgType := "file"
		if strings.HasPrefix(realType, "image/") {
			msgType = "image"
		} else if strings.HasPrefix(realType, "audio/") {
			msgType = "audio"
		} else if strings.HasPrefix(realType, "application/") {
			msgType = "document"
		}

		conversationID := conversationService.GetOrCreate(to)

		msg := models.Message{
			From:           to,
			ConversationID: conversationID,
			Type:           msgType,
			Body:           "",
			MediaID:        mediaID,
			MediaURL:       publicURL,
			Direction:      "outbound",
			Timestamp:      time.Now(),
		}

		messageService.SaveMessage(msg)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"url": publicURL,
		})
	}
}