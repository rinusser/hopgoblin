// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package log

import (
  "strings"
)


/*
  Type alias for log levels.
 */
type Level int
const (
  TRACE Level = iota
  DEBUG
  INFO
  WARN
  ERROR
  FATAL
  OFF
)


/*
  Converts Level to string.
 */
func (l Level) String() string {
  return [...]string{"TRACE","DEBUG","INFO","WARN","ERROR","FATAL","OFF"}[l]
}

/*
  Converts string to Level.
  Input string is case insensitive, will panic if value is invalid.
 */
func FromString(input string) Level {
  switch(strings.ToUpper(strings.TrimSpace(input))) {
    case "TRACE":
      return TRACE
    case "DEBUG":
      return DEBUG
    case "INFO":
      return INFO
    case "WARN":
      return WARN
    case "ERROR":
      return ERROR
    case "FATAL":
      return FATAL
    case "OFF":
      return OFF
  }
  panic("invalid level '"+input+"'")
}
