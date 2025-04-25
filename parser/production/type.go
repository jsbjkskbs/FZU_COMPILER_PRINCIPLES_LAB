package production

import (
	"slices"
)

type Production struct {
	Head Symbol
	Body []Symbol
}

func (p Production) Equals(other Production) bool {
	if p.Head != other.Head {
		return false
	}
	return slices.Equal(p.Body, other.Body)
}

type Symbol string

func (s Symbol) IsEpsilon() bool {
	return s == EPSILON
}

type Terminal string

func (t Terminal) IsEpsilon() bool {
	return t == EPSILON
}
