package action

import "strings"

type Path []string

func (p Path) StringSlice() []string {
	return p
}

func (p Path) Name() string {
	return p[len(p)-1]
}

func (p Path) String() string {
	return strings.Join(p, Delim)
}

const (
	All   = "*"
	Delim = "."
)
