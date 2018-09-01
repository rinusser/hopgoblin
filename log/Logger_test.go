// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package log

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "io"
  "io/ioutil"
  go_log "log"
  "os"
  "strings"
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
  assert.Equal(t,expected,ParseArguments(input).prefixes)
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
  assert.Equal(t,expected,ParseArguments(input).prefixes,"setting without prefix")
}

/*
  Makes sure log prefixes are matched correctly, and in the correct order (longer first).
 */
func TestGetEffectiveLevel(t *testing.T) {
  settings:=NewSettings()
  settings.prefixes=map[string]Level {
    "*":TRACE,
    "package.type.func1":DEBUG,
    "package.type":INFO,
    "package.type.func2":OFF,
  }
  assert.Equal(t,TRACE,getEffectiveLevel(settings,"otherpackage"),"fallback to default")
  assert.Equal(t,INFO, getEffectiveLevel(settings,"package.type"),"full match")
  assert.Equal(t,INFO, getEffectiveLevel(settings,"package.type.function"),"prefixed match")
  assert.Equal(t,TRACE,getEffectiveLevel(settings,"package.typ"),"cut short")
  assert.Equal(t,DEBUG,getEffectiveLevel(settings,"package.type.func1"),"better match (declared earlier)")
  assert.Equal(t,OFF,  getEffectiveLevel(settings,"package.type.func2"),"better match (declared later)")
}

/*
  Makes sure AssembleLogSettingsArg() sorts settings correctly (more precise first)
 */
func TestAssembleLogSettingsArg(t *testing.T) {
  input:=NewSettings()
  input.prefixes=map[string]Level {
    "zzz/":OFF,
    "*":DEBUG,
    "http/":OFF,
    "aaa/":OFF,
    "http/sometype":ERROR,
    "lll/":WARN,
  }
  assert.Equal(t,"--log=http/sometype=ERROR,http/=OFF,aaa/=OFF,lll/=WARN,zzz/=OFF,*=DEBUG",AssembleLogSettingsArg(input))
}


type testReceiver struct {
}

func (this testReceiver) valueReceiver() {
  Fatal("test")
}

func (this *testReceiver) pointerReceiver() {
  Error("test")
}

func TestFunctionNamePrettiness(t *testing.T) {
  r:=testReceiver{}
  reader,writer:=io.Pipe()
  go_log.SetOutput(writer)
  go func() {
    r.valueReceiver()
    r.pointerReceiver()
    reader.Close()
  }()
  data,_:=ioutil.ReadAll(reader)
  lines:=strings.Split(string(data),"\n")
  assert.Equal(t,3,len(lines),"number of log lines")
  assertMethodName(t,"hopgoblin/log.testReceiver.valueReceiver:",lines[0])
  assertMethodName(t,"hopgoblin/log.testReceiver.pointerReceiver:",lines[1])
  assert.Equal(t,"",lines[2])

  go_log.SetOutput(os.Stderr) //XXX Go's log package currently doesn't support fetching the previous writer, so let's hope it was os.Stderr..
}

func assertMethodName(t *testing.T, expected string, actual string) {
  raw_parts:=strings.Split(actual," ")
  parts:=[]string{}
  for _,part:=range raw_parts {
    if part!="" {
      parts=append(parts,part)
    }
  }
  assert.Equal(t,expected,parts[4])
}
