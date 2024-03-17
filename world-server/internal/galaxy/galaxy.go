package galaxy

import (
	"math/rand"

	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"
)

type Galaxy struct {
	systemsTree *kdtree.KDTree
	Systems     []StarSystem `json:"systems"`
	Routes      [][2]int32   `json:"routes"`
}

func (g *Galaxy) KNN(x float64, y float64, radius float64) []*StarSystem {
	around := g.systemsTree.KNN(&points.Point2D{X: x, Y: y}, 10)
	systems := make([]*StarSystem, len(around))
	for i, point := range around {
		systems[i] = &g.Systems[point.(*points.Point).Data.(int32)]
	}
	return systems
}

func (g *Galaxy) AssignNextHomeSystem(playerId int32) int {
	systemId := rand.Intn(len(g.Systems)) // TODO: Find next home system using a better method

	for g.Systems[systemId].Owner != -1 {
		systemId = rand.Intn(len(g.Systems))
	}

	g.Systems[systemId].Owner = playerId

	return systemId
}
