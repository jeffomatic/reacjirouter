// Package routestore provides an interface to a team-to-emoji-to-channel mapping.
// It is currently implemented as an in-memory routes, so the mapping will be
// lost after the process terminates.
package routestore

import (
	"sort"
	"strings"
	"sync"
)

var (
	routesByTeamMutex sync.Mutex
	routesByTeam      map[string]map[string]string
)

func init() {
	routesByTeam = make(map[string]map[string]string)
}

func Add(teamID string, emoji string, channelID string) {
	routesByTeamMutex.Lock()
	defer routesByTeamMutex.Unlock()

	teamRoutes, ok := routesByTeam[teamID]
	if !ok {
		teamRoutes = make(map[string]string)
		routesByTeam[teamID] = teamRoutes
	}

	teamRoutes[emoji] = channelID
}

type Route struct {
	Emoji     string
	ChannelID string
}

type RouteByEmoji []Route

func (a RouteByEmoji) Len() int           { return len(a) }
func (a RouteByEmoji) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a RouteByEmoji) Less(i, j int) bool { return strings.Compare(a[i].Emoji, a[j].Emoji) < 0 }

func listUnsorted(teamID string) []Route {
	routesByTeamMutex.Lock()
	defer routesByTeamMutex.Unlock()

	teamRoutes, ok := routesByTeam[teamID]
	if !ok {
		return nil
	}

	var res []Route
	for emoji, channel := range teamRoutes {
		res = append(res, Route{emoji, channel})
	}

	return res
}

func List(teamID string) []Route {
	list := listUnsorted(teamID)
	sort.Sort(RouteByEmoji(list))
	return list
}

func Get(teamID, emoji string) (string, bool) {
	routesByTeamMutex.Lock()
	defer routesByTeamMutex.Unlock()

	teamRoutes, ok := routesByTeam[teamID]
	if !ok {
		return "", false
	}

	channel, ok := teamRoutes[emoji]
	return channel, ok
}
