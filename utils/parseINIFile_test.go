// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package utils

import (
  "testing"
  "github.com/stretchr/testify/assert"
)


func TestParseINIData(t *testing.T) {
  input:=`asdF=1
    #ignored
    ;ignored too = 3

    []
    not ignored = 4
    [seCtion 1]
    K=v
    []

    1=
    `
  expected:=map[string]string {
    "asdf":"1",
    "not ignored":"4",
    "section 1.k":"v",
    "1":"",
  }

  assert.Equal(t,expected,*parseINIData([]byte(input)))
}

