// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package main

import (
  "bufio"
  "fmt"
  "net"
  "os"
  "strings"
  "time"
  "github.com/rinusser/hopgoblin/http"
  "github.com/rinusser/hopgoblin/log"
)


/*
  DummyHTTPProxy type.

  This would be better suited in the "dummyproxy" package, but needs parts of the "http" package. In package http's tests this
  results in import loops. So either the http package needs to be messed up, or this - better this.
*/
type DummyHTTPProxy struct {
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
    log.Error("proxy listener error: %s",err)
    return err
  }
  port=listener.Addr().(*net.TCPAddr).Port
  fmt.Fprintf(os.Stderr,"port: %d\n",port)

  listener.SetDeadline(time.Now().Add(60e9))
  log.Debug("waiting for connection on port %d",port)
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
