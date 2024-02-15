package entity

type Room struct {
	Name string
	X    int
	Y    int
}

type Position struct {
	Start *Room
	End   *Room
}

type Path []*Room

type FarmGraph struct {
	AntCount    int
	Start       *Room
	End         *Room
	Rooms       []*Room
	Positions   []*Position
	Paths       []*Path
	RoomNames   map[string]*Room
	Adadjacents map[*Room][]*Room
}
