package action

import "strings"

type Path []string

func (p *Path) StringSlice() []string {
	return *p
}

func (p *Path) String() string {
	return strings.Join(*p, Delim)
}

func (p *Path) Append(next string) {
	*p = append(*p, next)
}

const (
	All   = "*"
	Delim = "."
)
