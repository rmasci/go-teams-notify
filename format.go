package goteamsnotify

import (
	"bytes"
	"encoding/json"
	"errors"
)

// Even though Microsoft Teams doesn't show the additional newlines,
// https://messagecardplayground.azurewebsites.net/ DOES show the results
// as a formatted code block. Including the newlines now is an attempt at
// "future proofing" the codeblock support in MessageCard values sent to
// Microsoft Teams.
const (

	// msTeamsCodeBlockSubmissionPrefix is the prefix appended to text input
	// to indicate that the text should be displayed as a codeblock by
	// Microsoft Teams.
	msTeamsCodeBlockSubmissionPrefix string = "\n```\n"
	// msTeamsCodeBlockSubmissionPrefix string = "```"

	// msTeamsCodeBlockSubmissionSuffix is the suffix appended to text input
	// to indicate that the text should be displayed as a codeblock by
	// Microsoft Teams.
	msTeamsCodeBlockSubmissionSuffix string = "```\n"
	// msTeamsCodeBlockSubmissionSuffix string = "```"

	// msTeamsCodeSnippetSubmissionPrefix is the prefix appended to text input
	// to indicate that the text should be displayed as a code formatted
	// string of text by Microsoft Teams.
	msTeamsCodeSnippetSubmissionPrefix string = "`"

	// msTeamsCodeSnippetSubmissionSuffix is the suffix appended to text input
	// to indicate that the text should be displayed as a code formatted
	// string of text by Microsoft Teams.
	msTeamsCodeSnippetSubmissionSuffix string = "`"
)

// TryToFormatAsCodeBlock acts as a wrapper for FormatAsCodeBlock. If an
// error is encountered in the FormatAsCodeBlock function, this function will
// return the original string, otherwise if no errors occur the newly formatted
// string will be returned.
func TryToFormatAsCodeBlock(input string) string {

	result, err := FormatAsCodeBlock(input)
	if err != nil {
		logger.Printf("TryToFormatAsCodeBlock: error occurred when calling FormatAsCodeBlock: %v\n", err)
		logger.Println("TryToFormatAsCodeBlock: returning original string")
		return input
	}

	logger.Println("TryToFormatAsCodeBlock: no errors occurred when calling FormatAsCodeBlock")
	return result
}

// TryToFormatAsCodeSnippet acts as a wrapper for FormatAsCodeSnippet. If
// an error is encountered in the FormatAsCodeSnippet function, this function will
// return the original string, otherwise if no errors occur the newly formatted
// string will be returned.
func TryToFormatAsCodeSnippet(input string) string {

	result, err := FormatAsCodeSnippet(input)
	if err != nil {
		logger.Printf("TryToFormatAsCodeSnippet: error occurred when calling FormatAsCodeBlock: %v\n", err)
		logger.Println("TryToFormatAsCodeSnippet: returning original string")
		return input
	}

	logger.Println("TryToFormatAsCodeSnippet: no errors occurred when calling FormatAsCodeSnippet")
	return result
}

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

	var err error
	var byteSlice []byte

	switch {

	// required; protects against slice out of range panics
	case input == "":
		return "", errors.New("received empty string, refusing to format as code block")

	// If the input string is already valid JSON, don't double-encode and
	// escape the content
	case json.Valid([]byte(input)):
		logger.Printf("DEBUG: input string already valid JSON; input: %+v", input)
		logger.Printf("DEBUG: Calling json.RawMessage([]byte(input)); input: %+v", input)

		// FIXME: Is json.RawMessage() really needed if the input string is *already* JSON?
		// https://golang.org/pkg/encoding/json/#RawMessage seems to imply a different use case.
		byteSlice = json.RawMessage([]byte(input))
		//
		// From light testing, it appears to not be necessary:
		//
		// logger.Printf("DEBUG: Skipping json.RawMessage, converting string directly to byte slice; input: %+v", input)
		// byteSlice = []byte(input)

	default:
		logger.Printf("DEBUG: input string not valid JSON; input: %+v", input)
		logger.Printf("DEBUG: Calling json.Marshal(input); input: %+v", input)
		byteSlice, err = json.Marshal(input)
		if err != nil {
			return "", err
		}
	}

	logger.Println("DEBUG: byteSlice as string:", string(byteSlice))

	var prettyJSON bytes.Buffer

	logger.Println("DEBUG: calling json.Indent")
	err = json.Indent(&prettyJSON, byteSlice, "", "\t")
	if err != nil {
		return "", err
	}
	formattedJSON := prettyJSON.String()

	logger.Println("DEBUG: Formatted JSON:", formattedJSON)

	var codeContentForSubmission string

	// try to prevent "runtime error: slice bounds out of range"
	formattedJSONStartChar := 0
	formattedJSONEndChar := len(formattedJSON) - 1
	if formattedJSONEndChar < 0 {
		formattedJSONEndChar = 0
	}

	// handle cases where the formatted JSON string was not wrapped with
	// double-quotes
	switch {

	// if neither start or end character are double-quotes
	case string(formattedJSON[formattedJSONStartChar]) != `"` && string(formattedJSON[formattedJSONEndChar]) != `"`:
		codeContentForSubmission = prefix + string(formattedJSON) + suffix

	// if only start character is not a double-quote
	case string(formattedJSON[formattedJSONStartChar]) != `"`:
		logger.Println("[WARN]: escapedFormattedJSON is missing leading double-quote")
		codeContentForSubmission = prefix + string(formattedJSON)

	// if only end character is not a double-quote
	case string(formattedJSON[formattedJSONEndChar]) != `"`:
		logger.Println("[WARN]: escapedFormattedJSON is missing trailing double-quote")
		codeContentForSubmission = codeContentForSubmission + suffix

	default:
		codeContentForSubmission = prefix + string(formattedJSON[1:formattedJSONEndChar]) + suffix
	}

	logger.Printf("DEBUG: ... as-is:\n%s\n\n", string(formattedJSON))

	// this requires that the formattedJSON be at least two characters long
	if len(formattedJSON) > 2 {
		logger.Printf("DEBUG: ... without first and last characters: \n%s\n\n", string(formattedJSON[formattedJSONStartChar+1:formattedJSONEndChar]))
	}
	logger.Printf("DEBUG: formattedJSON is less than two chars: \n%s\n\n", string(formattedJSON))

	logger.Printf("DEBUG: codeContentForSubmission: \n%s\n\n", codeContentForSubmission)

	// err should be nil if everything worked as expected
	return codeContentForSubmission, err

}
