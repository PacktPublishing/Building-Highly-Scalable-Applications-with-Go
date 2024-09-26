package restv1

import "net/http"

// SpecHandler serves the OpenAPI spec from the generated code
func SpecHandler() http.HandlerFunc {
	// Use oapi-codegen's generated OpenAPI specification
	swagger, err := GetSwagger()
	if err != nil {
		panic("Failed to load embedded OpenAPI spec: " +
			err.Error())
	}

	// Marshal the spec into JSON format
	specBytes, err := swagger.MarshalJSON()
	if err != nil {
		panic("Failed to marshal OpenAPI spec to JSON: " +
			err.Error())
	}

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(specBytes) // Serve the spec as the response
	}
}
