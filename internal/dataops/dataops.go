package dataops

import (
	"math"
	"strconv"
	"strings"

	"lem-in/internal/entity"
	"lem-in/pkg/fsops"
	"lem-in/pkg/serror"
)

// CreateGraph creates a farm graph based on the input data
// It parses the game parameters, checks and stores room descriptions,
// performs additional checks on the positions of the rooms,
// and returns the populated farm graph
func CreateGraph(data []byte) (*entity.FarmGraph, error) {
	// Parse the game parameters from the input data
	gameParams, err := fsops.ParseFile(data)
	if err != nil {
		return nil, err
	}

	// Parse the number of ants from the first element of the game parameters
	antCount, err := strconv.Atoi(gameParams[0])
	if err != nil || antCount <= 0 {
		return nil, serror.ErrInvalidDataOrAntsNum
	}

	// Initialize the farm graph with the number of ants and an empty map of room names
	graph := &entity.FarmGraph{
		AntCount:  antCount,
		RoomNames: make(map[string]*entity.Room),
	}

	// Find the start position of the room descriptions in the game parameters
	startPosition := len(gameParams)

	// Iterate through the game parameters to check and store room descriptions until a dash is encountered
	for i, line := range gameParams[1:] {
		if strings.Contains(line, "-") {
			startPosition = i + 1
			break
		}

		// Check and store the room descriptions in the farm graph
		if err := checkRooms(line, graph); err != nil {
			return nil, err
		}
	}

	// Perform additional checks on the positions of the rooms in the farm graph
	if err := checkPositions(graph); err != nil {
		return nil, err
	}

	// Return error if no start room is found in the game parameters
	if startPosition == len(gameParams) {
		return nil, serror.ErrInvalidDataOrNoStartRoomFound
	}

	// Iterate through the game parameters to parse and store the positions of the rooms in the farm graph
	for _, line := range gameParams[startPosition:] {
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// Parse the position and add it to the farm graph
		pos, err := parsePosition(line, graph)
		if err != nil {
			return nil, serror.ErrInvalidDataOrNoStartRoomFound
		}

		graph.Positions = append(graph.Positions, pos)
	}

	// Return the populated farm graph
	return graph, nil
}

// checkPositions verifies that the graph has a valid start and end room.
func checkPositions(graph *entity.FarmGraph) error {
	// If the graph has no start room, return an error
	if graph.Start == nil {
		return serror.ErrInvalidDataOrNoStartRoomFound
	}

	// If the graph has no end room, return an error
	if graph.End == nil {
		return serror.ErrInvalidDataOrNoEndRoomFound
	}

	// If the graph has both start and end rooms, return no error
	return nil
}

// checkRooms checks the input line and updates the FarmGraph with room details.
// It returns an error if the input line is invalid or if there is a duplicate room.
func checkRooms(line string, f *entity.FarmGraph) (err error) {
	// Create a new Room instance
	room := &entity.Room{}

	// Check if the line starts with "##start"
	switch {
	case strings.HasPrefix(line, "##start"):
		// Return an error if the start room already exists
		if f.Start != nil {
			return serror.ErrInvalidDataFormat // Start room already exists
		}
		// Remove the "##start" prefix and parse the room details
		line = strings.TrimPrefix(line, "##start")
		room, err = parseRoom(line)
		if err != nil {
			return err
		}
		// Set the start room
		f.Start = room
	// Check if the line starts with "##end"
	case strings.HasPrefix(line, "##end"):
		// Return an error if the end room already exists
		if f.End != nil {
			return serror.ErrInvalidDataFormat // End room already exists
		}
		// Remove the "##end" prefix and parse the room details
		line = strings.TrimPrefix(line, "##end")
		room, err = parseRoom(line)
		if err != nil {
			return err
		}
		// Set the end room
		f.End = room
	// Check if the line starts with "#" or is an empty line
	case strings.HasPrefix(line, "#") || line == "":
		// Skip processing for comments and empty lines
		return nil
	default:
		// Parse the room details
		room, err = parseRoom(line)
		if err != nil {
			return err
		}

		// Check if the room name already exists
		if existingPoint := f.RoomNames[room.Name]; existingPoint != nil {
			return serror.ErrInvalidDataFormat // Room name already exists
		}

		// Check for duplicate room coordinates
		for _, p := range f.Rooms {
			if p.X == room.X && p.Y == room.Y {
				return serror.ErrInvalidDataFormat // Duplicate room coordinates
			}
		}
		// Append the new room to the list of rooms
		f.Rooms = append(f.Rooms, room)
	}

	// Set the room name to room mapping
	f.RoomNames[room.Name] = room

	return nil
}

// parseRoom parses the input string and returns a Room object or an error
func parseRoom(str string) (*entity.Room, error) {
	// Split the input string into coordinates
	coords := strings.Fields(str)

	// Check if the number of coordinates is correct
	if len(coords) != 3 {
		return nil, serror.ErrInvalidRoomFormat
	}

	// Extract room name, x coordinate, and y coordinate
	name, xStr, yStr := coords[0], coords[1], coords[2]

	// Convert x coordinate to integer
	x, err := strconv.Atoi(xStr)
	if err != nil {
		return nil, serror.ErrInvalidRoomFormat
	}

	// Convert y coordinate to integer
	y, err := strconv.Atoi(yStr)
	if err != nil {
		return nil, serror.ErrInvalidRoomFormat
	}

	// Check if the room name is valid
	if strings.HasPrefix(name, "L") {
		return nil, serror.ErrInvalidRoomFormat
	}

	// Return the Room object
	return &entity.Room{Name: name, X: x, Y: y}, nil
}

// parsePosition parses a string representing a position and returns a Position object.
func parsePosition(str string, graph *entity.FarmGraph) (*entity.Position, error) {
	// Split the string into 'from' and 'to' parts
	fromTo := strings.Split(str, "-")
	// Check if the string was split into exactly 2 parts
	if len(fromTo) != 2 {
		return nil, serror.ErrInvalidDataFormat
	}
	// Get the 'from' and 'to' room names
	from, to := fromTo[0], fromTo[1]
	// Get the 'from' room object from the map
	fromRoom := graph.RoomNames[from]
	// Check if the 'from' room exists
	if fromRoom == nil {
		return nil, serror.ErrInvalidDataFormat
	}
	// Get the 'to' room object from the map
	toRoom := graph.RoomNames[to]
	// Check if the 'to' room exists
	if toRoom == nil {
		return nil, serror.ErrInvalidDataFormat
	}
	// Return the Position object with the 'from' and 'to' rooms
	return &entity.Position{Start: fromRoom, End: toRoom}, nil
}

// SearchPaths is a recursive function that searches for valid paths in a farm graph.
// It takes a list of paths, a group of paths, a farm graph, and a list of path groups as input.
func SearchPaths(paths []*entity.Path, group []*entity.Path, graph *entity.FarmGraph, groups *[][]*entity.Path) {
	// If the group length equals the number of ants, add the group to the list of groups and return
	if len(group) == graph.AntCount {
		*groups = append(*groups, group)
		return
	}
	// Iterate through the paths
	for i, path := range paths {
		// Initialize a variable to track if the path is already taken
		pathIsTaken := false
		// Iterate through the paths in the group
		for _, takenPath := range group {
			// Iterate through the points in the current path
			for _, room := range *path {
				// Iterate through the points in the taken path, excluding the start and end points
				for _, takenPoint := range (*takenPath)[1 : len(*takenPath)-1] {
					// If the current point is the same as a point in the taken path, mark the path as taken
					if room == takenPoint {
						pathIsTaken = true
						break
					}
				}
				// If the path is already taken, break out of the loop
				if pathIsTaken {
					break
				}
			}
			// If the path is already taken, break out of the loop
			if pathIsTaken {
				break
			}
		}
		// If the path is not taken, recursively call SearchPaths with the updated paths and group
		if !pathIsTaken {
			newPaths := append(paths[:i], paths[i+1:]...)  // Create a new list of paths without the current path
			newGroup := append(group, path)                // Add the current path to the group
			SearchPaths(newPaths, newGroup, graph, groups) // Recursively call SearchPaths with the new paths and group
		}
	}
	// If the group is not empty, add the group to the list of groups
	if len(group) != 0 {
		*groups = append(*groups, group)
	}
}

// MoveAntsToEnd returns the optimal solution path for the FarmGraph
func MoveAntsToEnd(graph *entity.FarmGraph) []*entity.Path {
	// allPathGroups stores all the possible combinations of paths
	allPathGroups := [][]*entity.Path{}

	// If the first path has only 2 nodes, start searching from the second path
	if len(*graph.Paths[0]) == 2 {
		SearchPaths(graph.Paths[1:], []*entity.Path{graph.Paths[0]}, graph, &allPathGroups)
	} else {
		SearchPaths(graph.Paths, []*entity.Path{}, graph, &allPathGroups)
	}

	// paths stores the optimal path
	paths := []*entity.Path{}
	pathsTurns := math.MaxInt

	// Calculate the turns for each group of paths and find the optimal paths
	for _, group := range allPathGroups {
		bandwidth := 1
		for _, path := range group {
			bandwidth *= len(*path)
		}
		multiSum := 0
		for _, path := range group {
			multiSum += bandwidth / len(*path)
		}
		turns := graph.AntCount * bandwidth / multiSum
		if turns < pathsTurns {
			pathsTurns = turns
			paths = group
		}
	}

	// Return the optimal paths
	return paths
}
