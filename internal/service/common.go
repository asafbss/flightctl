package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/flightctl/flightctl/api/v1alpha1"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers/gorillamux"
)

const (
	Forbidden                      = "Forbidden"
	AuthorizationServerUnavailable = "Authorization server unavailable"
)

func validateAgainstSchema(ctx context.Context, obj []byte, objPath string) error {
	swagger, err := v1alpha1.GetSwagger()
	if err != nil {
		return err
	}
	// Skip server name validation
	swagger.Servers = nil

	url, err := url.Parse(objPath)
	if err != nil {
		return err
	}
	httpReq := &http.Request{
		Method: "PUT",
		URL:    url,
		Body:   io.NopCloser(bytes.NewReader(obj)),
		Header: http.Header{"Content-Type": []string{"application/json"}},
	}

	router, err := gorillamux.NewRouter(swagger)
	if err != nil {
		return err
	}
	route, pathParams, err := router.FindRoute(httpReq)
	if err != nil {
		return err
	}

	requestValidationInput := &openapi3filter.RequestValidationInput{
		Request:    httpReq,
		PathParams: pathParams,
		Route:      route,
	}
	return openapi3filter.ValidateRequest(ctx, requestValidationInput)
}

func ApplyJSONPatch[T any](ctx context.Context, obj T, newObj T, patchRequest v1alpha1.PatchRequest, objPath string) error {
	patch, err := json.Marshal(patchRequest)
	if err != nil {
		return err
	}
	jsonPatch, err := jsonpatch.DecodePatch(patch)
	if err != nil {
		return err
	}

	objJSON, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	newJSON, err := jsonPatch.Apply(objJSON)
	if err != nil {
		return err
	}

	//validate the new object against OpenAPI schema
	err = validateAgainstSchema(ctx, newJSON, objPath)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewReader(newJSON))
	decoder.DisallowUnknownFields()
	return decoder.Decode(&newObj)
}
