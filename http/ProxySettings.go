// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "fmt"
  "github.com/rinusser/hopgoblin/utils"
)


/*
  Proxy settings for HTTP clients.
 */
type ProxySettings struct {
  Host string
  Port int
}


/*
  Creates a new ProxySettings instance.
 */
func NewProxySettings(host string, port int) *ProxySettings {
  if port<1 || port>65535 {
    panic("invalid TCP port: "+string(port))
  }
  return &ProxySettings {
    Host: host,
    Port: port,
  }
}

/*
  Fetches the default upstream proxy settings.
 */
func GetDefaultProxySettings() *ProxySettings {
  host:=utils.GetConfigValue("proxy.host")
  portstr:=utils.GetConfigValue("proxy.port")
  if len(host)<1 || len(portstr)<1 {
    panic("invalid proxy settings: host="+host+", port="+portstr)
  }
  port:=0
  fmt.Sscanf(portstr,"%d",&port)
  if port<1 || port>65535 {
    panic("invalid proxy port \""+portstr+"%s\"")
  }
  return NewProxySettings(host,port)
}
