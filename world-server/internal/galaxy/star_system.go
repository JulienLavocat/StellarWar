package galaxy

import (
	"fmt"
	"math/rand"

	"com.stellarwar/world-server/internal/utils"
)

type StarSystem struct {
	Id    int32     `json:"id"`
	Name  string    `json:"name"`
	X     float64   `json:"x"`
	Y     float64   `json:"y"`
	To    utils.Set `json:"to"`
	Owner int32     `json:"owner"`
}

func NewStarSystem(seed string, id int32, x float64, y float64) StarSystem {
	rnd := rand.New(rand.NewSource(utils.HashString(fmt.Sprintf("%s-%f-%f", seed, x, y))))
	system := StarSystem{
		Id:    id,
		X:     x,
		Y:     y,
		To:    utils.NewSet(),
		Name:  generateName(rnd),
		Owner: -1,
	}

	return system
}
