// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

/*
  HTTP server imitating a proxy.
  Doesn't forward requests anywhere, instead responds with status 200 to CONNECT requests and status 418 to everything else.
  All incoming requests are printed to stdout, log messages to stderr.

  By default listens on TCP port 64086 - this can be changed with a command-line argument:

    dummyproxy --port <number>

  If you pass port number 0, the proxy will run on a random port. The actual listening port will be written to stderr, e.g.:

    $ dummyproxy --port 0
    port: 54321

  As this is intended for automated tests, the server will only handle one connection and then exit. There's a test helper
  (hopgoblin/http/dummyproxy.DummyProxyRunner) that simplifies usage in tests.

  See main() function on how to embed this server into another application without starting it in a separate process.
 */
package main

import (
  "flag"
  "fmt"
  "github.com/rinusser/hopgoblin/bootstrap"
  _ "github.com/rinusser/hopgoblin/log/appconfig" //keep: allows log configuration in application.ini
)


var port=flag.Int("port",64086,"TCP port to listen on")


func main() {
  bootstrap.Init()
  proxy:=NewHTTPProxy()
  proxy.Listen(findPort())
}


func findPort() int {
  if *port<0 || *port>65535 {
    panic(fmt.Sprintf("got invalid port number '%d'",*port))
  }
  return *port
}
