// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package utils

import (
  "testing"
  "github.com/stretchr/testify/assert"
)


func TestGetConfigValue(t *testing.T) {
  appConfiguration=&map[string]string {
    "aaa.b":  "t",
    "aaa.bb": "Uv",
    "aaa.bbb":"WxY",
  }
  assert.Equal(t,"",   GetConfigValue("aaa."),    "key shouldn't match longer config key")
  assert.Equal(t,"t",  GetConfigValue("aaa.b"),   "exact match (1)")
  assert.Equal(t,"Uv", GetConfigValue("aaa.bb"),  "exact match (2)")
  assert.Equal(t,"WxY",GetConfigValue("aaa.bbb"), "exact match (3)")
  assert.Equal(t,"",   GetConfigValue("aaa.bbbb"),"key shouldn't match shorter config key")

  appConfiguration=nil
  assert.Equal(t,"",GetConfigValue("aaa.bbb"),"nil application config should return empty values gracefully")
}

func TestGetConfigValuesByPrefix(t *testing.T) {
  config:=map[string]string {
    "pkg1":    "h",
    "pkg2":    "ij",
    "pkg2.":   "K",
    "pkg2.1":  "Ll",
    "pkg2.abc":"",
    "pkg2.def":"mNO P",
  }
  appConfiguration=&config

  assert.Equal(t,0,len(GetConfigValuesByPrefix("pkg3")), "unset keys shouldn't produce results")
  assert.Equal(t,0,len(GetConfigValuesByPrefix("pkg11")),"prefix should be matched fully")
  assert.Equal(t,0,len(GetConfigValuesByPrefix("pkg1.")),"trailing dots should still be required to match")
  assert.Equal(t,0,len(GetConfigValuesByPrefix("kg")),   "matching should start at the beginning")

  expected:=map[string]string {
    "":"h",
  }
  assert.Equal(t,expected,GetConfigValuesByPrefix("pkg1"),"should match only setting with same key")

  expected=map[string]string {
    "":   "K",
    "1":  "Ll",
    "abc":"",
    "def":"mNO P",
  }
  assert.Equal(t,expected,GetConfigValuesByPrefix("pKG2."),"should be case sensitive and exclude key with missing dot at end")

  appConfiguration=nil
  assert.Equal(t,0,len(GetConfigValuesByPrefix("pkg2.")),"nil application config should return empty values gracefully")
}
