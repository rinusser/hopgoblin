// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package log

import (
  "fmt"
  "log"
  "os"
  "regexp"
  "runtime"
  "strings"
  "time"
)


/*
  The default log level.
 */
var DefaultLevel=INFO

/*
  The default format string for timestamps in log messages, see time.Time.Format().
 */
var DefaultTimestampFormat="2006-01-02 15:04:05.000"

/*
  The currently active log settings.
 */
var CurrentSettings Settings


var methodNameMangler=regexp.MustCompile(`\(\*([a-zA-Z0-9]+)\)`)


func init() {
  CurrentSettings=NewSettings()
  log.SetFlags(0)
}


/*
  Determines the effective level for a given (qualified) function name.
 */
func getEffectiveLevel(settings Settings, function string) Level {
  for _,prefix:=range settings.getSortedPrefixes() {
    if strings.Index(function,prefix)==0 {
      return settings.Prefixes[prefix]
    }
  }

  value,found:=settings.Prefixes["*"]
  if found {
    return value
  } else {
    return DefaultLevel
  }
}


/*
  Turns given log settings into a command-line argument string.
 */
func AssembleLogSettingsArg(settings Settings) string {
  var rvs []string
  for _,prefix:=range settings.getSortedPrefixes() {
    rvs=append(rvs,fmt.Sprintf("%s=%s",prefix,settings.Prefixes[prefix].String()))
  }
  return "--log="+strings.Join(rvs,",")
}

/*
  Turns the current log settings into a command-line argument string.
 */
func AssemblePassthroughArg() string {
  return AssembleLogSettingsArg(CurrentSettings)
}


func _log(level Level, format string, data...interface{}) {
  pc,_,/*line*/_,_:=runtime.Caller(2)
  function_pretty:=getMethodName(pc)

  if level<getEffectiveLevel(CurrentSettings,function_pretty) {
    return
  }

  msg:=fmt.Sprintf(format,data...)
  log.Printf("%s% 6d %- 5s %s: %s",time.Now().Format(getCurrentTimestampFormat()),os.Getpid(),level.String(),function_pretty,msg)
}

func getCurrentTimestampFormat() string { //XXX could cache this
  rv:=CurrentSettings.TimestampFormat
  if rv=="" {
    rv=DefaultTimestampFormat
  }
  return rv
}

func getMethodName(pc uintptr) string {
  name:=runtime.FuncForPC(pc).Name()
  name=methodNameMangler.ReplaceAllString(name,"${1}")
  if strings.HasPrefix(name,"github.com/rinusser/") {
    name=name[20:]
  }
  return name
}


/*
  Logs a message at TRACE level.
 */
func Trace(format string, data...interface{}) {
  _log(TRACE,format,data...)
}

/*
  Logs a message at DEBUG level.
 */
func Debug(format string, data...interface{}) {
  _log(DEBUG,format,data...)
}

/*
  Logs a message at INFO level.
 */
func Info(format string, data...interface{}) {
  _log(INFO,format,data...)
}

/*
  Logs a message at WARN level.
 */
func Warn(format string, data...interface{}) {
  _log(WARN,format,data...)
}

/*
  Logs a message at ERROR level.
 */
func Error(format string, data...interface{}) {
  _log(ERROR,format,data...)
}

/*
  Logs a message at FATAL level.
 */
func Fatal(format string, data...interface{}) {
  _log(FATAL,format,data...)
}
