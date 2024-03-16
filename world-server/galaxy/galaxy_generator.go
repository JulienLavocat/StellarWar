package galaxy

import (
	"fmt"
	"math"
	"math/rand"

	"com.stellarwar/world-server/utils"
	"github.com/fogleman/poissondisc"
	"github.com/kyroy/kdtree"
	"github.com/kyroy/kdtree/points"
	"github.com/pradeep-pyro/triangle"
	"gonum.org/v1/gonum/graph/path"
	graphs "gonum.org/v1/gonum/graph/simple"
)

type GalaxyGenerator struct {
	Width          float64
	Height         float64
	SystemsDensity float64
	Seed           string

	rnd       *rand.Rand
	routes    [][2]int32
	systems   [][2]float64
	triangles [][3]int32
	lengths   [][3]float64
}

func getLength(a [2]float64, b [2]float64) float64 {
	abDstX := (a[0] - b[0])
	abDstY := (a[1] - b[1])
	return math.Sqrt(abDstX*abDstX + abDstY*abDstY)
}

func (m *GalaxyGenerator) Generate() Galaxy {
	m.rnd = rand.New(rand.NewSource(utils.HashString(m.Seed)))

	m.generateSystems()
	m.generateRoutes()

	systems := make([]StarSystem, len(m.systems))
	names := make(map[string]int)
	for i, point := range m.systems {
		system := NewStarSystem(m.Seed, int32(i), point[0], point[1])

		names[system.Name]++
		if names[system.Name] > 1 {
			system.Name = fmt.Sprintf("%s %d", system.Name, names[system.Name])
		}

		systems[i] = system
	}

	for _, route := range m.routes {
		systems[route[0]].To.Add(route[1])
		systems[route[1]].To.Add(route[0])
	}

	systemsPoints := make([]kdtree.Point, len(systems))
	for i, system := range systems {
		systemsPoints[i] = points.NewPoint([]float64{system.X, system.Y}, system.Id)
	}
	systemsTree := kdtree.New(systemsPoints)

	return Galaxy{
		Systems:     systems,
		Routes:      m.routes,
		systemsTree: systemsTree,
	}
}

func (m *GalaxyGenerator) generateSystems() {
	poissonSampling := poissondisc.Sample(-m.Width/2, -m.Height/2, m.Width/2, m.Height/2, m.SystemsDensity, 30, m.rnd)

	points := make([][2]float64, len(poissonSampling))
	for i, p := range poissonSampling {
		points[i] = [2]float64{p.X, p.Y}
	}

	triangles := triangle.Delaunay(points)

	lengths := make([][3]float64, len(triangles))
	for i, triangle := range triangles {
		a := points[triangle[0]]
		b := points[triangle[1]]
		c := points[triangle[2]]

		lengths[i] = [3]float64{getLength(a, b), getLength(b, c), getLength(c, a)}
	}

	m.systems = points
	m.triangles = triangles
	m.lengths = lengths
}

func (m *GalaxyGenerator) generateRoutes() {
	graph := graphs.NewWeightedUndirectedGraph(0, math.MaxFloat64)
	for i := 0; i < len(m.triangles); i++ {
		a := m.triangles[i][0]
		b := m.triangles[i][1]
		c := m.triangles[i][2]
		weights := m.lengths[i]
		graph.SetWeightedEdge(graphs.WeightedEdge{
			F: graphs.Node(int64(a)),
			T: graphs.Node(int64(b)),
			W: weights[0],
		})
		graph.SetWeightedEdge(graphs.WeightedEdge{
			F: graphs.Node(int64(b)),
			T: graphs.Node(int64(c)),
			W: weights[1],
		})
		graph.SetWeightedEdge(graphs.WeightedEdge{
			F: graphs.Node(int64(c)),
			T: graphs.Node(int64(a)),
			W: weights[2],
		})
	}

	mst := graphs.NewWeightedUndirectedGraph(0, math.MaxFloat64)
	path.Kruskal(mst, graph)
	edges := mst.Edges()

	var routes [][2]int32
	for {
		if !edges.Next() {
			break
		}
		routes = append(routes, [2]int32{int32(edges.Edge().From().ID()), int32(edges.Edge().To().ID())})
	}
	m.routes = routes
}
