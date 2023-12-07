package events

type Event struct {
	ID          string        `json:"id"`
	Source      string        `json:"source"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Location    EventLocation `json:"location"`
	StartTime   string        `json:"start_date_time"`
	EndTime     string        `json:"end_date_time"`
	Organizer   string        `json:"organizer"`
	Pricing     string        `json:"pricing"`
	Url         string        `json:"url"`
	Interests   []string      `json:"interests"`
}

type EventLocation struct {
	City      string `json:"city"`
	Country   string `json:"country"`
	Street    string `json:"street"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	Zip       string `json:"zip"`
}
