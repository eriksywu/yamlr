package utils

type Set map[interface{}]struct{}

func (s *Set) Has(v interface{}) bool {
	_, k := (*s)[v]
	return k
}

func (s *Set) Add(v interface{}) bool {
	if _, k := (*s)[v]; k {
		return false
	}
	(*s)[v] = struct{}{}
	return true
}

func (s *Set) Remove(v interface{}) bool {
	if _, k := (*s)[v]; k {
		return false
	}
	delete(*s, v)
	return true
}

func (s *Set) Intersect(other Set) *Set {
	intersection := NewSet()
	for v := range other {
		if s.Has(v) {
			intersection.Add(v)
		}
	}
	return intersection
}

func (s Set) ToSlice() []interface{} {
	slice := make([]interface{}, 0)
	for k, _ := range s {
		slice = append(slice, k)
	}
	return slice
}

func NewSet() *Set {
	set := make(Set)
	return &set
}
