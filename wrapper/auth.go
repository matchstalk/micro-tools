package wrapper

import (
	"github.com/micro/go-micro/server"
	"os"
	"strings"
	"context"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-log"
	"net/http"
	"github.com/micro/go-micro/errors"
	"github.com/matchstalk/micro-tools/library"
)

var (
	AuthPrefix = "X-Auth-"
	Secret = "auth"
)

func WrapperAuth(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		if os.Getenv("DISABLE_AUTH") == "true" {
			return fn(ctx, req, resp)
		} else {
			md, _ := metadata.FromContext(ctx)
			if strings.ContainsAny(md["Authorization"], "bearer") {
				if claims, err := library.VerifyJwt(Secret, md["Authorization"][7:]); err != nil {
					log.Log(err)
					return errors.New(req.Service() + "." + req.Method(), err.Error(), http.StatusUnauthorized)
				} else {
					for _, key := range []string{"User", "Role", "Telnet"} {
						if u, err := claims.Get(key); err == nil {
							md[AuthPrefix + key] = u.(string)
						}
					}
					ctx = metadata.NewContext(ctx, md)
				}
			}
			return fn(ctx, req, resp)
		}

	}
}

func Verfiy(ctx context.Context) (m map[string]string, ok bool) {
	if md, ok := metadata.FromContext(ctx); ok {
		ok = true
		for _, key := range []string{"User", "Role", "Telnet"} {
			if v, found := md[AuthPrefix + key]; found {
				m[key] = v
			}
		}
	}
	return
}