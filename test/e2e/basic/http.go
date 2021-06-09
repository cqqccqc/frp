package basic

import (
	"fmt"
	"net/http"

	"github.com/fatedier/frp/test/e2e/framework"
	"github.com/fatedier/frp/test/e2e/framework/consts"
	"github.com/fatedier/frp/test/e2e/mock/server/httpserver"
	"github.com/fatedier/frp/test/e2e/pkg/request"
	"github.com/fatedier/frp/tests/util"

	. "github.com/onsi/ginkgo"
)

var _ = Describe("[Feature: HTTP]", func() {
	f := framework.NewDefaultFramework()

	It("HTTP route by locations", func() {
		serverConf := consts.DefaultServerConfig
		vhostHTTPPort := f.AllocPort()
		serverConf += fmt.Sprintf(`
			vhost_http_port = %d
			`, vhostHTTPPort)

		fooPort := f.AllocPort()
		fooServer := httpserver.New(
			httpserver.WithBindPort(fooPort),
			httpserver.WithHandler(framework.SpecifiedHTTPBodyHandler([]byte("foo"))),
		)
		f.RunServer("", fooServer)

		barPort := f.AllocPort()
		barServer := httpserver.New(
			httpserver.WithBindPort(barPort),
			httpserver.WithHandler(framework.SpecifiedHTTPBodyHandler([]byte("bar"))),
		)
		f.RunServer("", barServer)

		clientConf := consts.DefaultClientConfig
		clientConf += fmt.Sprintf(`
			[foo]
			type = http
			local_port = %d
			custom_domains = normal.example.com
			locations = /,/foo

			[bar]
			type = http
			local_port = %d
			custom_domains = normal.example.com
			locations = /bar
			`, fooPort, barPort)

		f.RunProcesses([]string{serverConf}, []string{clientConf})

		// foo path
		framework.NewRequestExpect(f).Explain("foo path").Port(vhostHTTPPort).
			RequestModify(func(r *request.Request) {
				r.HTTP().HTTPHost("normal.example.com").HTTPPath("/foo")
			}).
			ExpectResp([]byte("foo")).
			Ensure()

		// bar path
		framework.NewRequestExpect(f).Explain("bar path").Port(vhostHTTPPort).
			RequestModify(func(r *request.Request) {
				r.HTTP().HTTPHost("normal.example.com").HTTPPath("/bar")
			}).
			ExpectResp([]byte("bar")).
			Ensure()

		// other path
		framework.NewRequestExpect(f).Explain("other path").Port(vhostHTTPPort).
			RequestModify(func(r *request.Request) {
				r.HTTP().HTTPHost("normal.example.com").HTTPPath("/other")
			}).
			ExpectResp([]byte("foo")).
			Ensure()
	})

	It("HTTP Basic Auth", func() {
		serverConf := consts.DefaultServerConfig
		vhostHTTPPort := f.AllocPort()
		serverConf += fmt.Sprintf(`
			vhost_http_port = %d
			`, vhostHTTPPort)

		clientConf := consts.DefaultClientConfig
		clientConf += fmt.Sprintf(`
			[test]
			type = http
			local_port = {{ .%s }}
			custom_domains = normal.example.com
			http_user = test
			http_pwd = test
			`, framework.HTTPSimpleServerPort)

		f.RunProcesses([]string{serverConf}, []string{clientConf})

		// not set auth header
		framework.NewRequestExpect(f).Port(vhostHTTPPort).
			RequestModify(func(r *request.Request) {
				r.HTTP().HTTPHost("normal.example.com")
			}).
			Ensure(func(resp *request.Response) bool {
				return resp.Code == 401
			})

		// set incorrect auth header
		framework.NewRequestExpect(f).Port(vhostHTTPPort).
			RequestModify(func(r *request.Request) {
				r.HTTP().HTTPHost("normal.example.com").HTTPHeaders(map[string]string{
					"Authorization": util.BasicAuth("test", "invalid"),
				})
			}).
			Ensure(func(resp *request.Response) bool {
				return resp.Code == 401
			})

		// set correct auth header
		framework.NewRequestExpect(f).Port(vhostHTTPPort).
			RequestModify(func(r *request.Request) {
				r.HTTP().HTTPHost("normal.example.com").HTTPHeaders(map[string]string{
					"Authorization": util.BasicAuth("test", "test"),
				})
			}).
			Ensure()
	})

	It("Wildcard domain", func() {
		// TODO

	})

	It("Subdomain", func() {
		// TODO
	})

	It("Modify headers", func() {
		serverConf := consts.DefaultServerConfig
		vhostHTTPPort := f.AllocPort()
		serverConf += fmt.Sprintf(`
			vhost_http_port = %d
			`, vhostHTTPPort)

		localPort := f.AllocPort()
		localServer := httpserver.New(
			httpserver.WithBindPort(localPort),
			httpserver.WithHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				w.Write([]byte(req.Header.Get("X-From-Where")))
			})),
		)
		f.RunServer("", localServer)

		clientConf := consts.DefaultClientConfig
		clientConf += fmt.Sprintf(`
			[test]
			type = http
			local_port = %d
			custom_domains = normal.example.com
			header_X-From-Where = frp
			`, localPort)

		f.RunProcesses([]string{serverConf}, []string{clientConf})

		// not set auth header
		framework.NewRequestExpect(f).Port(vhostHTTPPort).
			RequestModify(func(r *request.Request) {
				r.HTTP().HTTPHost("normal.example.com")
			}).
			ExpectResp([]byte("frp")). // local http server will write this X-From-Where header to response body
			Ensure()
	})
	It("Host Header Rewrite", func() {
		// TODO
	})
})
