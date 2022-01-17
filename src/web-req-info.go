package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var visit_count int = 0

// Server to return info seen in the HTTP request
func WebInfoServer(w http.ResponseWriter, r *http.Request) {

	var user_agent string = r.Header.Get("User-Agent")
	var cookies string = r.Header.Get("Cookie")
	var requester_ip_port string = r.RemoteAddr
	var xff string = r.Header.Get("X-Forwarded-For")
	var xfh string = r.Header.Get("X-Forwarded-Host")
	var query string = r.URL.RawQuery
	var query_exists bool = false
	var show_req_header bool = false

	// check whether req_hdr API is called
	if len(query) > 0 {
		query_exists = true
		show_req_header, _ = strconv.ParseBool(r.URL.Query().Get("req_hdr"))
	}

	// show IP address of the client only if the requested URL is "/"
	if r.URL.Path == "/" {
		if len(xff) > 0 {
			ips := strings.Split(xff, ", ")
			fmt.Fprintf(w, "%s", ips[0])
		} else {
			ip_port_slice := strings.Split(requester_ip_port, ":")
			fmt.Fprintf(w, "%s", strings.TrimSuffix(requester_ip_port, ":"+ip_port_slice[len(ip_port_slice)-1]))
		}
		// show more info about the request if the requested URL is not "/"
	} else {

		// show visitor count since last restart
		visit_count += 1
		fmt.Fprintf(w, "You are visitor number %d", visit_count)
		fmt.Fprintf(w, " since last restart\n\n")

		// show HTTP method
		fmt.Fprintf(w, "You requested to:\n%s\n\n", r.Method)

		// show HOST info
		if len(xfh) > 0 {
			fmt.Fprintf(w, "Your original request is for the host:\n%s\n\n", xfh)
		}

		// show URL requested
		fmt.Fprintf(w, "Your requested URL is:\n%s%s\n\n", r.Host, r.URL.Path)

		// show query
		if query_exists {
			fmt.Fprintf(w, "Query in your request:\n%s\n\n", query)
		}

		// show source IP address of IP packet
		fmt.Fprintf(w, "Your request is from this IP address:\n%s\n\n", requester_ip_port)

		// show IP addresses of the client and intermediate nodes
		if len(xff) > 0 {
			ips := strings.Split(xff, ", ")
			for i, ip := range ips {
				if i == 0 {
					fmt.Fprintf(w, "Your original IP address is:\n")
					fmt.Fprintf(w, "%s\n", ip)
				} else if i == 1 {
					fmt.Fprintf(w, "Your request is going through the following IP addresses:\n")
					fmt.Fprintf(w, "%s\n", ip)
				} else {
					fmt.Fprintf(w, "%s\n", ip)
				}
			}
			fmt.Fprintf(w, "\n")
		}

		// show cookie info
		if len(cookies) > 0 {
			fmt.Fprintf(w, "Your Cookies are:\n%s\n\n", cookies)
		}

		// show User-Agent info
		if len(user_agent) > 0 {
			fmt.Fprintf(w, "Your User-Agent is:\n%s\n\n", user_agent)
		}
	}

	// show full HTTP request header
	if show_req_header {
		fmt.Fprintf(w, "\n\nThe full request header is:\n")
		for key, value := range r.Header {
			for _, element := range value {
				fmt.Fprintf(w, "%s: %s\n", key, element)
			}
		}
		fmt.Fprintf(w, "\n")
	}

	// show basic info in log
	log.Printf("Received request from %s for path: %s from %s", requester_ip_port, r.URL.Path, user_agent)
}

func main() {

	var addr string = ":80"
	handler := http.HandlerFunc(WebInfoServer)

	log.Printf("Starting webserver on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Could not listen on port %s %v", addr, err)
	}
}
