package control

import (
	"testing"

	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/stretchr/testify/require"
)

func TestFanout(t *testing.T) {
	mod, err := NewFanout(mirror.ModuleContext{}, []byte(`[
		{"type": "control.identity", "config": {}},
		{"type": "control.identity", "config": {}}
	]`))
	require.NoError(t, err)

	expected := mirror.Request{Path: "test"}
	in := make(chan mirror.Request, 1)
	in <- expected

	mod.SetInput(in)
	close(in)

	out := []mirror.Request{}
	for r := range mod.Output() {
		out = append(out, r)
	}

	require.Equal(t, []mirror.Request{expected, expected}, out)
}
