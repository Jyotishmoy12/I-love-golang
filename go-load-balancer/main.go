// main package defines the entry point for the Go program
package main

import (
	"fmt"                 // for formatted I/O
	"net/http/httputil"   // for reverse proxy logic
	"net/url"             // for parsing backend URLs
	"os"                  // for program exit on fatal error
	"net/http"            // for HTTP request and response handling
)

// Server interface defines behavior for any backend server
type Server interface {
	Address() string                                 // returns server address (URL)
	IsAlive() bool                                   // returns if server is healthy
	Serve(w http.ResponseWriter, r *http.Request)    // handles HTTP request using proxy
}

// simpleServer is a basic implementation of the Server interface
type simpleServer struct {
	addr  string                   // backend server address
	proxy *httputil.ReverseProxy  // reverse proxy to forward requests to backend
}

// newSimpleServer creates a new simpleServer with a reverse proxy to the given address
func newSimpleServer(addr string) *simpleServer {
	serverUrl, err := url.Parse(addr)          // parse the address string into a URL object
	handleErr(err)                             // exit if parsing fails

	return &simpleServer{
		addr:  addr,                             // store the address
		proxy: httputil.NewSingleHostReverseProxy(serverUrl), // create a proxy to this backend
	}
}

// loadBalancer is the main structure that distributes traffic across servers
type loadBalancer struct {
	port            string    // port the load balancer listens on
	roundRobinCount int       // counter to track round-robin position
	servers         []Server  // list of servers to distribute requests to
}

// NewLoadBalancer constructs and returns a new loadBalancer
func NewLoadBalancer(port string, servers []Server) *loadBalancer {
	return &loadBalancer{
		roundRobinCount: 0,   // initialize counter
		port:            port,
		servers:         servers,
	}
}

// handleErr checks if an error exists, and if so, prints and exits the program
func handleErr(err error) {
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}

// getNextAvailableServer selects the next healthy server using round-robin
func (lb *loadBalancer) getNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)] // pick server based on current count
	for !server.IsAlive() {                                   // if server is not alive
		// try next server
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++ // increment counter for next round
	return server        // return selected healthy server
}

// serveProxy forwards the incoming HTTP request to the selected backend server
func (lb *loadBalancer) serveProxy(w http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer() // select backend server
	fmt.Printf("redirecting request to %s\n", targetServer.Address()) // log redirection
	targetServer.Serve(w, r) // forward request using the proxy
}

// Address returns the address of the simpleServer
func (s *simpleServer) Address() string {
	return s.addr
}

// IsAlive checks if the server is healthy
// This is a dummy check (always returns true), but can be replaced with real health logic
func (s *simpleServer) IsAlive() bool {
	return true
}

// Serve handles the HTTP request by passing it to the reverse proxy
func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req) // reverse proxy handles the request
}

// main is the entry point of the application
func main() {
	// Define list of backend servers
	servers := []Server{
		newSimpleServer("https://google.com/"),
		newSimpleServer("https://facebook.com/"),
		newSimpleServer("https://openai.com/"),
	}

	// Create a load balancer on port 8080 with the list of servers
	lb := NewLoadBalancer(":8080", servers)

	// Define HTTP handler function to forward requests
	handleRedirect := func(rw http.ResponseWriter, req *http.Request) {
		lb.serveProxy(rw, req) // serve each request via load balancer
	}

	// Register the handler function for the root path "/"
	http.HandleFunc("/", handleRedirect)

	// Log that the load balancer is running
	fmt.Printf("serving requests at 'localhost%s'\n", lb.port)

	// Start the HTTP server on specified port
	http.ListenAndServe(lb.port, nil)
}
