package slack

type Method string

var (
	AuthTest          Method = "auth.test"
	ChatPostMessage   Method = "chat.postMessage"
	ConversationsInfo Method = "conversations.info"
	OauthAccess       Method = "oauth.access"
)

type methodInfo struct {
	noAccessToken      bool
	formEncodedRequest bool
}

var infoByMethod map[Method]methodInfo

func init() {
	infoByMethod = map[Method]methodInfo{
		ConversationsInfo: {formEncodedRequest: true},
		OauthAccess:       {noAccessToken: true, formEncodedRequest: true},
	}
}
