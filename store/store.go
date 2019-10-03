// Package store provides an interface to a team-to-emoji-to-channel mapping.
// It is currently implemented as an in-memory store, so the mapping will be
// lost after the process terminates.
package store

import (
	"sort"
	"strings"
	"sync"
)

var (
	teamEmojiChannelMutex sync.Mutex
	teamEmojiChannel      map[string]map[string]string
)

func init() {
	teamEmojiChannel = make(map[string]map[string]string)
}

func Add(teamID string, emoji string, channelID string) {
	teamEmojiChannelMutex.Lock()
	defer teamEmojiChannelMutex.Unlock()

	emojiChannel, ok := teamEmojiChannel[teamID]
	if !ok {
		emojiChannel = make(map[string]string)
		teamEmojiChannel[teamID] = emojiChannel
	}

	emojiChannel[emoji] = channelID
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
	teamEmojiChannelMutex.Lock()
	defer teamEmojiChannelMutex.Unlock()

	emojiChannel, ok := teamEmojiChannel[teamID]
	if !ok {
		return nil
	}

	var res []Route
	for emoji, channel := range emojiChannel {
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
	teamEmojiChannelMutex.Lock()
	defer teamEmojiChannelMutex.Unlock()

	emojiChannel, ok := teamEmojiChannel[teamID]
	if !ok {
		return "", false
	}

	channel, ok := emojiChannel[emoji]
	return channel, ok
}
