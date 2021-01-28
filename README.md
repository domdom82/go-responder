[![Build Status](https://travis-ci.org/domdom82/go-responder.svg?branch=master)](https://travis-ci.org/domdom82/go-responder)

# go-responder
Configure and test server responses


## Purpose
During development, my programs often talk to remote services that behave weird from time to time. e.g. they respond
nicely for a couple of requests but then take a long time to respond or give 500s all of a sudden or give you the finger
via rate throttling.

Clients usually have to deal with all sorts of weird server behavior. So to make client development easier without having
to rely on the actual backend and its shenanigans, I wrote go-responder.

## How to use
go-responder allows you to simulate the following:
    
- HTTP Servers
- Websocket Servers
- Pure TCP Servers
    
For both http and websocket you can define an arbitrary number of endpoints, each with individual response patterns.
Pure tcp servers don't have endpoints, therefore only ports and read/write behavior on sockets can be configured.

As responses, the following types can be configured:

### Http Servers

```
servers:
  - http:
      port: 8080
      responses:
        /test:
          get:
            static:
              status: 200
              body: all good
          post:
            seq:
              - status: 500
                body: oops that did not work
              - status: 500
                body: try again buddy
              - status: 200
                body: whew now it worked
        /delayed:
          get:
            seq:
              - status: 200
                body: this response took 0.5 seconds
                delay: 500ms
              - status: 200
                body: this response took 1 second
                delay: 1s
        /slow:
          get:
            seq:
              - status: 200
                body: this is a slow response
                rate: 10B/s
              - status: 200
                body: this is even slower
                rate: 1B/s
        /loop:
          get:
            loop:
              - status: 200
                body: one
              - status: 200
                body: two
              - status: 200
                body: three
        /big:
          get:
            static:
              status: 200
              bigbody:
                size: 10M
                type: lorem
```

- Below `responses`, an arbitrary number of paths (e.g. `/alice`, `/bob`, ...) can be configured. The paths must be unique.
- Below the path, http methods `get`, `put`, `post`, `delete` can be configured separately.
- Below the method, one of `static`, `seq` or `loop` can be configured:
    - `static` means the server will always respond with this http response and status
    - `seq` means the server will respond with the next response taken from a sequence with every new request.
    - `loop` is just like `seq` except the sequence will be repeated when it reaches the end.

The actual response has the following properties:

- `status`  the http status code (e.g. 200, 400, 500, 301 etc.)
- `body`  a plain text body
- `bigbody`  an optional large body. used for simulation of large responses.
    - `type` can be one of `binary` (returns random binary data) or `lorem` (returns [lorem ipsum](https://en.wikipedia.org/wiki/Lorem_ipsum) text) 
    - `size` the size of the response to simulate given in human readable format (e.g. `1M`, `100K`, `25B`)
- `delay` an optional delay before the response is transmitted. used for slow server responses. 
          given in human readable format (e.g. `100ms`, `1s`, `1m25s`)
          
### Websocket Servers

```
servers:
  - websocket:
      port: 8082
      responses:
        /ws:
          read:
            bufsize: 1M
            delay: 1s
          write:
            bufsize: 1M
            type: binary
            delay: 1s
        /ws2:
          read:
            bufsize: 1M
          write:
            bufsize: 100B
            type: lorem
            delay: 1s
```

- Below `responses`, an arbitrary number of paths (e.g. `/alice`, `/bob`, ...) can be configured. The paths must be unique.
- Below the path, the socket operations `read` and `write` can be configured.

Websocket servers don't support sequences or loops currently. 

They simply keep reading and writing data with these parameters:

- `bufsize` the buffer size in human readable format (e.g. `1M`, `100K`, `25B`)
- `type` can be one of `binary` (returns random binary data) or `lorem` (returns [lorem ipsum](https://en.wikipedia.org/wiki/Lorem_ipsum) text) 
- `size` the size of the response to simulate given in human readable format (e.g. `1M`, `100K`, `25B`)
- `delay` an optional delay before the response is transmitted. used for slow server responses. 
          given in human readable format (e.g. `100ms`, `1s`, `1m25s`)


### TCP Servers

```
servers:
  - tcp:
      port: 8081
      responses:
        read:
          bufsize: 1M
        write:
          bufsize: 1M
          type: lorem
          delay: 500ms
```

- Below `responses`, the socket operations `read` and `write` can be configured.

TCP servers don't support sequences or loops currently. 

They simply keep reading and writing data with these parameters:

- `bufsize` the buffer size in human readable format (e.g. `1M`, `100K`, `25B`)
- `type` can be one of `binary` (returns random binary data) or `lorem` (returns [lorem ipsum](https://en.wikipedia.org/wiki/Lorem_ipsum) text) 
- `size` the size of the response to simulate given in human readable format (e.g. `1M`, `100K`, `25B`)
- `delay` an optional delay before the response is transmitted. used for slow server responses. 
          given in human readable format (e.g. `100ms`, `1s`, `1m25s`)
          
          
## How to run
Running the *go-responder* simulator is pretty straightforward:

1. `go get github.com/domdom82/go-responder`
1. `cd $GOPATH/src/github.com/domdom82/go-responder`
1. `cp config-template.yml config.yml` and adjust it to your needs.
1. `go build` 
1. `./go-responder` 
