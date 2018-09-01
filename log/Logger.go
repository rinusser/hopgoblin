// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package log

import (
  "flag"
  "fmt"
  "log"
  "os"
  "regexp"
  "runtime"
  "strings"
  "time"
  "github.com/rinusser/hopgoblin/bootstrap"
)


/*
  The default log level.
 */
var DefaultLevel=INFO //TODO: add config setting

/*
  The format string for timestamps in log messages, see time.Time.format().
 */
var TimestampFormat="2006-01-02 15:04:05.000" //TODO: add config setting


var settings Settings
var levelArguments LevelArguments
var methodNameMangler=regexp.MustCompile(`\(\*([a-zA-Z0-9]+)\)`)


func init() {
  flag.Var(&levelArguments,"log","comma-separated log levels: prefix=level; prefix '*' for global setting")
  bootstrap.AfterFlagParse(initHook)
}

func initHook() {
  settings=ParseArguments(levelArguments)

  log.SetFlags(0)
}


/*
  Parses a list of command-line argument strings (without the argument itself) into the internal settings structure.
 */
func ParseArguments(arguments []string) Settings {
  rv:=NewSettings()
  rv.prefixes["*"]=DefaultLevel
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
      rv.prefixes[prefix]=level
    }
  }
  return rv
}

/*
  Determines the effective level for a given (qualified) function name.
 */
func getEffectiveLevel(settings Settings, function string) Level {
  for _,prefix:=range settings.getSortedPrefixes() {
    if strings.Index(function,prefix)==0 {
      return settings.prefixes[prefix]
    }
  }
  value,found:=settings.prefixes["*"]
  if !found { panic("invalid settings, bootstrap incomplete") }
  return value
}


/*
  Turns given log settings into a command-line argument string.
 */
func AssembleLogSettingsArg(settings Settings) string {
  var rvs []string
  for _,prefix:=range settings.getSortedPrefixes() {
    rvs=append(rvs,fmt.Sprintf("%s=%s",prefix,settings.prefixes[prefix].String()))
  }
  return "--log="+strings.Join(rvs,",")
}

/*
  Turns the current log settings into a command-line argument string.
 */
func AssemblePassthroughArg() string {
  return AssembleLogSettingsArg(settings)
}


func _log(level Level, format string, data...interface{}) {
  pc,_,/*line*/_,_:=runtime.Caller(2)
  function_pretty:=getMethodName(pc)

  if level<getEffectiveLevel(settings,function_pretty) {
    return
  }

  msg:=fmt.Sprintf(format,data...)
  log.Printf("%s% 6d %- 5s %s: %s",time.Now().Format(TimestampFormat),os.Getpid(),level.String(),function_pretty,msg)
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
