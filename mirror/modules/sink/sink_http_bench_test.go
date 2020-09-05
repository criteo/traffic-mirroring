package sink

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/stretchr/testify/require"
)

func BenchmarkHTTP(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {}))

	mod, err := NewHTTP([]byte(`{
		"target_url": "` + server.URL + `",
		"timeout": "10s",
		"parallel": 10
	}`))
	require.NoError(b, err)
	in := make(chan mirror.Request)
	mod.SetInput(in)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		in <- mirror.Request{
			Method: mirror.GET,
			Path:   "/index.html",
			Headers: map[string]mirror.HeaderValues{
				"X-My-Header":    {Values: []string{"value1", "value2"}},
				"X-Other-Header": {Values: []string{"value3"}},
			},
		}
	}

	close(in)

	<-mod.Output()
}
