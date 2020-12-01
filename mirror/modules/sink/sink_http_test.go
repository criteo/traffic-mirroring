package sink

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/criteo/traffic-mirroring/mirror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTP(t *testing.T) {
	serv_count := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/index.html", r.URL.Path)
		assert.Equal(t, r.Header["X-My-Header"], []string{"value1", "value2"})
		assert.Equal(t, r.Header["X-Other-Header"], []string{"value33"})
		serv_count ++
	}))

	mod, err := NewHTTP(&mirror.ModuleContext{}, []byte(`{"target_url": "`+server.URL+`", "timeout": "10s"}`))
	require.NoError(t, err)

	in := make(chan mirror.Request, 1)
	in <- mirror.Request{
		Method: mirror.Method_GET,
		Path:   "/index.html",
		Headers: map[string]*mirror.HeaderValue{
			"X-My-Header":    {Values: []string{"value1", "value2"}},
			"X-Other-Header": {Values: []string{"value3"}},
		},
	}

	mod.SetInput(in)
	close(in)

	<-mod.Output()

	assert.Equal(t, serv_count, 0)
}
