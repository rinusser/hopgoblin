// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package log

import (
  "testing"
  "github.com/stretchr/testify/assert"
)


/*
  Makes sure ParseArguments() can handle both multiple arguments and comma-separated values.
 */
func TestParseArguments(t *testing.T) {
  input:=[]string {
    "*=DEBUG",
    "pkg1=TRACE,pkg2=INFO,pkg3=OFF",
  }
  expected:=map[string]Level {
    "*":DEBUG,
    "pkg1":TRACE,
    "pkg2":INFO,
    "pkg3":OFF,
  }
  assert.Equal(t,expected,ParseArguments(input).Prefixes,"multiple arguments should be combined")
}

/*
  Makes sure ParseArguments() reads a value without prefix as the global setting.
 */
func TestParseArgumentsWithoutKey(t *testing.T) {
  input:=[]string { "WARN,a=OFF" }
  expected:=map[string]Level {
    "*":WARN,
    "a":OFF,
  }
  assert.Equal(t,expected,ParseArguments(input).Prefixes,"setting without prefix should be read as global default")
}

