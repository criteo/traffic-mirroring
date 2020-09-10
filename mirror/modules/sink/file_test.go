package sink

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/stretchr/testify/require"
)

func TestFileJSON(t *testing.T) {
	f, err := ioutil.TempFile("/tmp", "http-mirror-test")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	mod, err := NewFile(mirror.ModuleContext{}, []byte(`{"path": "`+f.Name()+`", "format": "json"}`))
	require.NoError(t, err)

	reqs := []mirror.Request{
		{
			Method:      mirror.Method_GET,
			Path:        "/index.html",
			HttpVersion: mirror.HTTPVersion_HTTP1_1,
			Headers: map[string]*mirror.HeaderValue{
				"Host": {Values: []string{"127.0.0.1:10080"}},
			},
		},
		{
			Method:      mirror.Method_GET,
			Path:        "/index.css",
			HttpVersion: mirror.HTTPVersion_HTTP1_1,
			Headers: map[string]*mirror.HeaderValue{
				"Host": {Values: []string{"127.0.0.1:10080"}},
			},
		},
	}

	in := make(chan mirror.Request, len(reqs))
	for _, r := range reqs {
		in <- r
	}

	mod.SetInput(in)
	close(in)

	<-mod.Output()

	contents, err := ioutil.ReadFile(f.Name())
	require.NoError(t, err)

	parts := bytes.Split(contents, []byte{'\n'})
	require.Equal(t, len(reqs), len(parts)-1)
	parts = parts[:len(parts)-1]

	for i, part := range parts {
		req := mirror.Request{}
		err := json.Unmarshal(part, &req)
		require.NoError(t, err)
		require.Equal(t, reqs[i], req)
	}
}
