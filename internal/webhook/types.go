package webhook

type Message struct {
	From string `json:"from"`
	Type string `json:"type"`

	Text *struct {
		Body string `json:"body"`
	} `json:"text,omitempty"`

	Image *struct {
		ID string `json:"id"`
	} `json:"image,omitempty"`

	Audio *struct {
		ID string `json:"id"`
	} `json:"audio,omitempty"`

	Document *struct {
		ID string `json:"id"`
	} `json:"document,omitempty"`
}

// evento que aconteceu no processo de receber a mensagem
type Event struct {
	Entry []struct {
		Changes []struct {
			Value struct {
				Messages []Message `json:"messages"`
			} `json:"value"`
		} `json:"changes"`
	} `json:"entry"`
}
