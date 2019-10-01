package slack

import (
	"fmt"
	"regexp"
	"strings"
)

var channelParser *regexp.Regexp

func init() {
	channelParser = regexp.MustCompile(`^<#(\w+)\|([^>]+)>$`)
}

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

// Extracts an emoji name from a formatting string such as ":grin:".
func ExtractEmoji(s string) (string, bool) {
	if len(s) < 3 {
		return "", false
	}

	if s[0] != ':' || s[len(s)-1] != ':' {
		return "", false
	}

	return strings.Trim(s, ":"), true
}

// Extracts a channel ID from messsage token such as "<#C7LTCHERE|general>".
func ExtractChannelID(s string) (string, bool) {
	matches := channelParser.FindAllStringSubmatch(s, 1)
	if len(matches) == 0 {
		return "", false
	}

	submatches := matches[0]
	if len(submatches) != 3 {
		return "", false
	}

	return submatches[1], true
}
