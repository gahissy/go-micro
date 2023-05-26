package micro

import (
	"errors"
	"fmt"
	"github.com/gahissy/go-micro/h"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

type routerImpl struct {
	Router
	engine *echo.Echo
	//publicRoutes []string
}

type routerRequestContextImpl struct {
	RequestContext
	wrapper echo.Context
	ctx     *Ctx
}

type routeGroupImpl struct {
	wrapped *echo.Group
}

const EnvKey = "Env"
const AuthKey = "auth"

func NewEchoRouter(env *Env, config RouterConfig) Router {

	e := echo.New()

	impl := &routerImpl{
		engine: e,
	}
	if !env.IsProduction() {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}\n",
		}))
	}
	//e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(sentryecho.New(sentryecho.Options{}))
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	if config.Cors != nil && config.Cors.Enabled {
		e.Use(middleware.CORS())
	}
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
	if config.JwtConfig != nil {
		e.Use(echojwt.WithConfig(echojwt.Config{
			SigningKey: []byte(config.JwtConfig.Secret),
			ContextKey: "user",
			Skipper: func(c echo.Context) bool {
				path := h.NormalizeUri(c.Request().RequestURI)
				//TODO: make this configurable
				if c.Request().Method == "OPTIONS" {
					return true // preflight Request
				}
				if strings.HasSuffix(path, ".ico") {
					return true
				}
				if strings.HasPrefix(path, "/swagger/") {
					return true // swagger resources
				}
				if path == "/status" || path == "/" {
					return true // healthcheck and api info
				}
				if strings.HasPrefix(path, "/pub/") {
					return true // public resources
				}
				if path == "/auth" {
					return true // auth Request
				}
				for _, r := range config.PublicRoutes {
					if strings.EqualFold(path, h.NormalizeUri(r)) {
						return true
					}
				}

				return false
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
			ErrorHandler: func(c echo.Context, err error) error {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
			},
		}))
	}
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	for _, route := range config.Routes {
		route(impl)
	}

	return impl
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
		lport = h.GetEnv("PORT")
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
		return handleResponse(c, res, err)
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
			return handleResponse(c, res, err)
		})
}

func (r *routerImpl) Group(path string, permissions ...string) RouteGroup {
	group := r.engine.Group(h.NormalizeUri(path))

	if permissions != nil && len(permissions) > 0 && !h.Contains(permissions, RoleAuthenticated) {
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
		return echo.NewHTTPError(
			http.StatusBadRequest,
			h.Map{
				"kind":    "input.binding",
				"message": "input.binding.failed",
				"error":   err.Error(),
			},
		)
	}
	if err := validator.New().Struct(input); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errors []string
		for _, e := range validationErrors {
			errors = append(errors, fmt.Sprintf("%s: %s", e.Field(), e.Tag()))
		}
		return echo.NewHTTPError(
			http.StatusBadRequest,
			h.Map{
				"kind":    "validation",
				"message": "validation.failed",
				"errors":  errors,
			},
		)

	}
	return nil
}

func (r *routerRequestContextImpl) ShouldBind(input interface{}) {
	if err := r.Bind(input); err != nil {
		panic(err)
	}
}

func checkPermission(auth *Auth, permissions ...string) {
	if len(permissions) == 0 {
		return
	}
	if len(permissions) == 1 && permissions[0] == "*" {
		return
	}
	if !auth.Authenticated {
		if h.Contains(permissions, RoleAnonymous) {
			return
		}
		panic(echo.NewHTTPError(http.StatusForbidden, "permission.denied"))
	}
	if auth.Role == RoleAdmin {
		return
	}
	if h.Contains(permissions, RoleAuthenticated) {
		return
	}
	if h.IsStrEmpty(auth.Role) || !h.Contains(permissions, auth.Role) {
		panic(echo.NewHTTPError(http.StatusForbidden, "permission.denied"))
	}
}

func handleResponse(c echo.Context, res interface{}, err error) error {
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "resource.not_found")
		}

		if ferr, ok := err.(*h.FunctionalError); ok {
			// This is a FunctionalError, you can access the error code and message
			fmt.Printf("Functional error with code %d: %s\n", ferr.Code, ferr.Message)
			return c.JSON(http.StatusBadRequest, ferr)
		}

		if ferr, ok := err.(*h.ForbiddenError); ok {
			// This is a FunctionalError, you can access the error code and message
			return c.JSON(http.StatusBadRequest, ferr)
		}
		if ferr, ok := err.(*h.TechnicalError); ok {
			// This is a FunctionalError, you can access the error code and message
			return c.JSON(http.StatusInternalServerError, ferr)
		}
		if ferr, ok := err.(*h.UnauthorizedError); ok {
			// This is a FunctionalError, you can access the error code and message
			return c.JSON(http.StatusUnauthorized, ferr)
		}
		if ferr, ok := err.(*h.ResourceNotFoundError); ok {
			// This is a FunctionalError, you can access the error code and message
			return c.JSON(http.StatusNotFound, ferr)
		}
		if eerr, ok := err.(*echo.HTTPError); ok {
			return c.JSON(eerr.Code, eerr.Message)
		}
		return c.JSON(500, h.Map{
			"error": err.Error(),
			"kind":  "technical",
			"time":  time.Now().UTC(),
		})
	}
	return c.JSON(http.StatusOK, res)

}
