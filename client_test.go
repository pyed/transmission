package transmission

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/auth"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	cMux *http.ServeMux

	client *ApiClient

	cServer *httptest.Server
)

func RPCHandler(res http.ResponseWriter, req *http.Request) {
	if req.Header.Get("X-Transmission-Session-Id") == "" {
		res.Header().Set("X-Transmission-Session-Id", "123")
		res.WriteHeader(http.StatusConflict)
		return
	}
	fmt.Fprintf(res, `{"arguments":{},"result":"no method name"}`)
}

func cSetup() {
	// test server
	cMux = http.NewServeMux()
	m := martini.New()
	r := martini.NewRouter()
	r.Post("/transmission/rpc", RPCHandler)
	m.Action(r.Handle)
	m.Use(auth.Basic("test", "test"))
	cMux.Handle("/", m)
	cServer = httptest.NewServer(cMux)

	// github client configured to use test server
	client = NewClient(cServer.URL+"/transmission/rpc", "test", "test")
}

func cTeardown() {
	cServer.Close()
}

func TestPost(t *testing.T) {
	cSetup()
	defer cTeardown()
	Convey("Test Post is working correctly", t, func() {
		output, err := client.Post("")
		So(err, ShouldBeNil)
		So(string(output), ShouldEqual, `{"arguments":{},"result":"no method name"}`)
	})

	Convey("Test when auth is incorrect", t, func() {
		fakeClient := NewClient(cServer.URL, "testfake", "testfake")
		output, err := fakeClient.Post("")
		So(err, ShouldBeNil)
		So(string(output), ShouldEqual, "Not Authorized\n")
	})

}
