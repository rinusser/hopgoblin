// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package main

import (
  "github.com/rinusser/hopgoblin/bootstrap"
  "github.com/rinusser/hopgoblin/http"
  "github.com/rinusser/hopgoblin/log"
  _ "github.com/rinusser/hopgoblin/log/appconfig" //keep: allows log configuration in application.ini
  _ "github.com/rinusser/hopgoblin/sitehandlers" //keep this: package's initialization routine registers site handlers
)

func main() {
  bootstrap.Init()
  log.Info("starting server")
  server:=http.NewServer()
  server.AddAllRegisteredSiteHandlers()
  server.Listen(64080)
}
