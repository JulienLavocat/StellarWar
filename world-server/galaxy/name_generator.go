package galaxy

import (
	"math/rand"
	"strings"
)

var vowels = [][]string{
	{"b", "c", "d", "f", "g", "h", "i", "j", "k", "l", "m", "n", "p", "q", "r", "s", "t", "v", "w", "x", "y", "z"},
	{"a", "e", "o", "u"},
	{"br", "cr", "dr", "fr", "gr", "pr", "str", "tr", "bl", "cl", "fl", "gl", "pl", "sl", "sc", "sk", "sm", "sn", "sp", "st", "sw", "ch", "sh", "th", "wh"},
	{"ae", "ai", "ao", "au", "a", "ay", "ea", "ei", "eo", "eu", "e", "ey", "ua", "ue", "ui", "uo", "u", "uy", "ia", "ie", "iu", "io", "iy", "oa", "oe", "ou", "oi", "o", "oy"},
	{"turn", "ter", "nus", "rus", "tania", "hiri", "hines", "gawa", "nides", "carro", "rilia", "stea", "lia", "lea", "ria", "nov", "phus", "mia", "nerth", "wei", "ruta", "tov", "zuno", "vis", "lara", "nia", "liv", "tera", "gantu", "yama", "tune", "ter", "nus", "cury", "bos", "pra", "thea", "nope", "tis", "clite"},
	{"una", "ion", "iea", "iri", "illes", "ides", "agua", "olla", "inda", "eshan", "oria", "ilia", "erth", "arth", "orth", "oth", "illon", "ichi", "ov", "arvis", "ara", "ars", "yke", "yria", "onoe", "ippe", "osie", "one", "ore", "ade", "adus", "urn", "ypso", "ora", "iuq", "orix", "apus", "ion", "eon", "eron", "ao", "omia"},
}

var mtx = [][]int{
	{1, 1, 2, 2, 5, 5},
	{2, 2, 3, 3, 6, 6},
	{3, 3, 4, 4, 5, 5},
	{4, 4, 3, 3, 6, 6},
	// {3, 3, 4, 4, 2, 2, 5, 5},
	// {2, 2, 1, 1, 3, 3, 6, 6},
	// {3, 3, 4, 4, 2, 2, 5, 5},
	// {4, 4, 3, 3, 1, 1, 6, 6},
	// {3, 3, 4, 4, 1, 1, 4, 4, 5, 5},
	// {4, 4, 1, 1, 4, 4, 3, 3, 6, 6},
}

func generateName(rnd *rand.Rand) string {
	name := strings.Builder{}

	comp := mtx[rnd.Intn(len(mtx))]
	il := len(comp) / 2
	for i := 0; i < il; i++ {
		name.WriteString(vowels[comp[i*2]-1][rnd.Intn(len(vowels[comp[i*2]-1]))])
	}

	return strings.Title(name.String())
}
