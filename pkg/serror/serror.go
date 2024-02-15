package serror

import (
	"errors"
	"fmt"
	"log"
)

var (
	ErrNotEnoughOSArgs    = errors.New("ERROR: not enough os args")
	ErrInvalidFileName    = errors.New("ERROR: invalid a file name")
	ErrParseGameParamFile = errors.New("ERROR: while parsing the game params file")
	//
	ErrInvalidDataFormat             = errors.New("ERROR: invalid data format")
	ErrInvalidRoomFormat             = errors.New("ERROR: invalid room format")
	ErrInvalidDataOrAntsNum          = errors.New("ERROR: invalid data format, invalid number of Ants")
	ErrInvalidDataOrNoStartRoomFound = errors.New("ERROR: invalid data format, no start room found")
	ErrInvalidDataOrNoEndRoomFound   = errors.New("ERROR: invalid data format, no end room found")
)

// ErrorHandler handling error with specefic way
func ErrorHandler(err error) {
	if err != nil {
		switch {
		case errors.Is(err, ErrNotEnoughOSArgs):
			fmt.Println("Usage: go run main.go <exampleNum.txt>")
		case errors.Is(err, ErrInvalidFileName):
			fmt.Println("Wrong a file name. Please, try again!")
		default:
			log.Fatal(err)
		}
	}
}
