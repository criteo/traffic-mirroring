package expr

import (
	"encoding/json"
	"testing"

	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/stretchr/testify/require"
)

func TestJSONString(t *testing.T) {
	e := &StringExpr{}
	err := json.Unmarshal([]byte(`"str"`), e)
	require.NoError(t, err)

	v, err := e.Eval(mirror.Request{})
	require.NoError(t, err)

	require.Equal(t, "str", v)
}

func TestJSONStringBad(t *testing.T) {
	e := &StringExpr{}
	err := json.Unmarshal([]byte(`"{1}"`), e)
	require.NotNil(t, err)
}
