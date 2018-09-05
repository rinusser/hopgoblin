// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package appconfig

import (
  "github.com/rinusser/hopgoblin/bootstrap"
  "github.com/rinusser/hopgoblin/log"
  "github.com/rinusser/hopgoblin/utils"
)


func init() {
  bootstrap.AfterFlagParse(initHook)
}

func initHook() {
  log.CurrentSettings=log.MergeSettings(ParseApplicationConfiguration(),log.CurrentSettings)
}


/*
  Reads log settings from the application configuration (resources/application.ini).
  All settings are in the INI's "log" section.

  For example:

    [log]
    ; sets the default log level to DEBUG
    default_level = debug

    ; disables log messages for prefix "hopgoblin/http"
    levels.hopgoblin/http = OFF

    ; sets the log timestamp format, see Go's time.Time.Format() documentation
    timestamp_format = 2006-01-02 15:04:05.000
 */
func ParseApplicationConfiguration() log.Settings {
  rv:=log.NewSettings()
  rv.Prefixes["*"]=log.DefaultLevel
  levelstr:=utils.GetConfigValue("log.default_level")
  if levelstr!="" {
    rv.Prefixes["*"]=log.FromString(levelstr)
  }

  for prefix,levelstr:=range utils.GetConfigValuesByPrefix("log.levels.") {
    rv.Prefixes[prefix]=log.FromString(levelstr)
  }

  rv.TimestampFormat=utils.GetConfigValue("log.timestamp_format")

  return rv
}

