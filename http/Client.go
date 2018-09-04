// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "bufio"
  "crypto/tls"
  "crypto/x509"
  "fmt"
  "io/ioutil"
  "net"
  "os"
  "strings"
  "github.com/rinusser/hopgoblin/log"
  "github.com/rinusser/hopgoblin/utils"
)


/*
  HTTP client for outgoing connections.
 */
type Client struct {
  conn net.Conn                      //the outgoing connection
  *ProxySettings                     //proxy settings to use
  EnableCertificateVerification bool //whether remote certificates should be verified
}

/*
  Creates a default Client instance.
 */
func NewClient() *Client {
  return &Client {
    ProxySettings:GetDefaultProxySettings(),
    EnableCertificateVerification:true,
  }
}


var caCertPool *x509.CertPool=nil


/*
  Take proxy server settings from parent server instance.
 */
func (this *Client) CopyProxySettings(server *Server) {
  settings:=server.ProxySettings
  if settings!=nil {
    val:=*settings
    settings=&val
  }
  this.ProxySettings=settings
}

func sendHTTPStringAndParseResponse(request string, buf *bufio.ReadWriter) (*Response,error) {
  _,err:=buf.WriteString(request)
  if err!=nil {
    log.Error("ERROR: could not send request (%s)",err)
    return nil,err
  }

  err=buf.Flush()
  if err!=nil {
    log.Error("flush failed: %s",err)
    return nil,err
  }

  log.Trace("starting to read response...")
  response_text,err:=ReadHTTPMessageAsString(buf)
  log.Trace("finished reading")
  if err!=nil {
    return nil,err
  }

  response:=ParseResponse(response_text)
  return &response,nil
}

/*
  Forwards an HTTP request to the upstream proxy.
  Will use the HTTP CONNECT method to open a tunnel for SSL requests
 */
func (client *Client) ForwardRequest(request Request) (*Response,error) {
  log.Debug("connecting to proxy %s:%d\n",client.ProxySettings.Host,client.ProxySettings.Port)
  conn,err:=net.Dial("tcp",fmt.Sprintf("%s:%d",client.ProxySettings.Host,client.ProxySettings.Port))
  if err!=nil {
    log.Warn("could not connect to proxy %s:%d (%s)",client.ProxySettings.Host,client.ProxySettings.Port,err)
    return CreateSimpleResponse(502),nil
  }
  log.Trace("got connection to proxy")
  defer conn.Close()
  client.conn=conn
  request.Headers.Set("Connection","close")
  buf:=bufio.NewReadWriter(bufio.NewReader(conn),bufio.NewWriter(conn))

  if request.Is_ssl {
    log.Trace("handling https request, establishing tunnel through proxy..")
    host,found:=request.Headers.Get("Host")
    if !found {
      log.Error("no host header found in request, aborting")
      return nil,nil
    }

    response,err:=sendHTTPStringAndParseResponse("CONNECT "+host+":443 HTTP/1.1\r\n\r\n",buf)
    if err!=nil {
      log.Error("could not communicate with proxy: %s",err)
      return nil,nil
    }
    if response.Status!=200 {
      log.Warn("got status %d from proxy",response.Status)
      return response,nil //TODO: should this be a new, generic 503 maybe?
    }

    tlsconfig:=&tls.Config{
      InsecureSkipVerify:!client.EnableCertificateVerification,
      ServerName:host,
      RootCAs:GetCertificatePool(),
    }
    tlsconn:=tls.Client(conn,tlsconfig)
    log.Trace("performing TLS handshake...")
    err=tlsconn.Handshake()
    if err!=nil {
      log.Error("TLS handshake error: %v",err)
      return nil,err
    }
    buf=bufio.NewReadWriter(bufio.NewReader(tlsconn),bufio.NewWriter(tlsconn))
  }

  response,err:=sendHTTPStringAndParseResponse(request.ToString(),buf)
  return response,nil
}

/*
  Gets the CA certificate pool.
  The pool consists of CAs supplied by the system, additionally CAs loaded from the "certs" resource directory.
 */
func GetCertificatePool() *x509.CertPool {
  if caCertPool==nil {
    caCertPool,_=x509.SystemCertPool()
    if caCertPool==nil {
      caCertPool=x509.NewCertPool()
    }

    certspath:=utils.GetResourceDir("certs")
    certsdir,err:=os.Open(certspath)
    if err!=nil { panic(err) }
    defer certsdir.Close()
    certnames,err:=certsdir.Readdirnames(0)
    if err!=nil { panic(err) }
    for _,certname:=range certnames {
      if strings.Index(certname,"CA-")!=0 {
        continue
      }
      certs,err:=ioutil.ReadFile(certspath+certname)
      if err!=nil || !caCertPool.AppendCertsFromPEM(certs) {
        log.Warn("could not load CA certificate %s",certname)
      } else {
        log.Debug("loaded CA certificate %s",certname)
      }
    }
  } else {
    log.Trace("reusing cached CA certificate pool")
  }

  return caCertPool
}
