package sink

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/index.html", r.URL.Path)
		assert.Equal(t, r.Header["X-My-Header"], []string{"value1", "value2"})
		assert.Equal(t, r.Header["X-Other-Header"], []string{"value3"})
	}))

	mod, err := NewHTTP(mirror.ModuleContext{}, []byte(`{"target_url": "`+server.URL+`", "timeout": "10s"}`))
	require.NoError(t, err)

	in := make(chan mirror.Request, 1)
	in <- mirror.Request{
		Method: mirror.GET,
		Path:   "/index.html",
		Headers: map[string]mirror.HeaderValues{
			"X-My-Header":    {Values: []string{"value1", "value2"}},
			"X-Other-Header": {Values: []string{"value3"}},
		},
	}

	mod.SetInput(in)
	close(in)

	<-mod.Output()
}
