package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/jeffomatic/reacjirouter/routestore"
	"github.com/jeffomatic/reacjirouter/slack"
)

var spaceSplitter *regexp.Regexp

func init() {
	spaceSplitter = regexp.MustCompile(`\s+`)
}

func processCommand(teamID, text string) string {
	text = strings.TrimSpace(text)
	tokens := spaceSplitter.Split(text, 3)

	switch strings.ToLower(tokens[0]) {
	case "add":
		return processAddCommand(teamID, tokens)

	case "list":
		return processListCommand(teamID, tokens)

	case "help":
		return processHelpCommand()

	default:
		return "Sorry, I didn't recognize that command! Type \"help\" for instructions."
	}
}

func processNonMember(c *teamClient, targetChannelID string) string {
	userID, err := c.userID()
	if err != nil {
		log.Println("error looking up ID:", err)
		return fmt.Sprintf("We got an error trying to look up that channel.")
	}

	return fmt.Sprintf(`
*Oops!* Reacji Router isn't a member of <#%s>. Try inviting the app:

%s
/invite <@%s> <#%s>
%s
`,
		targetChannelID,
		"```",
		userID,
		targetChannelID,
		"```",
	)
}

func processAddCommand(teamID string, tokens []string) string {
	if len(tokens) != 3 {
		return `Could not understand add command`
	}

	emoji, ok := slack.ExtractEmoji(tokens[1])
	if !ok {
		return `Could not understand add command`
	}

	targetChannelID, ok := slack.ExtractChannelID(tokens[2])
	if !ok {
		return `Could not understand add command`
	}

	c := newTeamClient(teamID)
	var resp slack.ConversationsInfoResponse
	err := c.Call(slack.ConversationsInfo, slack.ConversationsInfoRequest{ChannelID: targetChannelID}, &resp)
	if err != nil {
		if apiErr := slack.GetAPIError(err); apiErr != nil {
			if apiErr.Error() == "channel_not_found" {
				return processNonMember(c, targetChannelID)
			}

			return fmt.Sprintf("We got an error from Slack trying to look up that channel: %s", apiErr.Error())
		}

		log.Println("unknown conversations.info error:", err)
		return fmt.Sprintf("We got an error trying to look up that channel.")
	}

	if !resp.Channel.IsMember {
		return processNonMember(c, targetChannelID)
	}

	// If we've gotten here, the bot user is a member of the target channel. We
	// can add the route and return a success message.
	routestore.Add(teamID, emoji, targetChannelID)
	return fmt.Sprintf("Okay, I'll send all messages with the :%s: reacji to <#%s>.", emoji, targetChannelID)
}

func processListCommand(teamID string, tokens []string) string {
	if len(tokens) > 1 {
		return `Could not understand add command`
	}

	text := `No reacji configured.`
	list := routestore.List(teamID)
	if len(list) > 0 {
		lines := make([]string, len(list))
		for _, pair := range list {
			lines = append(lines, fmt.Sprintf(":%s: :point_right: <#%s>", pair.Emoji, pair.ChannelID))
		}
		text = strings.Join(lines, "\n")
	}

	return text
}

func processHelpCommand() string {
	return `
*Instructions*

Add a new reaction route
> add :emoji: #channel

List reaction routes on this team
> list

Show help
> help
`
}
