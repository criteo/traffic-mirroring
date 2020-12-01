package sink

import (
	"github.com/criteo/traffic-mirroring/mirror"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTP(t *testing.T) {
	reqCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/index.html", r.URL.Path)
		assert.Equal(t, r.Header["X-My-Header"], []string{"value1", "value2"})
		assert.Equal(t, r.Header["X-Other-Header"], []string{"value3"})
		reqCount++
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

	require.Equal(t, 1, reqCount)
}
func TestHTTP_dynamic_target(t *testing.T) {
	reqCount1 := 0
	server1 := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/index.html", r.URL.Path)
		assert.Equal(t, r.Header["X-My-Header"], []string{"value1_1", "value1_2"})
		assert.Equal(t, r.Header["X-Other-Header"], []string{"value1_3"})
		reqCount1++
	}))

	reqCount2 := 0
	server2 := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/index.html", r.URL.Path)
		assert.Equal(t, r.Header["X-My-Header"], []string{"value2_1", "value2_2"})
		assert.Equal(t, r.Header["X-Other-Header"], []string{"value2_3"})
		reqCount2++
	}))

	mod, err := NewHTTP(&mirror.ModuleContext{}, []byte(`{"target_url": "{req.meta.target.string}", "timeout": "10s"}`))
	require.NoError(t, err)

	in := make(chan mirror.Request, 2)
	in <- mirror.Request{
		Method: mirror.Method_GET,
		Path:   "/index.html",
		Meta: map[string]*mirror.MetaValue{
			"target": {Value: &mirror.MetaValue_String_{String_: server1.URL}},
		},
		Headers: map[string]*mirror.HeaderValue{
			"X-My-Header":    {Values: []string{"value1_1", "value1_2"}},
			"X-Other-Header": {Values: []string{"value1_3"}},
		},
	}
	in <- mirror.Request{
		Method: mirror.Method_GET,
		Path:   "/index.html",
		Meta: map[string]*mirror.MetaValue{
			"target": {Value: &mirror.MetaValue_String_{String_: server2.URL}},
		},
		Headers: map[string]*mirror.HeaderValue{
			"X-My-Header":    {Values: []string{"value2_1", "value2_2"}},
			"X-Other-Header": {Values: []string{"value2_3"}},
		},
	}

	mod.SetInput(in)
	close(in)

	<-mod.Output()
	require.Equal(t, 1, reqCount1)
	require.Equal(t, 1, reqCount2)

}
