package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServeHTTP(t *testing.T) {
	for name, test := range map[string]struct {
		RequestURL         string
		ExpectedStatusCode int
		ExpectedHeader     http.Header
		ExpectedbodyString string
	}{
		"Request status": {
			RequestURL:         "/status",
			ExpectedStatusCode: http.StatusOK,
			ExpectedHeader:     http.Header{"Content-Type": []string{"application/json"}},
			ExpectedbodyString: `{"enabled":true}`,
		},
		"Hello world": {
			RequestURL:         "/hello",
			ExpectedStatusCode: http.StatusOK,
			ExpectedHeader:     http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}},
			ExpectedbodyString: "Hello World!",
		},
		"InvalidRequestURL": {
			RequestURL:         "/not_found",
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedHeader:     http.Header{"Content-Type": []string{"text/plain; charset=utf-8"}, "X-Content-Type-Options": []string{"nosniff"}},
			ExpectedbodyString: "404 page not found\n",
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			plugin := &Plugin{}

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", test.RequestURL, nil)
			plugin.ServeHTTP(nil, w, r)

			result := w.Result()
			require.NotNil(t, result)

			bodyBytes, err := ioutil.ReadAll(result.Body)
			require.Nil(t, err)
			bodyString := string(bodyBytes)

			assert.Equal(test.ExpectedbodyString, bodyString)
			assert.Equal(test.ExpectedStatusCode, result.StatusCode)
			assert.Equal(test.ExpectedHeader, result.Header)
		})
	}
}
