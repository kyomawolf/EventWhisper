package identities

type Identity struct {
	Sub             string           `json:"sub"`
	Name            string           `json:"name"`
	Location        string           `json:"location"`
	Interests       []string         `json:"interests"`
	Channels        []Channel        `json:"channels"`
	AnnouncedEvents []AnnouncedEvent `json:"announcedEvents"`
}

type Channel struct {
	ID          string           `json:"id"`
	ChannelName string           `json:"channelname"`
	Type        string           `json:"type"` // "directmessage" or "group"
	Specifics   ChannelSpecifics `json:"specifics"`
}

type ChannelSpecifics struct {
	ChatID string `json:"chatId"`
}

type AnnouncedEvent struct {
	ID          string `json:"id"`
	EventID     string `json:"eventid"`
	AnnouncedAt string `json:"announced_at"`
	DeleteAt    string `json:"delete_at"`
}
