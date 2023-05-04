// Copyright 2020-2023 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package bufcurl

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestErrorHasFilename(t *testing.T) {
	err := errors.New("example error")
	filename := "example.txt"

	// Test case where filename is not included in the error message.
	newErr := ErrorHasFilename(err, filename)
	expectedErr := fmt.Sprintf("%s: %v", filename, err)
	if newErr.Error() != expectedErr {
		t.Errorf("ErrorHasFilename(%v, %s) = %s, want %s", err, filename, newErr, expectedErr)
	}

	// Test case where filename is already included in the error message.
	err = fmt.Errorf("%s: %v", filename, err)
	newErr = ErrorHasFilename(err, filename)
	if newErr.Error() != err.Error() {
		t.Errorf("ErrorHasFilename(%v, %s) = %s, want %s", err, filename, newErr, err)
	}
}

func TestLineReader(t *testing.T) {
	type testData struct {
		input    string
		expected []string
	}

	tests := []testData{
		{
			input: "line 1\nline 2\nline 3\n",
			expected: []string{
				"line 1",
				"line 2",
				"line 3",
			},
		},
		{
			input:    "line 1\n",
			expected: []string{"line 1"},
		},
		{
			input:    "",
			expected: nil,
		},
	}

	for i, test := range tests {
		r := &lineReader{r: bufio.NewReader(strings.NewReader(test.input))}
		for _, expected := range test.expected {
			line, err := r.ReadLine()
			if err != nil {
				t.Errorf("TestLineReader(%d): unexpected error: %v", i, err)
			}
			if line != expected {
				t.Errorf("TestLineReader(%d): got %s, want %s", i, line, expected)
			}
		}

		// Test that EOF is returned after the last line.
		line, err := r.ReadLine()
		if err != io.EOF {
			t.Errorf("TestLineReader(%d): Unexpected error: %v", i, err)
		}
		if line != "" {
			t.Errorf("TestLineReader(%d): Unexpected result: %v", i, line)
		}
	}
}
