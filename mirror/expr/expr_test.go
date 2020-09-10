package expr

import (
	"testing"

	"github.com/shimmerglass/http-mirror-pipeline/mirror"
	"github.com/stretchr/testify/require"
)

var testCases = map[string]interface{}{
	"abc":         "abc",
	"abc{1}":      "abc1",
	"{1}abc":      "1abc",
	"{1}abc{1}":   "1abc1",
	"a{1}abc{1}a": "a1abc1a",

	"{1}":     int64(1),
	"{true}":  true,
	"{false}": false,

	"{req.path}":           "/index.html",
	"{req.method == GET}":  true,
	`{req.header("Host")}`: "www.google.com",
	`{req.meta.key1.int}`:  int64(42),
}

func TestParseTmpl(t *testing.T) {
	r := mirror.Request{
		Method: mirror.Method_GET,
		Path:   "/index.html",
		Headers: map[string]*mirror.HeaderValue{
			"Host": {Values: []string{"www.google.com"}},
		},
		Meta: map[string]*mirror.MetaValue{
			"key1": {Value: &mirror.MetaValue_Int{Int: 42}},
		},
	}

	for str, expected := range testCases {
		t.Run(str, func(t *testing.T) {
			f, err := ParseTmpl(str)
			require.NoError(t, err)

			v, err := f.Eval(r)
			require.NoError(t, err)

			require.Equal(t, expected, v)
		})
	}
}
