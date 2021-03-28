package openapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const spec = `

`

func TestOpenAPIMiddleware(t *testing.T) {

	api, err := NewFromData([]byte(spec))
	assert.NoError(t, err)

	router := mux.NewRouter().StrictSlash(true)
	ts := httptest.NewServer(api.Middleware(router))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	res, err = http.Get(ts.URL + "/spec.json")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

}
