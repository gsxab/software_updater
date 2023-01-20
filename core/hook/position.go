package hook

type PositionCommand int

const (
	FirstCmd PositionCommand = iota
	LastCmd
	PrevCmd
	NextCmd
)

type Position struct {
	Cmd PositionCommand
	Ref string
}

const (
	FirstStr = "first"
	LastStr  = "last"
	PrevStr  = "prev_to"
	NextStr  = "next_to"
)

var PositionCommandMap = map[PositionCommand]string{
	FirstCmd: FirstStr,
	LastCmd:  LastStr,
	PrevCmd:  PrevStr,
	NextCmd:  NextStr,
}

func First() *Position {
	return &Position{FirstCmd, ""}
}

func Last() *Position {
	return &Position{FirstCmd, ""}
}

func PrevTo(ref string) *Position {
	return &Position{PrevCmd, ref}
}

func NextTo(ref string) *Position {
	return &Position{NextCmd, ref}
}
