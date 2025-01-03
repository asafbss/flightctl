// Package v1alpha1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package v1alpha1

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"q/8OAAD//yG2pElomQEA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
