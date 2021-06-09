package framework

import (
	"net/http"

	"github.com/fatedier/frp/test/e2e/framework/consts"
	"github.com/fatedier/frp/test/e2e/pkg/request"
)

func SpecifiedHTTPBodyHandler(body []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Write(body)
	}
}

// NewRequest return a default request with default timeout and content.
func NewRequest() *request.Request {
	return request.New().
		Timeout(consts.DefaultTimeout).
		Body([]byte(consts.TestString))
}

func NewHTTPRequest() *request.Request {
	return request.New().HTTP().HTTPParams("GET", "", "/", nil)
}

func ExpectResponse(req *request.Request, expectResp []byte, explain ...interface{}) {
	ret, err := req.Do()
	ExpectNoError(err, explain...)
	ExpectEqualValues(expectResp, ret.Content, explain...)
}

func ExpectResponseError(req *request.Request, explain ...interface{}) {
	_, err := req.Do()
	ExpectError(err, explain...)
}

type RequestExpect struct {
	req *request.Request

	f           *Framework
	expectResp  []byte
	expectError bool
	explain     []interface{}
}

func NewRequestExpect(f *Framework) *RequestExpect {
	return &RequestExpect{
		req:         NewRequest(),
		f:           f,
		expectResp:  []byte(consts.TestString),
		expectError: false,
		explain:     make([]interface{}, 0),
	}
}

func (e *RequestExpect) RequestModify(f func(r *request.Request)) *RequestExpect {
	f(e.req)
	return e
}

func (e *RequestExpect) Protocol(protocol string) *RequestExpect {
	e.req.Protocol(protocol)
	return e
}

func (e *RequestExpect) PortName(name string) *RequestExpect {
	if e.f != nil {
		e.req.Port(e.f.PortByName(name))
	}
	return e
}

func (e *RequestExpect) Port(port int) *RequestExpect {
	if e.f != nil {
		e.req.Port(port)
	}
	return e
}

func (e *RequestExpect) ExpectResp(resp []byte) *RequestExpect {
	e.expectResp = resp
	return e
}

func (e *RequestExpect) ExpectError(expectErr bool) *RequestExpect {
	e.expectError = expectErr
	return e
}

func (e *RequestExpect) Explain(explain ...interface{}) *RequestExpect {
	e.explain = explain
	return e
}

type EnsureFunc func(*request.Response) bool

func (e *RequestExpect) Ensure(fns ...EnsureFunc) {
	ret, err := e.req.Do()
	if e.expectError {
		ExpectError(err, e.explain...)
		return
	}
	ExpectNoError(err, e.explain...)

	if len(fns) == 0 {
		ExpectEqualValues(e.expectResp, ret.Content, e.explain...)
	} else {
		for _, fn := range fns {
			ok := fn(ret)
			ExpectTrue(ok, e.explain...)
		}
	}
}
