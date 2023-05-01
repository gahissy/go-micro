package micro

import (
	"github.com/gahissy/go-micro/h"
	"net/http"
)

type RouteGroup interface {
	GET(path string, cb func(ctx RequestContext) (interface{}, error))
	POST(path string, cb func(ctx RequestContext) (Any, error))
	PATCH(path string, cb func(ctx RequestContext) (Any, error))
	PUT(path string, cb func(ctx RequestContext) (Any, error))
	DELETE(path string, cb func(ctx RequestContext) (Any, error))
}

type Router interface {
	RouteGroup
	Group(path string, roles ...string) RouteGroup
	Start(port ...string)
	Handler() http.Handler
}

type RequestContext interface {
	Ctx() *Ctx
	Bind(data Any) error
	ShouldBind(data Any)
	CheckPermission(roles ...string)
}

func Handle0(cb func(ctx *Ctx) (interface{}, error)) func(RequestContext) (interface{}, error) {
	return Handle1("", cb)
}

func Handle1(role string, cb func(ctx *Ctx) (Any, error)) func(RequestContext) (Any, error) {
	return func(ctx RequestContext) (Any, error) {
		if !h.IsStrEmpty(role) {
			ctx.CheckPermission(role)
		}
		return cb(ctx.Ctx())
	}
}

func Handle2[K interface{}](input K, cb func(input K, ctx *Ctx) (Any, error)) func(RequestContext) (Any, error) {
	return Handle3("", input, cb)
}

func Handle3[K interface{}](role string, input K, cb func(input K, ctx *Ctx) (Any, error)) func(RequestContext) (Any, error) {
	return func(ctx RequestContext) (Any, error) {
		if !h.IsStrEmpty(role) {
			ctx.CheckPermission(role)
		}
		cloned := input
		ctx.ShouldBind(&cloned)
		return cb(cloned, ctx.Ctx())
	}
}
