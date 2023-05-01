package micro

import (
	"fmt"
	"gahissy/studio/app/core"
	"github.com/gahissy/go-micro/h"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
	"os"
	"strings"
)

type routerImpl struct {
	Router
	engine *echo.Echo
}

type routerRequestContextImpl struct {
	RequestContext
	wrapper echo.Context
	ctx     *Ctx
}

type routeGroupImpl struct {
	wrapped *echo.Group
}

const EnvKey = "env"
const AuthKey = "auth"

func NewEchoRouter(env *Env) Router {
	e := echo.New()
	if os.Getenv("ENV") != "production" {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}\n",
		}))
	}
	//e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(sentryecho.New(sentryecho.Options{}))
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set(EnvKey, env)
			c.Set(AuthKey, &Auth{
				Id:            "guest",
				Role:          "guest",
				Authenticated: false,
			})
			return next(c)
		}
	})
	e.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(os.Getenv("ENCRYPTION_KEY")),
		ContextKey: "user",
		Skipper: func(c echo.Context) bool {
			//TODO: make this configurable
			return c.Request().Method == "OPTIONS" ||
				strings.HasPrefix(c.Path(), "/swagger/") ||
				c.Path() == "/status" || c.Path() == "/" || c.Path() == "/auth" ||
				strings.HasPrefix(c.Path(), "/pub/")
		},
		SuccessHandler: func(c echo.Context) {
			user := c.Get("user")
			if user == nil {
				return
			}
			claims := user.(*jwt.Token).Claims.(jwt.MapClaims)
			c.Set(AuthKey, &Auth{
				Id:            claims["sub"].(string),
				Role:          claims["role"].(string),
				Authenticated: true,
			})
		},
	}))
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	return &routerImpl{engine: e}
}

func (r *routerImpl) Handler() http.Handler {
	return r.engine
}

func (r *routerImpl) Start(port ...string) {
	r.engine.GET("/swagger/*", echoSwagger.WrapHandler)
	lport := ""
	if port != nil && len(port) > 0 {
		lport = port[0]
	}
	if h.IsStrEmpty(lport) {
		lport = os.Getenv("PORT")
	}
	if h.IsStrEmpty(lport) {
		lport = "8080"
	}
	r.engine.Logger.Fatal(r.engine.Start(fmt.Sprintf(":%s", lport)))
}

type Context struct {
	internal echo.Context
}

// ---------------------------------------------------------------------------------------------------

func (r *routeGroupImpl) GET(path string, cb func(ctx RequestContext) (Any, error)) {
	r.handle("GET", path, cb)
}

func (r *routeGroupImpl) POST(path string, cb func(ctx RequestContext) (Any, error)) {
	r.handle("POST", path, cb)
}

func (r *routeGroupImpl) PUT(path string, cb func(ctx RequestContext) (Any, error)) {
	r.handle("PUT", path, cb)
}

func (r *routeGroupImpl) PATCH(path string, cb func(ctx RequestContext) (Any, error)) {
	r.handle("PATCH", path, cb)
}

func (r *routeGroupImpl) DELETE(path string, cb func(ctx RequestContext) (Any, error)) {
	r.handle("DELETE", path, cb)
}

func (r *routeGroupImpl) handle(method string, path string, cb func(ctx RequestContext) (Any, error)) {
	r.wrapped.Match([]string{method}, h.NormalizeUri(path), func(c echo.Context) error {
		auth := c.Get(AuthKey).(*Auth)
		env := c.Get(EnvKey).(*Env)
		requestContext := &routerRequestContextImpl{
			wrapper: c,
			ctx: &Ctx{
				Env:  env,
				Auth: auth,
			},
		}
		res, err := cb(requestContext)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, res)
	})
}

// ---------------------------------------------------------------------------------------------------

func (r *routerImpl) GET(path string, cb func(ctx RequestContext) (Any, error)) {
	r.handle("GET", path, cb)
}

func (r *routerImpl) POST(path string, cb func(ctx RequestContext) (Any, error)) {
	r.handle("POST", path, cb)
}

func (r *routerImpl) PUT(path string, cb func(ctx RequestContext) (Any, error)) {
	r.handle("PUT", path, cb)
}

func (r *routerImpl) PATCH(path string, cb func(ctx RequestContext) (Any, error)) {
	r.handle("PATCH", path, cb)
}

func (r *routerImpl) DELETE(path string, cb func(ctx RequestContext) (Any, error)) {
	r.handle("DELETE", path, cb)
}

func (r *routerImpl) handle(method string, path string, cb func(ctx RequestContext) (Any, error)) {
	r.engine.Match(
		[]string{method}, h.NormalizeUri(path), func(c echo.Context) error {
			auth := c.Get(AuthKey).(*Auth)
			env := c.Get(EnvKey).(*Env)
			rctw := &routerRequestContextImpl{
				wrapper: c,
				ctx: &Ctx{
					Env:  env,
					Auth: auth,
				},
			}
			res, err := cb(rctw)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, res)
		})
}

func (r *routerImpl) Group(path string, permissions ...string) RouteGroup {
	group := r.engine.Group(h.NormalizeUri(path))

	if permissions != nil && len(permissions) > 0 && !h.Contains(permissions, core.Authenticated) {
		group.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				auth := c.Get(AuthKey).(*Auth)
				checkPermission(auth, permissions...)
				return next(c)
			}
		})
	}

	return &routeGroupImpl{wrapped: group}
}

func (r *routerRequestContextImpl) Ctx() *Ctx {
	return r.ctx
}

func (r *routerRequestContextImpl) CheckPermission(permissions ...string) {
	checkPermission(r.ctx.Auth, permissions...)
}

func (r *routerRequestContextImpl) Bind(input interface{}) error {
	if err := r.wrapper.Bind(input); err != nil {
		return err
	}
	if err := validator.New().Struct(input); err != nil {
		//validationErrors := err.(validator.ValidationErrors)
		return err
	}
	return nil
}

func (r *routerRequestContextImpl) ShouldBind(input interface{}) {
	if err := r.Bind(input); err != nil {
		panic(echo.NewHTTPError(http.StatusBadRequest, "validation.failed"))
	}
}

func checkPermission(auth *Auth, permissions ...string) {
	if len(permissions) == 0 {
		return
	}
	if !auth.Authenticated {
		panic(echo.NewHTTPError(http.StatusForbidden, "permission.denied"))
	}
	if auth.Role == core.Admin {
		return
	}
	if h.Contains(permissions, core.Authenticated) {
		return
	}
	if h.IsStrEmpty(auth.Role) || !h.Contains(permissions, auth.Role) {
		panic(echo.NewHTTPError(http.StatusForbidden, "permission.denied"))
	}
}
