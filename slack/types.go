package slack

type SlashCommandResponse struct {
	ResponseType string `json:"response_type"`
	Text         string `json:"text"`
}

type Event struct {
	ID          string `json:"event_id"`
	T           string `json:"type"`
	UserID      string `json:"user"`
	ChannelID   string `json:"channel"`
	Reaction    string
	Text        string `json:"text`
	ChannelType string `json:"channel_type"`
	Item        struct {
		T         string `json:"type"`
		ChannelID string `json:"channel"`
		Timestamp string `json:"ts"`
	}
}

type EventPayload struct {
	T         string `json:"type"`
	TeamID    string `json:"team_id"`
	Event     Event
	Challenge string // for URL verification
}

type AuthTestResponse struct {
	URL    string `json:"url"`
	TeamID string `json:"team_id"`
	UserID string `json:"user_id"`
}

type ChatPostEphemeralRequest struct {
	User    string `json:"user"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	AsUser  bool   `json:"as_user"`
}

type ChatPostMessageRequest struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
	AsUser  bool   `json:"as_user"`
}

type AccessTokenArgs struct {
	ClientID     string
	ClientSecret string
	Code         string
}

type OauthAccessResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TeamID      string `json:"team_id"`
}
