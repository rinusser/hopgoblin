// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package log

import (
  "flag"
  "fmt"
  "strings"
  "github.com/rinusser/hopgoblin/bootstrap"
)


var levelArguments LevelArguments


func init() {
  flag.Var(&levelArguments,"log","comma-separated log levels: prefix=level; prefix '*' for global setting")
  bootstrap.AfterFlagParse(initHook)
}

func initHook() {
  CurrentSettings=MergeSettings(CurrentSettings,ParseArguments(levelArguments))
}


/*
  Parses a list of command-line argument strings (without the argument itself) into the internal settings structure.
 */
func ParseArguments(arguments []string) Settings {
  rv:=NewSettings()
  for _,value:=range arguments {
    for _,setting:=range strings.Split(value,",") {
      setting=strings.TrimSpace(setting)
      if setting=="" {
        continue
      }
      var prefix string
      var levelstr string
      if strings.Index(setting,"=")>0 {
        parts:=strings.Split(setting,"=")
        if len(parts)!=2 { panic(fmt.Sprintf("invalid log setting '%s'",setting)) }
        prefix=parts[0]
        levelstr=parts[1]
      } else {
        prefix="*"
        levelstr=setting
      }
      level:=FromString(levelstr)
      rv.Prefixes[prefix]=level
    }
  }
  return rv
}

