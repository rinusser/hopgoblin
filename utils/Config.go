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


/*
  Fetches key/value pairs from the application configuration.

  Keys are returned without the prefix, for example reading the configuration

    [log]
    levels.pkg1 = trace
    levels.pkg2 = off

  with prefix "log.levels." will result in this:

    map[string]string {
      "pkg1":"trace",
      "pkg2":"off",
    }
   */
func GetConfigValuesByPrefix(prefix string) map[string]string {
  rv:=make(map[string]string)
  if appConfiguration==nil {
    return rv
  }

  prefix_length:=len(prefix)
  prefix=strings.ToLower(prefix)
  for key,value:=range *appConfiguration {
    if strings.HasPrefix(key,prefix) {
      rv[key[prefix_length:]]=value
    }
  }

  return rv
}
