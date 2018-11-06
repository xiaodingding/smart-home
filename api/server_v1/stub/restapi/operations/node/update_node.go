// Code generated by go-swagger; DO NOT EDIT.

package node

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// UpdateNodeHandlerFunc turns a function with the right signature into a update node handler
type UpdateNodeHandlerFunc func(UpdateNodeParams, interface{}) middleware.Responder

// Handle executing the request and returning a response
func (fn UpdateNodeHandlerFunc) Handle(params UpdateNodeParams, principal interface{}) middleware.Responder {
	return fn(params, principal)
}

// UpdateNodeHandler interface for that can handle valid update node params
type UpdateNodeHandler interface {
	Handle(UpdateNodeParams, interface{}) middleware.Responder
}

// NewUpdateNode creates a new http.Handler for the update node operation
func NewUpdateNode(ctx *middleware.Context, handler UpdateNodeHandler) *UpdateNode {
	return &UpdateNode{Context: ctx, Handler: handler}
}

/*UpdateNode swagger:route PUT /node/{id} node updateNode

update node

*/
type UpdateNode struct {
	Context *middleware.Context
	Handler UpdateNodeHandler
}

func (o *UpdateNode) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewUpdateNodeParams()

	uprinc, aCtx, err := o.Context.Authorize(r, route)
	if err != nil {
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}
	if aCtx != nil {
		r = aCtx
	}
	var principal interface{}
	if uprinc != nil {
		principal = uprinc
	}

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params, principal) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
