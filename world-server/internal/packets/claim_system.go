package packets

type ClaimSystemPayload struct {
	SystemId int   `json:"systemId"`
	PlayerId int32 `json:"playerId"`
}
