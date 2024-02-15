package core

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"lem-in/internal/entity"
)

// ParseAdadjacents initializes the adjacents map in entity.FarmGraph
func ParseAdadjacents(graph *entity.FarmGraph) {
	// Initialize the map
	graph.Adadjacents = make(map[*entity.Room][]*entity.Room)
	// Iterate through positions to create adjacency list
	for _, position := range graph.Positions {
		// Add end position to the adjacency list of start position
		graph.Adadjacents[position.Start] = append(graph.Adadjacents[position.Start], position.End)
		// Add start position to the adjacency list of end position
		graph.Adadjacents[position.End] = append(graph.Adadjacents[position.End], position.Start)
	}
}

// DepthSearch performs a depth search on the entity.FarmGraph to find paths from a start room to the end room
func DepthSearch(graph *entity.FarmGraph, currentRoom *entity.Room, visitedRooms map[*entity.Room]bool, lastVisitedRoom map[*entity.Room]*entity.Room) {
	// If the current room is the end room, construct the path and add it to the graph's list of paths
	if currentRoom == graph.End {
		path := entity.Path{graph.End}
		for previousRoom := lastVisitedRoom[currentRoom]; previousRoom != nil; previousRoom = lastVisitedRoom[previousRoom] {
			path = append(path, previousRoom)
		}
		// Reverse the path
		for i := 0; i < len(path)/2; i++ {
			path[i], path[len(path)-1-i] = path[len(path)-1-i], path[i]
		}
		graph.Paths = append(graph.Paths, &path)
		return
	}
	// Mark the current room as visited
	visitedRooms[currentRoom] = true
	// Recursively visit each unvisited adjacent room
	for _, childRoom := range graph.Adadjacents[currentRoom] {
		if !visitedRooms[childRoom] {
			lastVisitedRoom[childRoom] = currentRoom
			DepthSearch(graph, childRoom, visitedRooms, lastVisitedRoom)
		}
	}
	// Mark the current room as not visited after all its adjacent rooms have been visited
	visitedRooms[currentRoom] = false
}

// FindAllPaths finds all possible paths in the farm graph
func FindAllPaths(graph *entity.FarmGraph) error {
	// Create a map to keep track of visited rooms
	roomVisited := make(map[*entity.Room]bool)

	// Create a map to keep track of the last visited room
	lastRoom := make(map[*entity.Room]*entity.Room)

	// Parse the adjacent rooms and store them in the farm graph
	ParseAdadjacents(graph)

	// Search for all possible paths starting from the farm graph's starting point
	DepthSearch(graph, graph.Start, roomVisited, lastRoom)

	// Sort the paths based on their length
	sort.Slice(graph.Paths, func(i, j int) bool {
		return len(*graph.Paths[i]) < len(*graph.Paths[j])
	})

	// If no paths were found, return an error
	if len(graph.Paths) == 0 {
		return fmt.Errorf("no paths")
	}

	// Otherwise, return nil to indicate success
	return nil
}

// ResultOutput takes a solution of paths and returns a string representation
func ResultOutput(solution []*entity.Path, graph *entity.FarmGraph) string {
	// Initialize FarmGraph and paths
	var paths = make(map[int][]string)

	// Iterate through the ants and their paths
	for ant := 0; ant < graph.AntCount; ant++ {
		// Get the path of the current ant
		antPath := solution[ant%len(solution)]
		// Calculate the starting line number for the current ant
		pathStartLine := ant/len(solution) + 1

		// Find the shortest path for the current ant
		if graph.AntCount-len(solution) < ant {
			for antNew := ant + 1; antNew < ant+len(solution); antNew++ {
				// Get the path of the new ant
				newAntPath := solution[antNew%len(solution)]
				// Calculate the starting line number for the new ant
				newPathStartLine := antNew/len(solution) + 1
				// Compare the lengths of the paths and update the current ant's path if the new path is shorter
				if len(*newAntPath)+newPathStartLine < len(*antPath)+pathStartLine {
					antPath = newAntPath
					pathStartLine = newPathStartLine
				}
			}
		}

		// Add the ant's path to the paths map
		for i, room := range (*antPath)[1:] {
			// Generate the log entry for the current ant's movement and add it to the paths map
			paths[pathStartLine+i] = append(paths[pathStartLine+i], "L"+strconv.Itoa(ant+1)+"-"+room.Name)
		}
	}

	// Convert paths map to a slice of strings
	var pathsSlice []string
	for i := 1; i <= len(paths); i++ {
		// Join the log entries for each line and add it to the pathsSlice
		pathsSlice = append(pathsSlice, strings.Join(paths[i], " "))
	}

	// Join the paths and return as a string
	return strings.Join(pathsSlice, "\n")
}
