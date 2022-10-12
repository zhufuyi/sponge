package configs

import "testing"

func TestPath(t *testing.T) {
	p := Path(".")
	t.Log(p)
}
