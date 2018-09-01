// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package utils

import (
  "strings"
  "github.com/rinusser/hopgoblin/bootstrap"
)


var appConfiguration *map[string]string=nil

func init() {
  bootstrap.AfterFlagParse(initHook)
}

func initHook() {
  appConfiguration=ParseINIFile(GetResourceBaseDir()+"application.ini")
}


/*
  Fetches a value from the application configuration.
  By default the configuration is read from resources/application.ini (relative to application root).
 */
func GetConfigValue(key string) string {
  if appConfiguration==nil {
    return ""
  }
  value,_:=(*appConfiguration)[strings.ToLower(key)]
  return value
}
