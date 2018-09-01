// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "testing"
  "github.com/stretchr/testify/assert"
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
