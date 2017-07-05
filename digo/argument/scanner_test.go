package argument

import (
	"testing"
	"io"
	"github.com/nhooyr/botatouille/digo/argumentument"
)

func TestScanner(t *testing.T) {
	type testCase struct {
		name    string
		cmdLine string
		token   string
		rest    string
		err     error
	}
	testCases := []testCase{
		{"quotes", `"food" rest`, "food", "rest", nil},
		{"single quotes", `'food' rest`, "food", "rest", nil},
		{"spaces in quotes", `"foo bar d" rest`, "foo bar d", "rest", nil},
		{"multiple quotes", `xd"foo bar d"xd'foobar lol' rest`, "xdfoo bar dxdfoobar lol", "rest", nil},
		{"backspaces", `\s"\m\"\'" rest`, "sm\"'", "rest", nil},
		{"unexpected backspace", `lol\`, "lol", "", argument.ErrUnexpectedBackspace},
		{"EOF", "", "", "", io.EOF},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := argument.NewScanner(tc.cmdLine)
			err := s.Scan()
			if err != tc.err {
				t.Fatalf("expected %q but got %q", tc.err, err)
			}
			token := s.Token()
			if token != tc.token {
				t.Errorf("expected %q but got %q", tc.token, token)
			}
			rest := s.Rest()
			if rest != tc.rest {
				t.Errorf("expected %q but got %q", tc.rest, rest)
			}
		})
	}
}
