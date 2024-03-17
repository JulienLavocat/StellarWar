package packets

import "com.stellarwar/world-server/internal/galaxy"

type PublicClientPayload struct {
	Id             int32  `json:"id"`
	Color          [3]int `json:"color"`
	HomeSystem     int    `json:"homeSystem"`
	ClaimedSystems []int  `json:"claimedSystems"`
}

type SyncClientPayload struct {
	Galaxy  galaxy.Galaxy                 `json:"galaxy"`
	Players map[int32]PublicClientPayload `json:"players"`
	Self    PublicClientPayload           `json:"self"`
}
