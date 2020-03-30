package goteamsnotify

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

// FormatAsCodeBlock accepts an arbitrary string, quoted or not, and calls a
// helper function which attempts to format as a valid Markdown code block for
// submission to Microsoft Teams
func FormatAsCodeBlock(input string) (string, error) {

	if input == "" {
		return "", errors.New("received empty string, refusing to format")
	}

	result, err := formatAsCode(
		input,
		msTeamsCodeBlockSubmissionPrefix,
		msTeamsCodeBlockSubmissionSuffix,
	)

	return result, err

}

// FormatAsCodeSnippet accepts an arbitrary string, quoted or not, and calls a
// helper function which attempts to format as a single-line valid Markdown
// code snippet for submission to Microsoft Teams
func FormatAsCodeSnippet(input string) (string, error) {
	if input == "" {
		return "", errors.New("received empty string, refusing to format")
	}

	result, err := formatAsCode(
		input,
		msTeamsCodeSnippetSubmissionPrefix,
		msTeamsCodeSnippetSubmissionSuffix,
	)

	return result, err
}

// formatAsCode is a helper function which accepts an arbitrary string, quoted
// or not, a desired prefix and a suffix for the string and attempts to format
// as a valid Markdown formatted code sample for submission to Microsoft Teams
func formatAsCode(input string, prefix string, suffix string) (string, error) {

	// required; protects against slice out of range panics
	if input == "" {
		return "", errors.New("received empty string, refusing to format as code block")
	}

	logger.Printf("Calling json.Marshal(input); input: %+v", input)
	byteSlice, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	logger.Println("byteSlice as string:", string(byteSlice))

	var prettyJSON bytes.Buffer

	logger.Println("calling json.Indent")
	err = json.Indent(&prettyJSON, []byte(input), "", "\t")
	if err != nil {
		return "", err
	}
	formattedJSON := prettyJSON.String()

	logger.Println("Formatted JSON:", formattedJSON)

	var codeContentForSubmission string

	// handle cases where the formatted JSON string was not wrapped with
	// double-quotes
	switch {

	// if neither start or end character are double-quotes
	case string(formattedJSON[0]) != `"` && string(formattedJSON[len(formattedJSON)-1]) != `"`:
		codeContentForSubmission = prefix + string(formattedJSON) + suffix

	// if only start character is not a double-quote
	case string(formattedJSON[0]) != `"`:
		logger.Println("[WARN]: escapedFormattedJSON is missing leading double-quote")
		codeContentForSubmission = prefix + string(formattedJSON)

	// if only end character is not a double-quote
	case string(formattedJSON[len(formattedJSON)-1]) != `"`:
		logger.Println("[WARN]: escapedFormattedJSON is missing trailing double-quote")
		codeContentForSubmission = codeContentForSubmission + suffix

	default:
		// Guard against strings of length 1 to prevent out of range panics:
		// panic: runtime error: slice bounds out of range [1:0]
		minLength := 2
		if len(formattedJSON) < minLength {
			return "", fmt.Errorf(
				"formattedJSON is invalid length; got %d chars, want at least %d chars",
				len(formattedJSON),
				minLength,
			)
		}
		codeContentForSubmission = prefix + string(formattedJSON[1:len(formattedJSON)-1]) + suffix
	}

	logger.Printf("... as-is:\n%s\n\n", string(formattedJSON))
	logger.Printf("... without leading and trailing double-quotes: \n%s\n\n", string(formattedJSON[1:len(formattedJSON)-1]))
	logger.Printf("codeContentForSubmission: \n%s\n\n", codeContentForSubmission)

	// err should be nil if everything worked as expected
	return codeContentForSubmission, err

}
