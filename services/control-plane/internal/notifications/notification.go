package notifications

type Notification struct {
	ID      string `json:"id"`
	Channel string `json:"channel"`
	Target  string `json:"target"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	Status  string `json:"status"`
}
