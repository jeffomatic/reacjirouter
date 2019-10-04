package tokenstore

import "sync"

var (
	teamCredsMutex sync.Mutex
	teamCreds      map[string]string
)

func init() {
	teamCreds = make(map[string]string)
}

func Add(teamID string, accessToken string) {
	teamCredsMutex.Lock()
	defer teamCredsMutex.Unlock()

	teamCreds[teamID] = accessToken
}

func Get(teamID string) (string, bool) {
	teamCredsMutex.Lock()
	defer teamCredsMutex.Unlock()
	accessToken, ok := teamCreds[teamID]
	return accessToken, ok
}
