This is a simple web server that accepts web requests and returns useful information for testing purposes.

1. IP address of the client making web requst

`$ curl 127.0.0.1
127.0.0.1`


2. Useful information in the HTTP header of the web request received

`$ curl 127.0.0.1/info
You are visitor number 2 since last restart

You requested to:
GET

Your requested URL is:
127.0.0.1/info

Your request is from this IP address:
127.0.0.1:34148

Your original IP address is:
172.17.0.1

Your User-Agent is:
curl/7.58.0`


3. Full header of the HTTP request

`$ curl 127.0.0.1/info?req_hdr=true
You are visitor number 1 since last restart

You requested to:
GET

Your requested URL is:
127.0.0.1/info

Query in your request:
req_hdr=true

Your request is from this IP address:
127.0.0.1:34142

Your original IP address is:
172.17.0.1

Your User-Agent is:
curl/7.58.0


The full request header is:
X-Forwarded-For: 172.17.0.1
Accept-Encoding: gzip
Connection: close
User-Agent: curl/7.58.0
Accept: */*`


4. Simulate load and process time on the web service

sleep=500 means the web service will pause for 500ms.  No CPU consumption during this period.
load=500 means the web service will enter a loop for 500ms.  CPU consumption will be high during this period.

`$ curl "127.0.0.1:8888/info?sleep=500&load=500"
You are visitor number 6 since last restart

You requested to:
GET

Your requested URL is:
127.0.0.1:8888/info

Query in your request:
sleep=500&load=500

Your request is from this IP address:
127.0.0.1:34298

Your original IP address is:
172.17.0.1

Your User-Agent is:
curl/7.58.0`

