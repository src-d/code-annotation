package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/src-d/code-annotation/server/assets"

	"github.com/koding/websocketproxy"
)

// Static contains handlers to serve static using go-bindata
type Static struct {
	dir         string
	proxyTarget string
	options     options
}

// NewStatic creates new Static
func NewStatic(dir, serverURL, gaTrackingID string, proxyTarget string) *Static {
	return &Static{
		dir:         dir,
		proxyTarget: proxyTarget,
		options: options{
			ServerURL:    serverURL,
			GaTrackingID: gaTrackingID,
		},
	}
}

// struct which will be marshalled and exposed to frontend
type options struct {
	ServerURL    string      `json:"SERVER_URL"`
	GaTrackingID string      `json:"GA_TRACKING_ID"`
	InitalState  interface{} `json:"initialState"`
}

// ServeHTTP serves any static file from static directory or fallbacks on index.hml
func (s *Static) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.proxyTarget != "" {
		s.proxy(w, r)
		return
	}

	filepath := s.dir + r.URL.Path
	b, err := assets.Asset(filepath)
	if err != nil {
		s.ServeIndexHTML(nil)(w, r)
		return
	}
	s.serveAsset(w, r, filepath, b)
}

// ServeIndexHTML serves index.html file with initial state
func (s *Static) ServeIndexHTML(initialState interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filepath := s.dir + "/index.html"
		b, err := assets.Asset(filepath)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		options := s.options
		options.InitalState = initialState
		bData, err := json.Marshal(options)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b = bytes.Replace(b, []byte("window.REPLACE_BY_SERVER"), bData, 1)
		s.serveAsset(w, r, filepath, b)
	}
}

func (s *Static) serveAsset(w http.ResponseWriter, r *http.Request, filepath string, content []byte) {
	info, err := assets.AssetInfo(filepath)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	http.ServeContent(w, r, info.Name(), info.ModTime(), bytes.NewReader(content))
}

func (s *Static) proxy(w http.ResponseWriter, r *http.Request) {
	upgrade := r.Header["Upgrade"]
	if len(upgrade) > 0 && upgrade[0] == "websocket" {
		target, _ := url.Parse("ws://" + s.proxyTarget)
		websocketproxy.ProxyHandler(target).ServeHTTP(w, r)
		return
	}
	target, _ := url.Parse("http://" + s.proxyTarget)
	httputil.NewSingleHostReverseProxy(target).ServeHTTP(w, r)
}
