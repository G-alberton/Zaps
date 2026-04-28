package webhook

type Message struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	Type      string `json:"type"`
	Timestamp string `json:"timestamp"`

	Text     *TextMessage     `json:"text,omitempty"`
	Image    *MediaMessage    `json:"image,omitempty"`
	Audio    *MediaMessage    `json:"audio,omitempty"`
	Document *MediaMessage    `json:"document,omitempty"`
	Video    *MediaMessage    `json:"video,omitempty"`
	Sticker  *MediaMessage    `json:"sticker,omitempty"`
	Location *LocationMessage `json:"location,omitempty"`

	Context *ContextMessage `json:"context,omitempty"`
}

type TextMessage struct {
	Body string `json:"body"`
}

type MediaMessage struct {
	ID string `json:"id"`
}

type LocationMessage struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Name      string  `json:"name"`
	Address   string  `json:"address"`
}

type ContextMessage struct {
	ID string `json:"id"`
}

type Event struct {
	Entry []Entry `json:"entry"`
}

type Entry struct {
	Changes []Change `json:"changes"`
}

type Change struct {
	Value Value `json:"value"`
}

type Value struct {
	Messages []Message `json:"messages"`
	Statuses []Status  `json:"statuses,omitempty"`
	Contacts []Contact `json:"contacts"`
}

type Status struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Recipient string `json:"recipient_id"`
}

type Contact struct {
	Profile Profile `json:"profile"`
	WaID    string  `json:"wa_id"`
}

type Profile struct {
	Name string `json:"name"`
}
