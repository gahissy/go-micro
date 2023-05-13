package micro

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gahissy/go-micro/h"
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type TestSupport struct {
	Http HttpTestSupport
	t    *testing.T
}

type HttpReqOpts struct {
	Bearer      string
	QueryParams map[string]interface{}
	Data        interface{}
	Headers     map[string]string
}

type HttpTestSupport struct {
	internal *httpexpect.Expect
	t        *testing.T
	basePath string
	config   *HttpReqOpts
}

type HttpTestExpectations struct {
	internal *httpexpect.Response
}

func (ts *HttpTestSupport) Group(path string, config ...HttpReqOpts) *HttpTestSupport {
	var baseConfig *HttpReqOpts
	if config != nil {
		if len(config) > 1 {
			ts.t.Fatalf("more than 1 HttpReqOpts provided")
			return nil
		}
		baseConfig = &config[0]
	}
	return &HttpTestSupport{
		internal: ts.internal,
		t:        ts.t,
		basePath: path,
		config:   baseConfig,
	}
}

func (ts *HttpTestSupport) GET(path string, config ...HttpReqOpts) *HttpTestExpectations {
	return ts.Request(http.MethodGet, path, config...)
}

func (ts *HttpTestSupport) POST(path string, config ...HttpReqOpts) *HttpTestExpectations {
	return ts.Request(http.MethodPost, path, config...)
}

func (ts *HttpTestSupport) PATCH(path string, config ...HttpReqOpts) *HttpTestExpectations {
	return ts.Request(http.MethodPatch, path, config...)
}

func (ts *HttpTestSupport) DELETE(path string, config ...HttpReqOpts) *HttpTestExpectations {
	return ts.Request(http.MethodDelete, path, config...)
}

func (ts *HttpTestSupport) Request(method string, path string, config ...HttpReqOpts) *HttpTestExpectations {
	fullPath := h.JoinUrl(ts.basePath, path)
	req := ts.internal.Request(method, fullPath)

	if ts.config != nil {
		if ts.config.Bearer != "" {
			req.WithHeader("Authorization", "Bearer "+ts.config.Bearer)
		}
	}

	if config != nil {
		if len(config) > 1 {
			ts.t.Fatalf("more than 1 HttpReqOpts provided")
			return nil
		}
		c := config[0]
		if c.Bearer != "" {
			req.WithHeader("Authorization", "Bearer "+c.Bearer)
		}
		if c.Data != nil {
			req.WithJSON(c.Data)
		}
	}
	return &HttpTestExpectations{
		internal: req.Expect(),
	}
}

func (te *HttpTestExpectations) IsOk() *HttpTestExpectations {
	te.internal.Status(http.StatusOK)
	return te
}

func (te *HttpTestExpectations) IsBadRequest() *HttpTestExpectations {
	te.internal.Status(http.StatusBadRequest)
	return te
}

func (te *HttpTestExpectations) IsUnauthorized() *HttpTestExpectations {
	te.internal.Status(http.StatusUnauthorized)
	return te
}

func (te *HttpTestExpectations) IsForbidden() *HttpTestExpectations {
	te.internal.Status(http.StatusForbidden)
	return te
}

func (te *HttpTestExpectations) IsNotFound() *HttpTestExpectations {
	te.internal.Status(http.StatusNotFound)
	return te
}

func (te *HttpTestExpectations) Is(status int) *HttpTestExpectations {
	te.internal.Status(status)
	return te
}

func (te *HttpTestExpectations) JSON() *httpexpect.Value {
	return te.internal.JSON()
}

type FakerFactory = func(field string) interface{}

func FakeIt(obj interface{}, factory ...FakerFactory) {
	// Get the value of the object as a reflect.Value.
	val := reflect.ValueOf(obj)

	// If the object is a pointer, get the underlying value.
	if val.Kind() != reflect.Ptr {
		panic("obj must be a pointer")
	}

	// If the object is a pointer, get the underlying value.
	val = val.Elem()

	// Create a new instance of the object's type.

	// Loop over the object's fields.
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldName := val.Type().Field(i).Name
		fieldType := val.Type().Field(i).Type

		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		lFieldName := strings.ToLower(fieldName)
		// If the field is unexported, skip it.
		if !fieldVal.CanSet() {
			continue
		}

		var fakeData interface{} = nil

		if factory != nil {
			fakeData = factory[0](fieldName)
		}

		// If the field is another struct, recursively call FakeObject on it.
		if fieldType.Kind() == reflect.Struct {
			newFieldVal := reflect.New(fieldType)
			if newFieldVal.Elem().CanSet() {
				FakeIt(newFieldVal.Interface())
				fieldVal.Set(newFieldVal)
			}
			continue
		}

		if fakeData == nil {
			// Generate fake data based on the field's type.
			switch fieldType.Kind() {
			case reflect.String:
				if strings.HasSuffix(lFieldName, "id") {
					continue
				} else {
					fakeData = FakeString(lFieldName)
				}
			case reflect.Int:
				fakeData = gofakeit.IntRange(1, 10000)
			case reflect.Int64:
				fakeData = gofakeit.Int64()
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fakeData = gofakeit.Uint64()
			case reflect.Float32:
				fakeData = gofakeit.Float32Range(0, 10000)
			case reflect.Float64:
				fakeData = gofakeit.Float64Range(0, 10000)
			case reflect.Bool:
				fakeData = gofakeit.Bool()
			default:
				// If we don't know how to generate fake data for this type,
				// skip the field.
				continue
			}
		}

		// Set the fake data on the new object.
		if fieldVal.Kind() == reflect.Ptr {
			newFieldVal := reflect.New(fieldType)
			if newFieldVal.Elem().CanSet() {
				newFieldVal.Elem().Set(reflect.ValueOf(fakeData))
				fieldVal.Set(newFieldVal)
			}
		} else {
			// Otherwise, set the fake data directly on the field.
			newFieldVal := val.FieldByName(fieldName)
			if newFieldVal.CanSet() {
				newFieldVal.Set(reflect.ValueOf(fakeData))
			}
		}
	}

	// Return the new object as an interface{}.
}

func FakeString(fieldName string) string {
	if strings.Contains(fieldName, "email") {
		return gofakeit.Email()
	} else if strings.Contains(fieldName, "phone") {
		return gofakeit.Phone()
	} else if strings.Contains(fieldName, "url") {
		return gofakeit.URL()
	} else if strings.Contains(fieldName, "address") {
		return gofakeit.Address().Address
	} else if strings.Contains(fieldName, "image") {
		return gofakeit.ImageURL(200, 200)
	} else if strings.Contains(fieldName, "brand") {
		return gofakeit.Company()
	} else if strings.Contains(fieldName, "description") {
		return gofakeit.SentenceSimple()
	} else if strings.Contains(fieldName, "currency") {
		return gofakeit.CurrencyShort()
	} else if strings.Contains(fieldName, "color") {
		return gofakeit.HexColor()
	} else {
		return gofakeit.Word()
	}
}
