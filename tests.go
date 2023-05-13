package micro

import (
	"github.com/gavv/httpexpect/v2"
	"net/http/httptest"
	"testing"
)

func (a *App) UseTest(t *testing.T, handler func(e *TestSupport)) {
	srv := httptest.NewServer(a.Handler())
	defer srv.Close()
	e := httpexpect.Default(t, srv.URL)
	handler(&TestSupport{
		t: t,
		Http: HttpTestSupport{
			t:        t,
			internal: e,
		},
	})
}
