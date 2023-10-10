package jira

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		name     string
		given    string
		expected string
	}{
		{
			name:     "Hello World",
			given:    "hello-world.json",
			expected: "hello-world.html",
		},
		{
			name:     "ADF",
			given:    "adf.json",
			expected: "adf.html",
		},
	}

	// Read the file
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			givenFilePath := fmt.Sprintf("testdata/%s", tc.given)
			jsonData, err := os.ReadFile(givenFilePath)

			if err != nil {
				t.Errorf("Failed to read given JSON data: %v", err)
			}

			expectedFilePath := fmt.Sprintf("testdata/%s", tc.expected)
			expectedData, err := os.ReadFile(expectedFilePath)
			expectedSummary := string(expectedData)

			if err != nil {
				t.Errorf("Failed to read expected data: %v", err)
			}

			// Unmarshal JSON data into the Document struct
			var doc Description
			err = json.Unmarshal([]byte(jsonData), &doc)
			if err != nil {
				t.Errorf("Error unmarshaling JSON: %v", err)
			}

			// Call the SummarizeAndValidate method
			summary := doc.String()
			if err != nil {
				t.Errorf("Error in SummarizeAndValidate: %v", err)
			}

			// Check if the summarized result matches the expected summary
			if summary != expectedSummary {
				t.Errorf("Summarized output is incorrect.\nGot: %s\nExpected: %s", summary, expectedSummary)
			}
		})
	}
}
