This is a simple web server that accepts web requests and returns useful information for testing purposes.

You can build an image by yourself using the sample Dockerfile in this repo or to pull the image from fefefe8888/go-web-req-info.


1. IP address of the client making web requst

```
$ curl 172.16.0.1
172.17.0.1
```


2. Useful information in the HTTP header of the web request received

```
$ curl 172.16.0.1/info
You are visitor number 2 since last restart

You requested to:
GET

Your requested URL is:
172.16.0.1/info

Your request is from this IP address:
127.0.0.1:34148

Your original IP address is:
172.17.0.1

Your User-Agent is:
curl/7.58.0
```


3. Full header of the HTTP request

```
$ curl 172.16.0.1/info?req_hdr=true
You are visitor number 1 since last restart

You requested to:
GET

Your requested URL is:
172.16.0.1/info

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
Accept: */*
```


4. Simulate load and process time on the web service

- sleep=500 means the web service will pause for 500ms.  No CPU consumption during this period.
- load=500 means the web service will enter a loop for 500ms.  CPU consumption will be high during this period.

```
$ curl "172.16.0.1/info?sleep=500&load=500"
You are visitor number 6 since last restart

You requested to:
GET

Your requested URL is:
172.16.0.1/info

Query in your request:
sleep=500&load=500

Your request is from this IP address:
127.0.0.1:34298

Your original IP address is:
172.17.0.1

Your User-Agent is:
curl/7.58.0
```


5. Control the size of the HTTP response

- data=1000 adds 1,000 bytes to the HTTP response

```
$ curl "172.16.0.1/?data=1000"
172.17.0.1
===== processing data parameter =====
Packing 1000 bytes of data
aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
```


6. Generate a DNS lookup on the server

- domain=www.google.com triggers the server to make a DNS lookup for www.google.com

```
$ curl "172.16.0.1/?domain=www.google.com"
172.17.0.1
===== processing domain parameter =====
Looking up domain: www.google.com
www.google.com resolves to 142.250.31.103
www.google.com resolves to 142.250.31.147
www.google.com resolves to 142.250.31.105
www.google.com resolves to 142.250.31.106
www.google.com resolves to 142.250.31.99
www.google.com resolves to 142.250.31.104
www.google.com resolves to 2607:f8b0:4004:c08::93
www.google.com resolves to 2607:f8b0:4004:c08::68
www.google.com resolves to 2607:f8b0:4004:c08::67
www.google.com resolves to 2607:f8b0:4004:c08::6a
```


7. Download malware sample from WildFire

```
$ curl "172.16.0.1/?cmd=download"
```


8. Trigger information leak event

```
$ curl "172.16.0.1/?cmd=leak"
```


9. Trigger unexpected listening port event

```
$ curl "172.16.0.1/?cmd=nc"
```


10. Trigger modified binary event

```
$ curl "172.16.0.1/?cmd=modified"
```
