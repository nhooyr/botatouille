package command

import (
	"testing"
	"github.com/nhooyr/botatouille/digo/commandmmand"
)

func TestIs(t *testing.T) {
	t.Run("foo", func(t *testing.T) {
		_, ok := command.Is("foo")
		if ok {
			t.Error("expected not ok")
		}
	})
	t.Run("!foo", func(t *testing.T) {
		_, ok := command.Is("!foo")
		if !ok {
			t.Error("expected ok")
		}
	})
}