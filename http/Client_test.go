// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "crypto/x509"
  "time"
  "github.com/rinusser/hopgoblin/http/dummyproxy"
)


/*
  Makes sure .CopyProxySettings() makes a snapshot copy of the server's settings
 */
func TestCopyProxySettings(t *testing.T) {
  server:=NewServer()
  server.ProxySettings=nil

  client:=NewClient()
  client.ProxySettings=NewProxySettings("1.2.3.4",1234)

  client.CopyProxySettings(server)
  assert.Nil(t,client.ProxySettings,"nil should have been copied to client")

  server.ProxySettings=NewProxySettings("2.3.4.5",2345)
  client.CopyProxySettings(server)
  assertProxySettings(t,client,"2.3.4.5",2345,"settings should have been copied")

  server.ProxySettings.Host="asdf"
  server.ProxySettings.Port=999
  assertProxySettings(t,client,"2.3.4.5",2345,"settings should have remained unchanged")
}

func assertProxySettings(t *testing.T, client *Client, host string, port int, message string) {
  assert.NotNil(t,client.ProxySettings,message+" (reference)")
  assert.Equal(t,host,client.ProxySettings.Host,message+" (host)")
  assert.Equal(t,port,client.ProxySettings.Port,message+" (port)")
}


/*
  Makes sure certificate hostname mismatches throw errors only if .EnableCertificateVerification is true
 */
func TestCertificateVerificationSetting(t *testing.T) {
  proxyrunner:=dummyproxy.NewDummyProxyRunner()
  proxyrunner.StartRandom()
  time.Sleep(5e8)

  client:=NewClient()
  client.ProxySettings=NewProxySettings("127.0.0.1",proxyrunner.Port)

  request:=Request {
    Method:"GET",
    Url:"/no_encoding/certcheck",
    Is_ssl:true,
    message: message {
      Protocol: "HTTP/1.1",
      Headers: NewHeaders(),
    },
  }

  client.EnableCertificateVerification=true
  request.Headers.Set("Host","invalidhost.localhost")
  response,err:=client.ForwardRequest(request);

  assert.NotNil(t,err)
  err,ok:=err.(x509.HostnameError)
  assert.True(t,ok)
  assert.Nil(t,response)

  proxyrunner.StartRandom()
  time.Sleep(5e8)

  client.ProxySettings.Port=proxyrunner.Port
  client.EnableCertificateVerification=true
  request.Headers.Set("Host","proxied.local")
  response,err=client.ForwardRequest(request);

  assert.Nil(t,err)
  assert.NotNil(t,response)
}
