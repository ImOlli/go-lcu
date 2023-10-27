package proxy

import (
	"crypto/tls"
	"github.com/ImOlli/go-lcu/lcu"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type LCUProxy struct {
	// The hostname and port of the proxy
	host string

	// The hostname and port of the LCU.
	LCUHost string

	// The authentication token of the LCU.
	AuthToken string

	// Indicates if the startup message should be shown
	DisableStartUpMessage bool

	// Indicates if CORS should be disabled. When this property is set to true.
	// The Proxy will respond with the Header (Access-Control-Allow-Origin: *)
	DisableCORS bool

	// Indicates that the certificate check should be disabled
	DisableCertCheck bool
}

// CreateProxy creates a new reverse proxy which can be used to get rid of https or/and the self-signed certificate.
// Also, it can be used to disable cors for local web development
//
// This function calls internally the lcu.FindLCUConnectInfo to retrieve the auth-token and port of the LCU.
func CreateProxy(host string) (*LCUProxy, error) {
	info, err := lcu.FindLCUConnectInfo()

	if err != nil {
		return nil, err
	}

	return CreateCustomProxy(host, "127.0.0.1:"+info.Port, info.AuthToken), nil
}

// CreateCustomProxy creates a new reverse proxy which can be used to get rid of https or/and the self-signed certificate.
// Also, it can be used to disable cors for local web development
//
// Other than the function CreateProxy this method doesn't call the lcu.FindLCUConnectInfo function but depends
// on the specific auth-token and port
func CreateCustomProxy(host string, lcuHost string, authToken string) *LCUProxy {
	var proxy = &LCUProxy{
		host:             host,
		LCUHost:          lcuHost,
		AuthToken:        authToken,
		DisableCertCheck: true,
	}

	handler := func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			r.Host = proxy.LCUHost
			r.SetBasicAuth("riot", proxy.AuthToken)
			p.ServeHTTP(w, r)
		}
	}

	httpProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Host:   proxy.LCUHost,
		Scheme: "https",
	})
	httpProxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: proxy.DisableCertCheck},
	}
	httpProxy.ModifyResponse = func(response *http.Response) error {
		if proxy.DisableCORS {
			response.Header.Set("Access-Control-Allow-Origin", "*")
		}

		return nil
	}
	http.HandleFunc("/", handler(httpProxy))

	return proxy
}

// ListenAndServe Starts the proxy and waits for incoming requests.
//
// To disable the startup message set DisableStartupMessage to true in the LCUProxy struct
func (proxy *LCUProxy) ListenAndServe() error {
	if !proxy.DisableStartUpMessage {
		log.Println("Starting LCU-Proxy on " + proxy.host)
	}

	return http.ListenAndServe(proxy.host, nil)
}
