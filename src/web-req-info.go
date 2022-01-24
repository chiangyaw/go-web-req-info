package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var visit_count int = 0

// Server to return info seen in the HTTP request
func WebInfoServer(w http.ResponseWriter, r *http.Request) {

	var show_req_header bool = false

	// count number of visitors since last restart
	visit_count += 1

	// show basic info in log
	log.Printf("Received request from %s for path: %s from %s", r.RemoteAddr, r.URL.Path, r.Header.Get("User-Agent"))

	if r.URL.Path == "/" {
		// show IP address of the client only if the requested URL is "/"
		print_sender_ip(w, r)
	} else if r.URL.Path == "/info" {
		// show more info about the request if the requested URL is "/info"
		print_req_info(w, r)
	} else {
		fmt.Fprintf(w, "Invalid URL")
	}

	// check whether query parameters exist
	if len(r.URL.RawQuery) > 0 {

		// req_hdr exists - return full HTTP request header
		show_req_header, _ = strconv.ParseBool(r.URL.Query().Get("req_hdr"))

		// sleep exists - sleep for the duration of the received value
		if sleep_duration_ms, err := strconv.Atoi(r.URL.Query().Get("sleep")); err == nil {
			delay_response(w, sleep_duration_ms)
		}

		// load exists - loop for the duration of the received value
		if load_duration_ms, err := strconv.Atoi(r.URL.Query().Get("load")); err == nil {
			add_load(w, load_duration_ms)
		}

		// domain exists - lookup IP address of the domain name
		if domain_name := r.URL.Query().Get("domain"); len(domain_name) > 0 {
			lookup_domain(w, domain_name)
		}

		// cmd exists - run specific commands
		if cmd_opt := r.URL.Query().Get("cmd"); len(cmd_opt) > 0 {
			run_cmd(w, cmd_opt)
		}
	}

	// show full HTTP request header
	if show_req_header {
		fmt.Fprintf(w, "\nThe full request header is:\n")
		for key, value := range r.Header {
			for _, element := range value {
				fmt.Fprintf(w, "%s: %s\n", key, element)
			}
		}
		fmt.Fprintf(w, "\n")
	}
}

// print IP address of HTTP request sender
func print_sender_ip(w http.ResponseWriter, r *http.Request) {

	var xff string = r.Header.Get("X-Forwarded-For")
	var requester_ip_port string = r.RemoteAddr

	if len(xff) > 0 {
		ips := strings.Split(xff, ", ")
		fmt.Fprintf(w, "%s", ips[0])
	} else {
		ip_port_slice := strings.Split(requester_ip_port, ":")
		fmt.Fprintf(w, "%s", strings.TrimSuffix(requester_ip_port, ":"+ip_port_slice[len(ip_port_slice)-1]))
	}
}

// print more info about the request if the requested URL is "/info"
func print_req_info(w http.ResponseWriter, r *http.Request) {

	var user_agent string = r.Header.Get("User-Agent")
	var cookies string = r.Header.Get("Cookie")
	var requester_ip_port string = r.RemoteAddr
	var xff string = r.Header.Get("X-Forwarded-For")
	var xfh string = r.Header.Get("X-Forwarded-Host")
	var query string = r.URL.RawQuery

	// show visitor count since last restart
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
	if len(r.URL.RawQuery) > 0 {
		fmt.Fprintf(w, "Query in your request:\n%s\n\n", query)
	}

	// show source IP address of IP packet
	fmt.Fprintf(w, "Your request is from this IP address:\n%s\n\n", requester_ip_port)

	// show IP addresses of the client and intermediate nodes
	print_xff_ips(w, xff)

	// show cookie info
	if len(cookies) > 0 {
		fmt.Fprintf(w, "Your Cookies are:\n%s\n\n", cookies)
	}

	// show User-Agent info
	if len(user_agent) > 0 {
		fmt.Fprintf(w, "Your User-Agent is:\n%s\n\n", user_agent)
	}
}

// print IP addresses from XFF header
func print_xff_ips(w http.ResponseWriter, xff string) {

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
}

// delay reponse to HTTP request by sleeping for a duration specified by the sleep parameter
func delay_response(w http.ResponseWriter, sleep_duration_ms int) {

	log.Printf("Sleeping for %v ms", sleep_duration_ms)
	time.Sleep(time.Duration(sleep_duration_ms) * time.Millisecond)
}

// add load when processing the HTTP request for a duration specified by the load parameter
func add_load(w http.ResponseWriter, load_duration_ms int) {

	var i int64 = 0

	log.Printf("Consuming CPU for %v ms", load_duration_ms)
	for start := time.Now(); time.Since(start) < (time.Duration(load_duration_ms) * time.Millisecond); {
		i++
	}
}

// lookup domain name based on the value of the domain parameter
func lookup_domain(w http.ResponseWriter, domain_name string) {

	ips, err := net.LookupIP(domain_name)
	if err == nil {
		for _, ip := range ips {
			log.Printf("%v resolves to %v", domain_name, ip)
		}
	} else {
		log.Printf("domain lookup error: %v", err)
	}
}

// run command based on the value of the cmd parameter
func run_cmd(w http.ResponseWriter, cmd_opt string) {

	fmt.Fprintf(w, "\n===== processing cmd parameter =====\n")

	switch cmd_opt {

	case "download":
		// run curl command to download sample malware file and trigger unexpected process, DNS, file system access and malware event
		cmd_out, err := exec.Command("curl", "http://sg.wildfire.paloaltonetworks.com/publicapi/test/elf", "-o", "/tmp/malware-sample").Output()
		fmt.Fprintf(w, "Running command \"curl http://sg.wildfire.paloaltonetworks.com/publicapi/test/elf\"\n")
		if err == nil {
			fmt.Fprintf(w, "%s\n", cmd_out)
		} else {
			fmt.Fprintf(w, "%s\n", err)
		}
	case "leak":
		// run cat /etc/passwd to trigger information leakage event
		cmd_out, err := exec.Command("cat", "/etc/passwd").Output()
		fmt.Fprintf(w, "Running command \"cat /etc/passwd\"\n")
		if err == nil {
			fmt.Fprintf(w, "%s\n", cmd_out)
		} else {
			fmt.Fprintf(w, "%s\n", err)
		}
	case "listen":
		// run nc command to trigger unexpected listening port event
		ctx_duration := 60 * time.Second
		tcp_port := "11111"

		ctx, cancel := context.WithTimeout(context.Background(), ctx_duration)
		defer cancel()

		cmd_out, err := exec.CommandContext(ctx, "nc", "-lvp", tcp_port).Output()
		fmt.Fprintf(w, "Running command \"nc -lvp %s\"\n", tcp_port)
		if err == nil {
			fmt.Fprintf(w, "%s\n", cmd_out)
		} else {
			fmt.Fprintf(w, "%s\n", err)
		}
	case "modified":
		// touch a binary and run it to trigger modified binary event
		cmd_out, err := exec.Command("touch", "/bin/ls").Output()
		fmt.Fprintf(w, "Running command \"touch /bin/ls\"\n")
		if err == nil {
			fmt.Fprintf(w, "%s\n", cmd_out)
		} else {
			fmt.Fprintf(w, "%s\n", err)
		}
		cmd_out, err = exec.Command("ls").Output()
		fmt.Fprintf(w, "Running command \"ls\"\n")
		if err == nil {
			fmt.Fprintf(w, "%s\n", cmd_out)
		} else {
			fmt.Fprintf(w, "%s\n", err)
		}
	default:
		fmt.Fprintf(w, "Invalid value: %s\n", cmd_opt)
	}
}

func main() {

	var addr string = ":80"
	handler := http.HandlerFunc(WebInfoServer)

	log.Printf("Starting webserver on %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Could not listen on port %s %v", addr, err)
	}
}
