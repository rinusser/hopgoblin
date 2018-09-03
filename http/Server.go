// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "bufio"
  "crypto/tls"
  "errors"
  "fmt"
  "net"
  "regexp"
  "strings"
  "time"
  "github.com/rinusser/hopgoblin/log"
  "github.com/rinusser/hopgoblin/utils"
)


/*
  HTTP Server type: will listen for incoming connections, acting like a proxy server.
 */
type Server struct {
  listener net.Listener //will be set to low-level socket listener
  siteHandlers []SiteHandler //list of site handlers; register with AddSiteHandler()
  Shutdown chan bool    //used to shut down the server instance during tests
  *ProxySettings        //upstream proxy settings
  tlsconfig tls.Config
}

/*
  Create a new server instance with defaults.
 */
func NewServer() *Server {
  rv:=&Server {
    listener: nil,
    Shutdown: make(chan bool),
    ProxySettings: GetDefaultProxySettings(),
  }

  resdir:=utils.GetResourceDir("certs")
  certfile:=utils.GetConfigValue("server.default_certificate_file")
  keyfile:=utils.GetConfigValue("server.default_certificate_key")
  if certfile=="" || keyfile=="" { panic("need to set default TLS certificate and key files in application config") }
  default_cert,err:=tls.LoadX509KeyPair(resdir+certfile,resdir+keyfile)
  if err!=nil {
    panic(fmt.Sprintf("error loading default keypair: %s",err))
  }

  rv.tlsconfig=tls.Config {
    Certificates: []tls.Certificate{default_cert},
    ClientAuth: tls.VerifyClientCertIfGiven,
    NameToCertificate: map[string]*tls.Certificate{},
    ServerName: "hopgoblin.localhost",
  }

  return rv
}


/*
  Registers a site handler.
 */
func (this *Server) AddSiteHandler(h SiteHandler) {
  this.siteHandlers=append(this.siteHandlers,h)

  for host,cert:=range h.GetCertificateMap() {
    this.tlsconfig.Certificates=append(this.tlsconfig.Certificates,*cert)
    this.tlsconfig.NameToCertificate[host]=cert
  }
}

/*
  Adds all registered site handlers.
 */
func (this *Server) AddAllRegisteredSiteHandlers() {
  for _,h:=range GetRegisteredSiteHandlers() {
    this.AddSiteHandler(h)
  }
}


/*
  Starts listening to incoming connections on the given port.

  This method won't return until a boolean "true" is received sent over the Server.Shutdown channel.
 */
func (server *Server) Listen(port uint16) error {
  var err error
  addr:=net.TCPAddr{net.IPv4(127,0,0,1),int(port),""}
  listener,err:=net.ListenTCP("tcp",&addr)
  server.listener=listener
  if err!=nil {
    log.Fatal("unable to listen on port %d",port)
    return err
  }
  log.Debug("listening on %s.\n",listener.Addr().String())
  for {
    listener.SetDeadline(time.Now().Add(1e9))
    log.Trace("waiting for connection...")
    conn,err:=listener.AcceptTCP()
    if err,ok:=err.(net.Error);ok&&err.Timeout() {
      select {
        case <-server.Shutdown:
          log.Debug("received shutdown signal")
          return nil
        default:
          continue
      }
    } else if err!=nil {
      log.Error("failed to accept connection",err)
      continue
    }
    log.Trace("got connection, spawning handler")
    go server.handleConnection(conn)
  }
  return nil
}

func (server *Server) handleConnection(conn net.Conn) {
  log.Trace("handler spawned, waiting for data...")
  defer conn.Close()

  buf:=bufio.NewReadWriter(bufio.NewReader(conn),bufio.NewWriter(conn))
  request,err:=server.readRequest(buf)
  if err!=nil || request==nil {
    return
  }

  host:=""
  if request.Method=="CONNECT" {
    url_parts:=strings.Split(request.Url,":")
    host=url_parts[0]
  } else {
    regex:=regexp.MustCompile(`^([a-z]+://)?([^:/]+)?(:[0-9]+)?/`)
    matches:=regex.FindAllStringSubmatch(request.Url,-1)
    host=matches[0][2]
  }

  var handler *SiteHandler=nil
  for _,h:=range server.siteHandlers {
    if h.HandlesHost(host) {
      handler=&h
      break
    }
  }

  response:=NewResponse()
  if handler==nil {
    if host!="detectportal.firefox.com" {
      log.Debug("denied %s to %s",request.Method,request.Url)
    }
    response.Status=403
    response.Body=[]byte("go away")
    server.WriteAndFlush(buf,response.ToString())
    return
  }

  log.Debug("allowing %s to %s",request.Method,request.Url)

  if request.Method=="CONNECT" {
    response.Status=200
    server.WriteAndFlush(buf,response.ToString())
    host:=request.Url
    colon_pos:=strings.Index(host,":")
    if colon_pos>0 {
      host=host[0:colon_pos]
    }
    buf,request,err=server.startSSLServer(conn,host)
    if err!=nil {
      return
    }
  }

  (*handler).HandleRequest(server,buf,request)
}

/*
  Performs the server-side part of an SSL/TLS handshake on an existing connection.

  This function returns new net.Conn and bufio.ReadWriter instances for the encrypted connection: use only those after a successful
  handshake.
 */
func (this *Server) UpgradeServerConnectionToSSL(conn net.Conn, host string) (net.Conn,*bufio.ReadWriter,error) {
  var tlsconn *tls.Conn
  tlsconfig:=*(&this.tlsconfig)
  tlsconfig.ServerName=host
  tlsconn=tls.Server(conn,&tlsconfig)
  log.Debug("performing TLS handshake...")

  err:=tlsconn.Handshake()
  if err!=nil {
    return nil,nil,err
  }

  buf:=bufio.NewReadWriter(bufio.NewReader(tlsconn),bufio.NewWriter(tlsconn))

  return tlsconn,buf,nil
}

func (server *Server) startSSLServer(conn net.Conn, host string) (*bufio.ReadWriter,*Request,error) {
  _,buf,err:=server.UpgradeServerConnectionToSSL(conn,host)
  if err!=nil {
    log.Debug("TLS handshake failed: %s",err)
    return nil,nil,err
  }

  request,err:=server.readRequest(buf)
  if err!=nil {
    log.Debug("could not read TLS'd request: %s",err)
    return nil,nil,err
  }
  request.Is_ssl=true
  return buf,request,nil
}

func (server *Server) readRequest(buf *bufio.ReadWriter) (*Request,error) {
  log.Trace("reading http request..")
  request_text,err:=ReadHTTPMessageAsString(buf)
  if err!=nil {
    return nil,err
  }

  request:=ParseRequest(request_text)
  if request!=nil {
    return request,nil
  } else {
    return nil,errors.New("could not parse request")
  }
}

/*
  Send data to a connected client.
 */
func (server *Server) WriteAndFlush(buf *bufio.ReadWriter, response string) error { //TODO: why public?
  _,err:=buf.WriteString(response)
  if err!=nil {
    log.Debug("can't write to connection",err)
  }
  err=buf.Flush()
  if err!=nil {
    log.Debug("flush failed",err)
  }
  return err
}

