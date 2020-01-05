# Go PASSWORD HASH REST API Exercise
A RESTful API for simple password hash service with Go

## Installation & Run
```bash
# Download this project
go get github.com/mcasado/hashservice
```

```bash
# Build and Run
cd hashservice
go build
./hashservice --listen-addr 3000

# API Endpoint : http://127.0.0.1:3000
```

## Structure
```
├── app
│   ├── handlers.go       // API core handlers
│   ├── comtroller.go     // Providing state to handler
│   ├── middleware.go     // Logging and tracing wrapping handler
│   ├── storage.go        // In memory state storage
│   ├── persist.go        // Our API core handlers
│   ├── router.go         // Request Routing 
│   ├── server.go         // Http server 
│   └── util.go           // Utilities
└── main.go               // main driver
```

## API

## `POST /hash`

#### Creates a base64 encoded password hash given the password in a form field. The request returns an id right away that allows to retrieve the passwrod hash with a get request   

#### Request
```
    curl  -v http://localhost:8000/hash  -X POST -d 'password=angryMonkey'
```

#### Response
```
> POST /hash HTTP/1.1
> User-Agent: curl/7.24.0 (x86_64-apple-darwin11.2.0) libcurl/7.24.0 OpenSSL/1.0.2a zlib/1.2.8 libidn/1.22
> Host: localhost:8000
> Accept: */*
> Content-Length: 20
> Content-Type: application/x-www-form-urlencoded
>
* upload completely sent off: 20 out of 20 bytes
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=utf-8
< X-Request-Id: 1578203585383597000
< Date: Sun, 05 Jan 2020 05:53:05 GMT
< Content-Length: 1
<
* Connection #0 to host localhost left intact
1
```

## `GET /hash/:id`
#### Retrieve the base64 encoded password hash given the id returned when it was created 

#### Request
```
    curl -i -H 'Accept: text/plain' http://localhost:3000/hash/1
```

#### Response
```
 > GET /hash/1 HTTP/1.1
 > User-Agent: curl/7.24.0 (x86_64-apple-darwin11.2.0) libcurl/7.24.0 OpenSSL/1.0.2a zlib/1.2.8 libidn/1.22
 > Host: localhost:8000
 > Accept: */*
 >
 < HTTP/1.1 200 OK
 < Content-Type: text/plain; charset=utf-8
 < X-Request-Id: 1578202471685844000
 < Date: Sun, 05 Jan 2020 05:34:31 GMT
 < Content-Length: 88
 <
 * Connection #0 to host localhost left intact
 ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==*
```

## `GET /stats`
#### Provides a statistics endpoint to get basic information about the password hashes requests.

#### Request
```
    curl -i -H 'Accept: text/plain' http://localhost:3000/stats | jq .
```

#### Response
```
< HTTP/1.1 200 OK
< Content-Type: application/json
< X-Request-Id: 1578204141824410000
< Date: Sun, 05 Jan 2020 06:02:21 GMT
< Content-Length: 635
<
{ [data not shown]
100   635  100   635    0     0   141k      0 --:--:-- --:--:-- --:--:--  620k
* Connection #0 to host localhost left intact
* Closing connection #0
{
  "pid": 23416,
  "uptime": "29m50.481843542s",
  "uptime_sec": 1790.481843542,
  "time": "2020-01-04 23:02:21.824416 -0700 MST m=+1790.482455838",
  "unixtime": 1578204141,
  "total_method_status_code_count": {
    "GET:hash": {
      "200": 8
    },
    "GET:stats": {
      "200": 2
    },
    "POST:hash": {
      "200": 3
    }
  },
  "total_count": 13,
  "total_method_response_time_sec": {
    "GET:hash": {
      "200": 0.000522117
    },
    "GET:stats": {
      "200": 0.00032646
    },
    "POST:hash": {
      "200": 0.000182051
    }
  },
  "total_response_time_sec": 0.001030628,
  "average_method_response_time_sec": {
    "GET:hash": {
      "200": 6.5264625e-05
    },
    "GET:stats": {
      "200": 0.00016323
    },
    "POST:hash": {
      "200": 6.068366666666666e-05
    }
  },
  "average_response_time_sec": 7.927907692307693e-05
}
```

## `GET /health`
#### Provides a health endpoint for the service 

#### Request
```
    curl -i -H 'Accept: text/plain' http://localhost:3000/health | jq .
```

#### Response
```
< HTTP/1.1 200 OK
< Content-Type: application/json; charset=utf-8
< X-Content-Type-Options: nosniff
< X-Request-Id: 1578204339502429000
< Date: Sun, 05 Jan 2020 06:05:39 GMT
< Content-Length: 15
<
{ [data not shown]
100    15  100    15    0     0   2963      0 --:--:-- --:--:-- --:--:-- 15000
* Connection #0 to host localhost left intact
* Closing connection #0
{
  "alive": true
}
```

## `GET /shutdown`
#### Provides a `graceful` service shutdown endpoint for the service 

#### Request
```
    curl -i -H 'Accept: text/plain' http://localhost:3000/shutdown
```

#### Response
```
 HTTP/1.1 200 OK
< X-Request-Id: 1578204563827789000
< Date: Sun, 05 Jan 2020 06:09:23 GMT
< Content-Length: 17
< Content-Type: text/plain; charset=utf-8
<
* Connection #0 to host localhost left intact
shutting down ...
```

## Todo

- [x] Support basic REST APIs.
- [x] Graceful shutdown
- [x] Health endpoint
- [ ] Write the tests for all APIs.
- [x] Organize the code with packages
- [ ] Make docs with GoDoc
- [ ] Building a deployment process 