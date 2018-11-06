// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	errors "github.com/go-openapi/errors"
	loads "github.com/go-openapi/loads"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	security "github.com/go-openapi/runtime/security"
	spec "github.com/go-openapi/spec"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/e154/smart-home/api/server_v1/stub/restapi/operations/index"
	"github.com/e154/smart-home/api/server_v1/stub/restapi/operations/node"
)

// NewServerAPI creates a new Server instance
func NewServerAPI(spec *loads.Document) *ServerAPI {
	return &ServerAPI{
		handlers:            make(map[string]map[string]http.Handler),
		formats:             strfmt.Default,
		defaultConsumes:     "application/json",
		defaultProduces:     "application/json",
		ServerShutdown:      func() {},
		spec:                spec,
		ServeError:          errors.ServeError,
		BasicAuthenticator:  security.BasicAuth,
		APIKeyAuthenticator: security.APIKeyAuth,
		BearerAuthenticator: security.BearerAuth,
		JSONConsumer:        runtime.JSONConsumer(),
		JSONProducer:        runtime.JSONProducer(),
		HTMLProducer: runtime.ProducerFunc(func(w io.Writer, data interface{}) error {
			return errors.NotImplemented("html producer has not yet been implemented")
		}),
		NodeAddNodeHandler: node.AddNodeHandlerFunc(func(params node.AddNodeParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation NodeAddNode has not yet been implemented")
		}),
		NodeDeleteNodeByIDHandler: node.DeleteNodeByIDHandlerFunc(func(params node.DeleteNodeByIDParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation NodeDeleteNodeByID has not yet been implemented")
		}),
		NodeGetNodeByIDHandler: node.GetNodeByIDHandlerFunc(func(params node.GetNodeByIDParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation NodeGetNodeByID has not yet been implemented")
		}),
		NodeGetNodeListHandler: node.GetNodeListHandlerFunc(func(params node.GetNodeListParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation NodeGetNodeList has not yet been implemented")
		}),
		IndexIndexHandler: index.IndexHandlerFunc(func(params index.IndexParams) middleware.Responder {
			return middleware.NotImplemented("operation IndexIndex has not yet been implemented")
		}),
		NodeUpdateNodeHandler: node.UpdateNodeHandlerFunc(func(params node.UpdateNodeParams, principal interface{}) middleware.Responder {
			return middleware.NotImplemented("operation NodeUpdateNode has not yet been implemented")
		}),

		// Applies when the "access_token" header is set
		APIKeyAuthAuth: func(token string) (interface{}, error) {
			return nil, errors.NotImplemented("api key auth (ApiKeyAuth) access_token from header param [access_token] has not yet been implemented")
		},

		// default authorizer is authorized meaning no requests are blocked
		APIAuthorizer: security.Authorized(),
	}
}

/*ServerAPI Server API */
type ServerAPI struct {
	spec            *loads.Document
	context         *middleware.Context
	handlers        map[string]map[string]http.Handler
	formats         strfmt.Registry
	defaultConsumes string
	defaultProduces string
	Middleware      func(middleware.Builder) http.Handler

	// BasicAuthenticator generates a runtime.Authenticator from the supplied basic auth function.
	// It has a default implemention in the security package, however you can replace it for your particular usage.
	BasicAuthenticator func(security.UserPassAuthentication) runtime.Authenticator
	// APIKeyAuthenticator generates a runtime.Authenticator from the supplied token auth function.
	// It has a default implemention in the security package, however you can replace it for your particular usage.
	APIKeyAuthenticator func(string, string, security.TokenAuthentication) runtime.Authenticator
	// BearerAuthenticator generates a runtime.Authenticator from the supplied bearer token auth function.
	// It has a default implemention in the security package, however you can replace it for your particular usage.
	BearerAuthenticator func(string, security.ScopedTokenAuthentication) runtime.Authenticator

	// JSONConsumer registers a consumer for a "application/json" mime type
	JSONConsumer runtime.Consumer

	// JSONProducer registers a producer for a "application/json" mime type
	JSONProducer runtime.Producer
	// HTMLProducer registers a producer for a "text/html" mime type
	HTMLProducer runtime.Producer

	// APIKeyAuthAuth registers a function that takes a token and returns a principal
	// it performs authentication based on an api key access_token provided in the header
	APIKeyAuthAuth func(string) (interface{}, error)

	// APIAuthorizer provides access control (ACL/RBAC/ABAC) by providing access to the request and authenticated principal
	APIAuthorizer runtime.Authorizer

	// NodeAddNodeHandler sets the operation handler for the add node operation
	NodeAddNodeHandler node.AddNodeHandler
	// NodeDeleteNodeByIDHandler sets the operation handler for the delete node by Id operation
	NodeDeleteNodeByIDHandler node.DeleteNodeByIDHandler
	// NodeGetNodeByIDHandler sets the operation handler for the get node by Id operation
	NodeGetNodeByIDHandler node.GetNodeByIDHandler
	// NodeGetNodeListHandler sets the operation handler for the get node list operation
	NodeGetNodeListHandler node.GetNodeListHandler
	// IndexIndexHandler sets the operation handler for the index operation
	IndexIndexHandler index.IndexHandler
	// NodeUpdateNodeHandler sets the operation handler for the update node operation
	NodeUpdateNodeHandler node.UpdateNodeHandler

	// ServeError is called when an error is received, there is a default handler
	// but you can set your own with this
	ServeError func(http.ResponseWriter, *http.Request, error)

	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	ServerShutdown func()

	// Custom command line argument groups with their descriptions
	CommandLineOptionsGroups []swag.CommandLineOptionsGroup

	// User defined logger function.
	Logger func(string, ...interface{})
}

// SetDefaultProduces sets the default produces media type
func (o *ServerAPI) SetDefaultProduces(mediaType string) {
	o.defaultProduces = mediaType
}

// SetDefaultConsumes returns the default consumes media type
func (o *ServerAPI) SetDefaultConsumes(mediaType string) {
	o.defaultConsumes = mediaType
}

// SetSpec sets a spec that will be served for the clients.
func (o *ServerAPI) SetSpec(spec *loads.Document) {
	o.spec = spec
}

// DefaultProduces returns the default produces media type
func (o *ServerAPI) DefaultProduces() string {
	return o.defaultProduces
}

// DefaultConsumes returns the default consumes media type
func (o *ServerAPI) DefaultConsumes() string {
	return o.defaultConsumes
}

// Formats returns the registered string formats
func (o *ServerAPI) Formats() strfmt.Registry {
	return o.formats
}

// RegisterFormat registers a custom format validator
func (o *ServerAPI) RegisterFormat(name string, format strfmt.Format, validator strfmt.Validator) {
	o.formats.Add(name, format, validator)
}

// Validate validates the registrations in the ServerAPI
func (o *ServerAPI) Validate() error {
	var unregistered []string

	if o.JSONConsumer == nil {
		unregistered = append(unregistered, "JSONConsumer")
	}

	if o.JSONProducer == nil {
		unregistered = append(unregistered, "JSONProducer")
	}

	if o.HTMLProducer == nil {
		unregistered = append(unregistered, "HTMLProducer")
	}

	if o.APIKeyAuthAuth == nil {
		unregistered = append(unregistered, "AccessTokenAuth")
	}

	if o.NodeAddNodeHandler == nil {
		unregistered = append(unregistered, "node.AddNodeHandler")
	}

	if o.NodeDeleteNodeByIDHandler == nil {
		unregistered = append(unregistered, "node.DeleteNodeByIDHandler")
	}

	if o.NodeGetNodeByIDHandler == nil {
		unregistered = append(unregistered, "node.GetNodeByIDHandler")
	}

	if o.NodeGetNodeListHandler == nil {
		unregistered = append(unregistered, "node.GetNodeListHandler")
	}

	if o.IndexIndexHandler == nil {
		unregistered = append(unregistered, "index.IndexHandler")
	}

	if o.NodeUpdateNodeHandler == nil {
		unregistered = append(unregistered, "node.UpdateNodeHandler")
	}

	if len(unregistered) > 0 {
		return fmt.Errorf("missing registration: %s", strings.Join(unregistered, ", "))
	}

	return nil
}

// ServeErrorFor gets a error handler for a given operation id
func (o *ServerAPI) ServeErrorFor(operationID string) func(http.ResponseWriter, *http.Request, error) {
	return o.ServeError
}

// AuthenticatorsFor gets the authenticators for the specified security schemes
func (o *ServerAPI) AuthenticatorsFor(schemes map[string]spec.SecurityScheme) map[string]runtime.Authenticator {

	result := make(map[string]runtime.Authenticator)
	for name, scheme := range schemes {
		switch name {

		case "ApiKeyAuth":

			result[name] = o.APIKeyAuthenticator(scheme.Name, scheme.In, o.APIKeyAuthAuth)

		}
	}
	return result

}

// Authorizer returns the registered authorizer
func (o *ServerAPI) Authorizer() runtime.Authorizer {

	return o.APIAuthorizer

}

// ConsumersFor gets the consumers for the specified media types
func (o *ServerAPI) ConsumersFor(mediaTypes []string) map[string]runtime.Consumer {

	result := make(map[string]runtime.Consumer)
	for _, mt := range mediaTypes {
		switch mt {

		case "application/json":
			result["application/json"] = o.JSONConsumer

		}
	}
	return result

}

// ProducersFor gets the producers for the specified media types
func (o *ServerAPI) ProducersFor(mediaTypes []string) map[string]runtime.Producer {

	result := make(map[string]runtime.Producer)
	for _, mt := range mediaTypes {
		switch mt {

		case "application/json":
			result["application/json"] = o.JSONProducer

		case "text/html":
			result["text/html"] = o.HTMLProducer

		}
	}
	return result

}

// HandlerFor gets a http.Handler for the provided operation method and path
func (o *ServerAPI) HandlerFor(method, path string) (http.Handler, bool) {
	if o.handlers == nil {
		return nil, false
	}
	um := strings.ToUpper(method)
	if _, ok := o.handlers[um]; !ok {
		return nil, false
	}
	if path == "/" {
		path = ""
	}
	h, ok := o.handlers[um][path]
	return h, ok
}

// Context returns the middleware context for the server API
func (o *ServerAPI) Context() *middleware.Context {
	if o.context == nil {
		o.context = middleware.NewRoutableContext(o.spec, o, nil)
	}

	return o.context
}

func (o *ServerAPI) initHandlerCache() {
	o.Context() // don't care about the result, just that the initialization happened

	if o.handlers == nil {
		o.handlers = make(map[string]map[string]http.Handler)
	}

	if o.handlers["POST"] == nil {
		o.handlers["POST"] = make(map[string]http.Handler)
	}
	o.handlers["POST"]["/node"] = node.NewAddNode(o.context, o.NodeAddNodeHandler)

	if o.handlers["DELETE"] == nil {
		o.handlers["DELETE"] = make(map[string]http.Handler)
	}
	o.handlers["DELETE"]["/node/{id}"] = node.NewDeleteNodeByID(o.context, o.NodeDeleteNodeByIDHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/node/{id}"] = node.NewGetNodeByID(o.context, o.NodeGetNodeByIDHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"]["/node"] = node.NewGetNodeList(o.context, o.NodeGetNodeListHandler)

	if o.handlers["GET"] == nil {
		o.handlers["GET"] = make(map[string]http.Handler)
	}
	o.handlers["GET"][""] = index.NewIndex(o.context, o.IndexIndexHandler)

	if o.handlers["PUT"] == nil {
		o.handlers["PUT"] = make(map[string]http.Handler)
	}
	o.handlers["PUT"]["/node/{id}"] = node.NewUpdateNode(o.context, o.NodeUpdateNodeHandler)

}

// Serve creates a http handler to serve the API over HTTP
// can be used directly in http.ListenAndServe(":8000", api.Serve(nil))
func (o *ServerAPI) Serve(builder middleware.Builder) http.Handler {
	o.Init()

	if o.Middleware != nil {
		return o.Middleware(builder)
	}
	return o.context.APIHandler(builder)
}

// Init allows you to just initialize the handler cache, you can then recompose the middelware as you see fit
func (o *ServerAPI) Init() {
	if len(o.handlers) == 0 {
		o.initHandlerCache()
	}
}
