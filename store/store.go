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

type Pair struct {
	Emoji     string
	ChannelID string
}

type PairByEmoji []Pair

func (a PairByEmoji) Len() int           { return len(a) }
func (a PairByEmoji) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a PairByEmoji) Less(i, j int) bool { return strings.Compare(a[i].Emoji, a[j].Emoji) < 0 }

func listUnsorted(teamID string) []Pair {
	teamEmojiChannelMutex.Lock()
	defer teamEmojiChannelMutex.Unlock()

	emojiChannel, ok := teamEmojiChannel[teamID]
	if !ok {
		return nil
	}

	var res []Pair
	for emoji, channel := range emojiChannel {
		res = append(res, Pair{emoji, channel})
	}

	return res
}

func List(teamID string) []Pair {
	list := listUnsorted(teamID)
	sort.Sort(PairByEmoji(list))
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
