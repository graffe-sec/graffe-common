package state

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type Recommendation struct {
	severity int
	category string
	action   string
}

type Result struct {
	assetID         string
	recommendations []Recommendation
}

type Scan struct {
	region  string
	results []Result
}

type State struct {
	timestamp string
	scans     []Scan
}

const stateFolderPath = ".\\.state"

// WriteState : writes a new time based state file
func WriteState(accountID string, state *State) error {
	output, err := json.MarshalIndent(state, "", " ")
	if err != nil {
		return fmt.Errorf("Issue convertion state input to JSON: %s", err)
	}

	// Create the state folder if it doesn't already exist.
	_ = os.Mkdir(stateFolderPath, os.ModeDir)

	// Construct the output file path.
	var outputFilePath string = fmt.Sprintf(
		".\\%s\\%s.json",
		stateFolderPath,
		accountID,
	)

	// Archive the latest state file if one exists.
	source, err := os.Open(outputFilePath)
	if err == nil {
		// Read the contents of the state file.
		byteValue, _ := ioutil.ReadAll(source)

		// Unmarshal the state file into a State struct.
		var oldState State
		json.Unmarshal(byteValue, &oldState)

		// Extract the timestamp of the file to be used in the name of the archived state.
		var archiveFilePath string = fmt.Sprintf(
			".\\%s\\%s-%s.json",
			stateFolderPath,
			accountID,
			oldState.timestamp,
		)

		// Create the new file to archive to.
		destination, err := os.Create(archiveFilePath)
		if err != nil {
			return fmt.Errorf("Issue creating archived state <%s>: %s", archiveFilePath, err)
		}

		// Archive the state.
		_, err = io.Copy(destination, source)
		if err != nil {
			return fmt.Errorf("Issue archiving state <%s>: %s", archiveFilePath, err)
		}
	}

	// Write the supplied state to the destination.
	err = ioutil.WriteFile(outputFilePath, output, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
