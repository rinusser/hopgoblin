// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package main

import (
  "flag"
  "fmt"
  "net"
  "strings"
  "github.com/rinusser/hopgoblin/bootstrap"
  "github.com/rinusser/hopgoblin/http"
  "github.com/rinusser/hopgoblin/log"
  _ "github.com/rinusser/hopgoblin/log/appconfig" //keep: allows log configuration in application.ini
  _ "github.com/rinusser/hopgoblin/sitehandlers" //keep this: package's initialization routine registers site handlers
  "github.com/rinusser/hopgoblin/utils"
)


/*
  The default IP address to listen on. Only used as fallback if no other value could be found.
 */
var DefaultListenAddress="127.0.0.1"

/*
  The default TCP port to listen on. Only used as fallback.
 */
var DefaultListenPort="64080"


var iparg   = flag.String("ip","","IP address to listen on")
var portarg = flag.String("port","","TCP port to listen on")


func main() {
  bootstrap.Init()
  server:=http.NewServer()
  server.AddAllRegisteredSiteHandlers()

  addr:=getListeningAddress()
  if addr==nil {
    return
  }

  log.Info("starting server")
  server.Listen(addr)
}


func getListeningAddress() *net.TCPAddr {
  addrstr:=getSettingOrDefault(iparg,"server.listen_address",DefaultListenAddress)
  if !isListeningAddress(addrstr) {
    log.Fatal("%s is not a local interface address",addrstr)
    return nil
  }

  portstr:=getSettingOrDefault(portarg,"server.listen_port",DefaultListenPort)

  addr,err:=net.ResolveTCPAddr("tcp",fmt.Sprintf("%s:%s",addrstr,portstr))
  if err!=nil {
    log.Fatal("%s",err)
    return nil
  }
  if addr==nil {
    log.Fatal("could not resolve %s:%s",addrstr,portstr)
    return nil
  }

  return addr
}

func getSettingOrDefault(arg *string, configkey string, def string) string {
  value:=*arg
  if value=="" {
    value=utils.GetConfigValue(configkey)
  }
  if value=="" {
    value=def
  }
  return value
}

func isListeningAddress(needle string) bool {
  if needle=="0.0.0.0" {
    return true
  }

  haystack,_:=net.InterfaceAddrs()
  needle+="/"

  for _,straw:=range haystack {
    if strings.HasPrefix(straw.String(),needle) {
      return true
    }
  }

  return false
}
