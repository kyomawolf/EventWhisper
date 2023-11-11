package events

type Event struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	StartTime   string   `json:"start_time"`
	EndTime     string   `json:"end_time"`
	Organizer   string   `json:"organizer"`
	Interest    []string `json:"interest"`
}
