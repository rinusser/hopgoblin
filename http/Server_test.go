// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "bufio"
  "crypto/tls"
  "fmt"
  "io/ioutil"
  "net"
  go_http "net/http"
  "net/url"
  "os"
  "regexp"
  "strings"
  "time"
  "github.com/rinusser/hopgoblin/http/dummyproxy"
  "github.com/rinusser/hopgoblin/log"
  "github.com/rinusser/hopgoblin/utils"
)


type ServerTestProxySiteHandler struct {
}

func (h ServerTestProxySiteHandler) HandlesHost(host string) bool {
  return host=="proxied.local"
}

func (h ServerTestProxySiteHandler) HandleRequest(server *Server, browserio *bufio.ReadWriter, request *Request) {
  log.Debug("handling request")
  log.Trace("request.is_ssl=%t",request.Is_ssl)
  client:=NewClient()
  client.CopyProxySettings(server)
  response,err:=client.ForwardRequest(*request)
  if err!=nil { panic(err) }
  server.WriteAndFlush(browserio,response.ToString())
}

func (this ServerTestProxySiteHandler) GetCertificateMap() map[string]*tls.Certificate {
  cert:=utils.LoadCertificate(utils.GetResourceDir("certs"),"test")
  return map[string]*tls.Certificate {
    "proxied.local":cert,
  }
}


type ServerTestDirectSiteHandler struct {
}

func (h ServerTestDirectSiteHandler) HandlesHost(host string) bool {
  return host=="direct.local"
}

func (h ServerTestDirectSiteHandler) HandleRequest(server *Server, browserio *bufio.ReadWriter, request *Request) {
  log.Debug("received %s request to %s",request.Method,request.Url)
  log.Trace("request.is_ssl=%t",request.Is_ssl)
  response:=NewResponse()
  response.Status=200
  if strings.Index(request.Url,"/no_encoding")>=0 {
    response.Body=[]byte(request.Url)
  } else if strings.Index(request.Url,"/chunked")>=0 {
    response.Headers.Set("Transfer-Encoding","chunked")
    response.Body=ChunkEncodeBody(request.Url,1,1)
  } else {
    panic("unhandled encoding")
  }
  server.WriteAndFlush(browserio,response.ToString())
}

func (this ServerTestDirectSiteHandler) GetCertificateMap() map[string]*tls.Certificate {
  cert:=utils.LoadCertificate(utils.GetResourceDir("certs"),"test")
  return map[string]*tls.Certificate {
    "direct.local":cert,
  }
}


func runServer(port int) (*Server, *go_http.Client) {
  server:=NewServer()
  server.AddSiteHandler(ServerTestProxySiteHandler{})
  server.AddSiteHandler(ServerTestDirectSiteHandler{})
  go server.Listen(uint16(port))

  proxy_url:=fmt.Sprintf("http://127.0.0.1:%d",port)
  os.Setenv("http_proxy", proxy_url)
  os.Setenv("https_proxy",proxy_url)

  transport:=&go_http.Transport {
    Proxy: func(req *go_http.Request) (*url.URL, error) { return url.Parse(proxy_url) },
    TLSClientConfig: &tls.Config{InsecureSkipVerify:true},
  }
  client:=&go_http.Client{Transport: transport}

  time.Sleep(5e8)
  return server,client
}


/*
  Makes sure HTTP requests to unhandled domains are denied with HTTP 403.
 */
func TestServerHTTP403(t *testing.T) {
  server,client:=runServer(64080)
  defer func() { server.Shutdown<-true }()

  response,err:=client.Get("http://does.not.exist/asdf")
  assert.Nil(t,err,"http request should have worked")
  if err!=nil {
    return
  }
  assert.Equal(t,403,response.StatusCode)
}

/*
  Makes sure HTTPS requests to unhandled domains are denied with HTTP 403.
 */
func TestServerHTTPS403(t *testing.T) {
  server,client:=runServer(64081)
  defer func() { server.Shutdown<-true }()

  _,err:=client.Get("https://does.not.exist/asdf")
  assert.Equal(t,"Get https://does.not.exist/asdf: Forbidden",err.Error(),"error response")
}


/*
  Makes sure HTTP(S) requests handled by the site handler directly are working correctly.
 */
func TestServerDirect(t *testing.T) {
  cases:=[][]string {
    //request url                             expected response body
    {"http://direct.local/no_encoding/http",  "http://direct.local/no_encoding/http"},
    {"https://direct.local/no_encoding/https","/no_encoding/https"},
    {"http://direct.local/chunked/http",      "http://direct.local/chunked/http"},
    {"https://direct.local/chunked/https",    "/chunked/https"},
  }
  for key,c:=range cases {
    runServerDirectTest(t,64090+key,c[0],c[1])
  }
}


func runServerDirectTest(t *testing.T, port int, url string, expectation string) {
  server,client:=runServer(port)
  defer func() { server.Shutdown<-true }()
  if !server.SupportsEncryption && url[0:8]=="https://" {
    log.Warn("skipping test case: encryption not supported")
    return
  }
  runServerDirectAssertions(t,client,port,url,expectation)
}

func runServerDirectAssertions(t *testing.T, client *go_http.Client, port int, url string, expectation string) {
  response,err:=client.Get(url)
  assert.Nil(t,err,"http request should have succeeded")
  if err==nil { return }
  defer response.Body.Close()
  assert.Equal(t,"200 OK",response.Status,"HTTP status")
  body,err:=ioutil.ReadAll(response.Body)
  assert.Nil(t,err,"reading body failed")
  assert.Equal(t,expectation,string(body),"HTTP response body")
}


/*
  Makes sure HTTP requests forwarded by the site handler are working correctly.
 */
func TestServerHTTPProxy(t *testing.T) {
  expected:=[]string{"POST http://proxied.local/asdf HTTP/1.1"}
  expected_body:="some/body"
  runServerProxyTest(t,64085,"http://proxied.local/asdf",expected,&expected_body)
}

/*
  Makes sure HTTPS requests forwarded by the site handler are working correctly.
 */
func TestServerHTTPSProxy(t *testing.T) {
  expected:=[]string{
    "CONNECT proxied.local:443 HTTP/1.1",
    "POST /asdf HTTP/1.1",
  }
  expected_body:="some/body"
  runServerProxyTest(t,64087,"https://proxied.local/asdf",expected,&expected_body)
}

/*
  Makes sure HTTP requests with a chunked response are working correctly.
 */
func TestServerHTTPProxyChunked(t *testing.T) {
  expected:=[]string{"POST http://proxied.local/chunked/fdsa HTTP/1.1"}
  expected_body:="some/body"
  runServerProxyTest(t,64088,"http://proxied.local/chunked/fdsa",expected,&expected_body)
}

/*
  Makes sure HTTPS requests with a chunked response are working correctly.
 */
func TestServerHTTPSProxyChunked(t *testing.T) {
  expected:=[]string{
    "CONNECT proxied.local:443 HTTP/1.1",
    "POST /chunked/fdsa HTTP/1.1",
  }
  expected_body:="some/body"
  runServerProxyTest(t,64089,"https://proxied.local/chunked/fdsa",expected,&expected_body)
}

func runServerProxyTest(t *testing.T, port int, url string, expected []string, expected_body *string) {
  server,client:=runServer(port)
  defer func() { server.Shutdown<-true }()
  if !server.SupportsEncryption && url[0:8]=="https://" {
    log.Warn("skipping test case: encryption not supported")
    return
  }

  proxyrunner:=dummyproxy.NewDummyProxyRunner()
  err:=proxyrunner.Start()
  if err!=nil { panic(err) }

  server.ProxySettings=NewProxySettings("127.0.0.1",64086)

  time.Sleep(5e8)

  log.Trace("starting http POST..")
  postbody:="some/body"
  result,err:=client.Post(url,"text/plain",strings.NewReader(postbody))
  defer result.Body.Close()
  log.Trace("finished http POST")
  assert.Nil(t,err,"http.Post() failed")

  proxyresult,proxyerrout,err:=proxyrunner.ReadAndWait()
  if len(proxyerrout)>0 {
    fmt.Fprintln(os.Stderr,string(proxyerrout))
  }

  if err!=nil {
    log.Error("%v",err)
    panic(err)
  }
  assert.Equal(t,418,result.StatusCode)

  request_lines:=findRequestLines(strings.TrimSpace(string(proxyresult)))
  assert.Equal(t,expected,request_lines)

  if expected_body!=nil {
    last_body:=findLastBody(string(proxyresult))
    assert.Equal(t,*expected_body,last_body)
  }
}

func findRequestLines(input string) []string {
  rv:=[]string{strings.TrimSpace(strings.Split(input,"\n")[0])}
  regex:=regexp.MustCompile(`(GET|POST|PUT|DELETE|HEAD|TRACE|OPTIONS|CONNECT) +[^ ]+ +HTTP/[0-9]\.[0-9]`)
  matches:=regex.FindAllStringSubmatch(input[14:],-1)
  for _,match:=range matches {
    rv=append(rv,match[0])
  }
  return rv
}

func findLastBody(input string) string {
  return input[strings.LastIndex(input,"\r\n\r\n")+4:]
}


/*
  Makes sure aborted connections are handled gracefully.

  More cases could be added, supposing it's possible to even trigger them with Go. A separate dummyclient could be created again,
  maybe that'd help reaching more server error states.

  TODO: Closing the client connection before the upstream response is received should trigger an error in the server but doesn't.
        This can definitely be triggered with a browser - why not here?
 */
func TestConnectionAborts(t *testing.T) {
  runAbortTest(t,0)
  runAbortTest(t,1)
  runAbortTest(t,99)
}

func runAbortTest(t *testing.T, num int) {
  port:=64100+num
  server,client:=runServer(port)
  defer func() { server.Shutdown<-true }()
  if !server.SupportsEncryption {
    log.Warn("skipping test case: encryption not supported")
    return
  }

  proxyrunner:=dummyproxy.NewDummyProxyRunner()
  err:=proxyrunner.Start()
  if err!=nil { panic(err) }
  server.ProxySettings=NewProxySettings("127.0.0.1",64086)

  openAbortedConnection(t,port,num)
  runServerDirectAssertions(t,client,port,"http://direct.local/no_encoding/http","http://direct.local/no_encoding/http")
}

func openAbortedConnection(t *testing.T, port int, num int) {
  step:=0

  log.Trace("running step %d...",step)
  conn,err:=net.Dial("tcp",fmt.Sprintf("127.0.0.1:%d",port))
  assert.Nil(t,err)
  defer conn.Close()
  conn.Write([]byte("CONNECT proxied.local:443 HTTP/1.1\r\n\r\n"))
  if step++;num<step { return }

  time.Sleep(5e8)
  log.Trace("running step %d...",step)
  raw_buf:=bufio.NewReadWriter(bufio.NewReader(conn),bufio.NewWriter(conn))
  ReadHTTPMessageAsString(raw_buf)
  tlsconfig:=&tls.Config {
    InsecureSkipVerify:true,
    ServerName:"direct.local",
  }
  tlsconn:=tls.Client(conn,tlsconfig)
  tlsconn.Write([]byte("GET /no_encoding/delayed HTTP/1.1"))
  if step++;num<step { return }

  time.Sleep(5e8)
  log.Trace("running step %d...",step)
  tlsconn.Write([]byte("\r\nHost: proxied.local\r\n\r\n"))
  log.Trace("force closing connection...");
  tlsconn.Close()
  conn.Close()
  time.Sleep(5e9)
}
