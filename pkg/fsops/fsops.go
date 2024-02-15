package fsops

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"lem-in/pkg/serror"
)

// GetOSArgsAndInterpreter reads the command line arguments and extracts the file name from the second argument
// If there are not enough command line arguments, it returns an error ErrNotEnoughOSArgs
func GetOSArgsAndInterpreter() (fileName string, err error) {
	// Get the command line arguments
	osArgs := os.Args

	// Check if there are enough arguments
	if len(osArgs) <= 1 {
		// If not enough arguments, return the error ErrNotEnoughOSArgs
		return "", serror.ErrNotEnoughOSArgs
	}

	// Return the second command line argument as the file name and no error
	return osArgs[1], nil
}

// ReadAndPrintFileContent reads the file content from the file name obtained from GetOSArgsAndInterpreter
// and prints it to Stdout. It returns an error if there are any issues with the file or the reading process.
func ReadAndPrintFileContent() error {
	// Get the file name and any error from the command line arguments
	fileName, err := GetOSArgsAndInterpreter()
	if err != nil {
		return err
	}

	// Open the file for reading
	file, err := os.Open(fileName)
	if err != nil {
		return serror.ErrInvalidFileName // Return an error if the file cannot be opened
	}
	defer file.Close() // Ensure the file is closed after reading

	// Create a new scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Iterate through each line in the file and print it to Stdout
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// Check for any errors that occurred during the file reading process
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil // Return nil if the file was read successfully
}

// ReadFile and return []byte of the file content
// If something going wrong, rerurn ErrInvalidFileName, ErrNotEnoughOSArgs
func ReadFile() (data []byte, err error) {
	fileName, err := GetOSArgsAndInterpreter()
	if err != nil {
		return data, err
	}

	data, err = os.ReadFile(fileName)
	if err != nil {
		return nil, serror.ErrInvalidFileName
	}

	return data, nil
}

// ParseFile parses the input data and returns an array of strings after performing various string manipulations.
func ParseFile(data []byte) ([]string, error) {
	// Remove carriage returns
	line := strings.ReplaceAll(string(data), "\r", "")

	// Replace double newlines with single newline
	line = strings.ReplaceAll(line, "\n\n", "\n")

	// Remove multiple spaces with a single space
	regexSpaces := regexp.MustCompile(` +`)
	line = regexSpaces.ReplaceAllLiteralString(line, " ")

	// Trim leading spaces
	regexLeadingSpaces := regexp.MustCompile(`^ +`)
	line = regexLeadingSpaces.ReplaceAllLiteralString(line, "")

	// Remove space before newline
	line = strings.ReplaceAll(line, " \n", "\n")

	// Replace specific patterns
	line = strings.NewReplacer("##start\n", "##start", "##end\n", "##end").Replace(line)

	// Split the string into an array of lines
	lines := strings.Split(line, "\n")

	// Check if the number of lines is less than 4
	if len(lines) < 4 {
		return nil, serror.ErrParseGameParamFile
	}

	// Return the array of strings
	return lines, nil
}
