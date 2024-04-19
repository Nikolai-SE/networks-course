package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"nikolai/proxy-server/internal/blacklist"
	"nikolai/proxy-server/internal/cache"
)

var (
	logFileName = "proxy-%v.log"
	transport   = http.DefaultTransport

	CacheControl = "Cache-Control"
	LastModified = "Last-Modified"
	Etag         = "Etag"
)

type Proxy struct {
	address   string
	port      int
	log       *log.Logger
	cache     cache.Cache
	blackList *[]blacklist.Forbiden
}

func NewProxy(address string, port int, w io.Writer, blackList *[]blacklist.Forbiden) Proxy {
	return Proxy{
		address:   address,
		port:      port,
		log:       log.New(w, fmt.Sprintf("%s:%d", address, port), log.Ldate|log.Ltime),
		cache:     cache.NewCache(),
		blackList: blackList,
	}
}

// ServeHTTP implements http.Handler.
func (p Proxy) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	if blacklist.IsInBlackList(request, *p.blackList) {
		defer request.Body.Close()
		w.WriteHeader(http.StatusForbidden)
		p.log.Printf("Blacklist (method: %v, url: %v)\n", request.Method, request.URL)
		return
	}

	switch m := request.Method; m {
	case http.MethodGet:
		p.serveHTTPGet(w, request)
	case http.MethodPost:
		p.serveHTTPPost(w, request)
	default:
		p.log.Printf("Tunneling (method: %v, url: %v)\n", m, request.URL)
		handleTunneling(w, request)
	}
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(dest_conn, client_conn)
	go transfer(client_conn, dest_conn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func (p Proxy) serveHTTPGet(w http.ResponseWriter, r *http.Request) {
	proxyReq, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy GET request", http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	defer proxyReq.Body.Close()

	// Copy the headers from the original request to the proxy request
	copyHeaders(r.Header, proxyReq.Header)

	cachedData := p.CheckInHash(proxyReq)

	// Send the proxy request using the custom transport
	resp, err := transport.RoundTrip(proxyReq)
	if err != nil {
		http.Error(w, "Error sending proxy GET request", http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	defer resp.Body.Close()

	// Copy the headers from the proxy response to the original response
	copyHeaders(resp.Header, w.Header())

	if cachedData != nil && resp.StatusCode == http.StatusNotModified {
		w.WriteHeader(http.StatusOK)
		w.Write(*cachedData)
		p.log.Printf("Cache match. Method: %s, Url: %s, Status: %d (%s)\n",
			r.Method, r.RequestURI, resp.StatusCode, resp.Status)
		return
	}

	// Set the status code of the original response to the status code of the proxy response
	w.WriteHeader(resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error read GET respounce", http.StatusInternalServerError)
		fmt.Fprintln(w, err)
	}
	w.Write(data)

	p.SaveToCache(r, resp, data)

	p.log.Printf("Method: %s, Url: %s, Status: %d (%s)\n", r.Method, r.RequestURI, resp.StatusCode, resp.Status)
}

func (p Proxy) CheckInHash(r *http.Request) *[]byte {
	cached := p.cache.Get(r.URL.String())
	if cached == nil {
		return nil
	}

	r.Header.Add("If-Modified-Since", cached.ModifiedSince)
	r.Header.Add("If-None-Match", cached.Etag)
	return &cached.Value
}

func (p Proxy) SaveToCache(r *http.Request, resp *http.Response, data []byte) {
	cacheControl := resp.Header.Get(CacheControl)
	etag := resp.Header.Get(Etag)
	lastModified := resp.Header.Get(LastModified)
	if len(lastModified) > 0 && len(etag) > 0 {
		var duration time.Duration = time.Duration(1) * time.Minute // Default
		if len(cacheControl) > 0 {
			v, err := strconv.Atoi(cacheControl)
			if err != nil {
				v = 0
			}

			duration = time.Duration(v) * time.Second
		}
		p.cache.Add(r.URL.String(), data, duration, lastModified, etag)
	}
}

func (p Proxy) serveHTTPPost(w http.ResponseWriter, r *http.Request) {
	// Create a new HTTP request with the same method, URL, and body as the original request
	proxyReq, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}
	defer proxyReq.Body.Close()

	// Copy the headers from the original request to the proxy request
	copyHeaders(r.Header, proxyReq.Header)

	// Send the proxy request using the custom transport
	resp, err := transport.RoundTrip(proxyReq)
	if err != nil {
		http.Error(w, "Error sending proxy POST request", http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}
	defer resp.Body.Close()

	// Copy the headers from the proxy response to the original response
	copyHeaders(resp.Header, w.Header())

	// Set the status code of the original response to the status code of the proxy response
	w.WriteHeader(resp.StatusCode)

	// Copy the body of the proxy response to the original response
	io.Copy(w, resp.Body)

	p.log.Printf("Method: %s, Url: %s, Status: %d (%s)\n", r.Method, r.RequestURI, resp.StatusCode, resp.Status)
}

func copyHeaders(s, d http.Header) {
	for name, values := range s {
		for _, value := range values {
			d.Add(name, value)
		}
	}
}

func RunProxy(address string, port int, blacklistPath string) {
	logFile, err := os.OpenFile(fmt.Sprintf(logFileName, time.Now().Truncate(time.Microsecond)), os.O_RDWR|os.O_SYNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	proxy := NewProxy(address, port, logFile, blacklist.GetBlackList(blacklistPath))
	http.ListenAndServe(fmt.Sprintf("%s:%d", address, port), proxy)
}
