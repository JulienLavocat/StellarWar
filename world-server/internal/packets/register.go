package packets

import (
	"github.com/StellarWar/world-server/internal/shared"
)

func NewRegisterPacket(playerId shared.PlayerId) *shared.ClientPacket {
	return &shared.ClientPacket{
		Type:     "register",
		Payload:  nil,
		PlayerId: playerId,
	}
}

func NewUnregisterPacket(playerId shared.PlayerId) *shared.ClientPacket {
	return &shared.ClientPacket{
		Type:     "unregister",
		Payload:  nil,
		PlayerId: playerId,
	}
}
