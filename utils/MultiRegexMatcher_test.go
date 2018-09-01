// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package utils

import (
  "testing"
  "github.com/stretchr/testify/assert"
)


/*
  Makes sure .MatchesAnyRegex() works correctly.
 */
func TestHostRegexes(t *testing.T) {
  m:=NewMultiRegexMatcher([]string {
    `^www\.asdf\.com$`,
    `^[a-z]+\.example\.com$`,
  })

  expected_successes:=[]string { "www.asdf.com", "www.example.com", "asdf.example.com" }
  expected_failures:=[]string { "www1.asdf.com", "ww.asdf.com", "v1.example.com" }

  for _,value:=range expected_successes {
    assert.True(t,m.MatchesAnyRegex(value),"should have matched: "+value)
  }

  for _,value:=range expected_failures {
    assert.False(t,m.MatchesAnyRegex(value),"shouldn't have matched: "+value)
  }
}

