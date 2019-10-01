package slack

import (
	"fmt"
	"strings"
)

type MessageLink struct {
	TeamURL   string
	ChannelID string
	Timestamp string
}

func (m MessageLink) String() string {
	return fmt.Sprintf(
		"%s/archives/%s/p%s",
		strings.TrimSuffix(m.TeamURL, "/"), // tolerate the trailing suffix returned by auth.test
		m.ChannelID,
		strings.Replace(m.Timestamp, ".", "", 1),
	)
}
