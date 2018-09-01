// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package log

import (
  "strings"
)


/*
  Custom datatype for flag package.
 */
type LevelArguments []string


/*
  Formats argument list into single comma-separated string.
 */
func (this *LevelArguments) String() string {
  if this!=nil {
    return strings.Join(*this,",")
  }
  return ""
}

/*
  Unused, but required for flag interface.
 */
func (this *LevelArguments) Value() []string {
  return *this
}

/*
  Called by flag package.
  Adds passed settings to the existing list of values.
 */
func (this *LevelArguments) Set(value string) error {
  for _,part:=range strings.Split(value,",") {
    *this=append(*this,part)
  }
  return nil
}
