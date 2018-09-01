// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package utils

import (
  "regexp"
)


/*
  Regex-matching for site handlers. Use this to easily match hostnames against a list of regular expression.
 */
type MultiRegexMatcher struct {
  regexes []string
  compiledRegexes []*regexp.Regexp
}

/*
  Creates a new MultiRegexMatcher instance. Pass any regular expressions that should be matched against.
 */
func NewMultiRegexMatcher(regexes []string) MultiRegexMatcher {
  rv:=MultiRegexMatcher {
    regexes: regexes,
  }
  rv.compileRegexes()
  return rv
}


/*
  Turns the list of regex strings into Regexp instances.
 */
func (this *MultiRegexMatcher) compileRegexes() {
  for _,value:=range this.regexes {
    this.compiledRegexes=append(this.compiledRegexes,regexp.MustCompile(value))
  }
}

/*
  Checks if the given string matches against any of the (compiled) regular expressions.
 */
func (this *MultiRegexMatcher) MatchesAnyRegex(needle string) bool {
  for _,regex:=range this.compiledRegexes {
    if regex.MatchString(needle) {
      return true
    }
  }
  return false
}
