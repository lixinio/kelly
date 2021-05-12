package swagger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFileNode(t *testing.T) {
	testcases := []struct {
		input  string
		expect struct {
			filename string
			path     string
		}
	}{
		{
			input: "swagger.yaml:test",
			expect: struct {
				filename string
				path     string
			}{
				filename: "swagger.yaml",
				path:     "test",
			},
		},
	}

	for _, testcase := range testcases {
		filename, path, err := parseFileNode(testcase.input)
		require.Equal(t, nil, err)
		require.Equal(t, testcase.expect.filename, filename)
		require.Equal(t, testcase.expect.path, path)
	}
}

func TestPathEditor(t *testing.T) {
	editor := newPathEditor()
	testcases := []struct {
		input  string
		expect string
	}{
		{
			input:  "/aaa/bb/:cc/:dd/ee",
			expect: "/aaa/bb/{cc}/{dd}/ee",
		},
	}

	for _, testcase := range testcases {
		require.Equal(t, testcase.expect, editor.update(testcase.input))
	}
}

func TestTagOptionsContains(t *testing.T) {
	testcases := []struct {
		input  string
		key    string
		expect bool
	}{
		{
			input:  "abc,def",
			key:    "abc",
			expect: true,
		},
		{
			input:  "abc,def",
			key:    "abcd",
			expect: false,
		},
		{
			input:  "abc,def",
			key:    "0abc",
			expect: false,
		},
		{
			input:  "abc,def",
			key:    "abc,",
			expect: false,
		},
		{
			input:  "abc,def",
			key:    "def",
			expect: true,
		},
		{
			input:  "abc",
			key:    "abc",
			expect: true,
		},
		{
			input:  "abc,",
			key:    "abc",
			expect: true,
		},
		{
			input:  "abc,def,",
			key:    "abc",
			expect: true,
		},
		{
			input:  ",abc,def,",
			key:    "abc",
			expect: true,
		},
	}

	for _, testcase := range testcases {
		to := tagOptions(testcase.input)
		require.Equal(t, testcase.expect, to.contains(testcase.key))
	}
}

func TestTagOptionsGetValue(t *testing.T) {
	testcases := []struct {
		input  string
		key    string
		expect string
		result bool
	}{
		{
			input:  "abc,def",
			key:    "abc",
			expect: "",
			result: false,
		},
		{
			input:  "abc=def",
			key:    "abc",
			expect: "def",
			result: true,
		},
		{
			input:  "abc=def,hij",
			key:    "abc",
			expect: "def",
			result: true,
		},
		{
			input:  "abc=def,hij",
			key:    "hij",
			expect: "",
			result: false,
		},
	}

	for _, testcase := range testcases {
		to := tagOptions(testcase.input)
		name, exists := to.getValue(testcase.key)
		require.Equal(t, testcase.expect, name)
		require.Equal(t, testcase.result, exists)
	}
}
