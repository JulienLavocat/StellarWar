package packets

type PlayerLeftPayload struct {
	Id           int32 `json:"id"`
	OwnedSystems []int `json:"ownedSystems"`
}
