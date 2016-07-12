package transmission

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/auth"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	tMux               *http.ServeMux
	transmissionClient *TransmissionClient
	tServer            *httptest.Server
)

func tSetup(output string) {
	tMux = http.NewServeMux()
	tServer = httptest.NewServer(tMux)
	m := martini.New()
	r := martini.NewRouter()
	r.Post("/transmission/rpc", func() string {
		return output
	})
	m.Action(r.Handle)
	m.Use(auth.Basic("test", "test"))
	tMux.Handle("/", m)

	transmissionClient, _ = New(tServer.URL+"/transmission/rpc", "test", "test")
}

func tTeardown() {
	tServer.Close()
}

func TestGetTorrents(t *testing.T) {
	tSetup(`{"arguments":{"torrents":[{"eta":-1,"id":5,
  "leftUntilDone":0,"name":"Test",
  "rateDownload":0,"rateUpload":0,"status":6,"uploadRatio":0.3114}]},
  "result":"success"}`)
	defer tTeardown()

	Convey("Test get list torrents", t, func() {
		torrents, err := transmissionClient.GetTorrents()
		So(err, ShouldBeNil)
		So(len(torrents), ShouldEqual, 1)
	})
}

func TestRemoveTorrent(t *testing.T) {
	tSetup(`{"arguments":{"torrents":[{"eta":-1,"id":5,
  "leftUntilDone":0,"name":"Test",
  "rateDownload":0,"rateUpload":0,"status":6,"uploadRatio":0.3114}]},
  "result":"success"}`)
	defer tTeardown()

	Convey("Test removing torrent", t, func() {
		name, err := transmissionClient.DeleteTorrent(5, true)
		So(err, ShouldBeNil)
		So(name, ShouldEqual, "Test")
	})
}

func TestAddTorrentByFilename(t *testing.T) {
	tSetup(`{"arguments":{"torrent-added":
  {"hashString":"875a2d90068c32b4ce7992eaf56cd03f5be0d193",
  "id":23,"name":"Test Name"}}
  ,"result":"success"}`)
	defer tTeardown()

	Convey("Test adding torrent by filename", t, func() {
		addCmd := NewAddCmdByFilename("/tmp/file")

		result, err := transmissionClient.ExecuteAddCommand(addCmd)
		So(err, ShouldBeNil)

		So(result.Name, ShouldEqual, "Test Name")
		So(result.ID, ShouldEqual, 23)
	})
}

func TestAddTorrentByMagnet(t *testing.T) {
	tSetup(`{"arguments":{"torrent-added":
  {"hashString":"875a2d90068c32b4ce7992eaf56cd03f5be0d193",
  "id":23,"name":"CentOS 7.0 x64"}}
  ,"result":"success"}`)
	defer tTeardown()

	Convey("Test adding torrent by magnet link", t, func() {
		addCmd := NewAddCmdByURL("magnet:?xt=urn:btih:1354ac45bfb3e644a04d69cc519e83283bd3ac6a&dn=CentOS+7.0+x64&tr=udp%3A%2F%2Ftracker.openbittorrent.com%3A80")

		result, err := transmissionClient.ExecuteAddCommand(addCmd)
		So(err, ShouldBeNil)

		So(result.Name, ShouldEqual, "CentOS 7.0 x64")
		So(result.ID, ShouldEqual, 23)
	})
}
