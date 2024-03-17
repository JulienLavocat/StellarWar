package utils

import "encoding/json"

type Set struct {
	list map[int32]struct{} // empty structs occupy 0 memory
}

func (s *Set) Has(v int32) bool {
	_, ok := s.list[v]
	return ok
}

func (s *Set) Add(v int32) {
	s.list[v] = struct{}{}
}

func (s *Set) Remove(v int32) {
	delete(s.list, v)
}

func (s *Set) Clear() {
	s.list = make(map[int32]struct{})
}

func (s *Set) Size() int {
	return len(s.list)
}

func NewSet() Set {
	s := Set{}
	s.list = make(map[int32]struct{})
	return s
}

func (s *Set) MarshalJSON() ([]byte, error) {
	keys := make([]int32, len(s.list))
	i := 0
	for k := range s.list {
		keys[i] = k
		i++
	}
	return json.Marshal(keys)
}
