// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

/*
  HTTP server imitating a proxy.
  Doesn't forward requests anywhere, instead responds with status 200 to CONNECT requests and status 418 to everything else.
  All incoming requests are printed to stdout, log messages to stderr.

  By default listens on TCP port 64086 - this can be changed with a command-line argument:

    dummyproxy --port <number>

  As this is intended for automated tests, the server will only handle one connection and then exit. There's a test helper
  (hopgoblin/http/dummyproxy.DummyProxyRunner) that simplifies usage in tests.

  See main() function on how to embed this server into another application without starting it in a separate process.
 */
package main

import (
  "bufio"
  "flag"
  "fmt"
  "net"
  "strings"
  "time"
  "github.com/rinusser/hopgoblin/bootstrap"
  "github.com/rinusser/hopgoblin/http"
  "github.com/rinusser/hopgoblin/log"
  _ "github.com/rinusser/hopgoblin/log/appconfig" //keep: allows log configuration in application.ini
)


var port=flag.Int("port",64086,"TCP port to listen on")


/*
  DummyHTTPProxy type.
 */
type DummyHTTPProxy struct { //TODO: move to http/dummyproxy package
}

/*
  Creates a new proxy instance.
 */
func NewHTTPProxy() *DummyHTTPProxy {
  return &DummyHTTPProxy {}
}


/*
  Starts the dummy proxy, will listen on the passed TCP port, handle a connection and return.
 */
func (proxy *DummyHTTPProxy) Listen(port int) error {
  addr:=net.TCPAddr{net.IPv4(127,0,0,1),port,""}
  listener,err:=net.ListenTCP("tcp",&addr)
  if err!=nil {
    return err
  }

  listener.SetDeadline(time.Now().Add(60e9))
  log.Info("waiting for connection on port %d",port)
  conn,err:=listener.AcceptTCP()
  if err!=nil {
    log.Error("proxy error: ",err)
    return err
  }
  log.Debug("proxy got connection")
  return proxy.handleConnection(conn)
}

func (proxy *DummyHTTPProxy) handleConnection(conn net.Conn) error {
  defer conn.Close()
  buf:=bufio.NewReadWriter(bufio.NewReader(conn),bufio.NewWriter(conn))
  text,err:=http.ReadHTTPMessageAsString(buf)
  if err!=nil {
    return err
  }
  fmt.Println(text)

  if text[0:8]=="CONNECT " {
    buf.WriteString("HTTP/1.1 200 OK\r\n\r\n")
    buf.Flush()

    log.Debug("upgrading to TLS..")
    server:=http.NewServer()
    _,buf,err=server.UpgradeServerConnectionToSSL(conn,"dummyproxy.local")
    if err!=nil {
      return nil
    }

    text,err=http.ReadHTTPMessageAsString(buf)
    if err!=nil {
      return err
    }
    fmt.Println(text)
  }

  if strings.Index(text,"/delayed")>=0 {
    log.Debug("got delay request")
    time.Sleep(2e9)
  }

  buf.WriteString("HTTP/1.1 418 I'm a teapot\r\n")

  if strings.Index(text,"/chunked")>=0 {
    buf.WriteString("Transfer-Encoding: chunked\r\n")
  }

  buf.WriteString("\r\n")

  body:=[]byte(strings.TrimSpace(text))
  if strings.Index(text,"/chunked")>=0 {
    body=http.ChunkEncodeBody(text,1,1)
  }

  buf.Write(body)
//  fmt.Fprintln(os.Stderr,[]byte(body))
//  log.Trace("last 3 bytes in body: %X",body[len(body)-3:])

  buf.Flush()
  return nil
}


func main() {
  bootstrap.Init()
  proxy:=NewHTTPProxy()
  proxy.Listen(findPort())
}


func findPort() int {
  if *port<1 || *port>65535 {
    panic(fmt.Sprintf("got invalid port number '%d'",*port))
  }
  return *port
}
