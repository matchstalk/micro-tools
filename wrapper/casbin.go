package wrapper

import (
	"net/http"
	"context"
	"github.com/casbin/casbin"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-micro/metadata"
	"github.com/cicdi-go/sso/src/utils"
	"github.com/micro/go-micro/errors"
)

var (
	CasbinAdapter *Casbin
)

// Wrapper is the router wrapper, prefer this method if you want to use casbin to your entire iris application.
// Usage:
// [...]
// app.WrapRouter(casbinMiddleware.Wrapper())
// app.Get("/dataset1/resource1", myHandler)
// [...]
func WrapperCasbin(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		if !CasbinAdapter.Check(ctx, req) {
			return errors.New(req.Service() + "." + req.Method(), "当前用户没有访问权限!", http.StatusForbidden)
		}
		return fn(ctx, req, resp)
	}
}


// Casbin is the auth services which contains the casbin enforcer.
type Casbin struct {
	enforcer *casbin.Enforcer
}

// Check checks the username, request's method and path and
// returns true if permission grandted otherwise false.
func (c *Casbin) Check(ctx context.Context, r server.Request) bool {
	username := Username(ctx)
	method := r.Method()
	service := r.Service()
	return c.enforcer.Enforce(username, service, method)
}

// Username gets the username from the basicauth.
func Username(ctx context.Context) string {
	md, _ := metadata.FromContext(ctx)
	username := md[utils.Config.AuthPrefix + "Username"]
	//role := md[utils.Config.AuthPrefix + "Role"]
	//telnet := md[utils.Config.AuthPrefix + "Telnet"]
	return username
}