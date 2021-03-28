package openapi

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
)

const specPath = "/spec.json"
const docPath = "/"

const docTemplate = `
<!DOCTYPE html>
<html>

<body>
    <redoc spec-url="./spec.json"></redoc>
    <script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"> </script>
</body>

</html>
`

// OpenAPI ...
type OpenAPI struct {
	spec *openapi3.Swagger
}

// NewFromFile creates openapi validation middleware
func NewFromFile(path string) (*OpenAPI, error) {
	spec, err := openapi3.NewSwaggerLoader().LoadSwaggerFromFile(path)
	if err != nil {
		return nil, err
	}
	return newFromSpec(spec)
}

// NewFromData creates openapi validation middleware
func NewFromData(data []byte) (*OpenAPI, error) {
	spec, err := openapi3.NewSwaggerLoader().LoadSwaggerFromData(data)
	if err != nil {
		return nil, err
	}
	return newFromSpec(spec)
}

func newFromSpec(spec *openapi3.Swagger) (*OpenAPI, error) {
	api := &OpenAPI{}
	api.spec = spec
	return api, nil
}

func (api *OpenAPI) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// serve static spec request
		if r.URL.Path == specPath && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(api.spec)
			return
		}

		// serve documentation request
		if r.URL.Path == docPath && r.Method == http.MethodGet {
			io.WriteString(w, docTemplate)
			return
		}

		next.ServeHTTP(w, r)
	})
}
