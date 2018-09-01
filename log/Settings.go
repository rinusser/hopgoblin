// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package log

import (
  "sort"
)


/*
  Type for log settings.
 */
type Settings struct {
  prefixes map[string]Level
}

/*
  Create a new log Settings instance.
 */
func NewSettings() Settings {
  return Settings {
    prefixes:make(map[string]Level,0),
  }
}

/*
  Returns the registered log prefixes, sorted by descending length.
 */
func (this *Settings) getSortedPrefixes() []string { //XXX: could cache this
  keys:=make([]string,0)
  for key,_:=range this.prefixes {
    keys=append(keys,key)
  }
  sort.Sort(prefixesType(keys))
  return keys
}


type prefixesType []string

/*
  required by sort.Interface
 */
func (p prefixesType) Len() int {
  return len(p)
}

/*
  required by sort.Interface
 */
func (p prefixesType) Swap(i, j int) {
  p[i],p[j]=p[j],p[i]
}

/*
  required by sort.Interface
 */
func (p prefixesType) Less(i, j int) bool {
  len1:=len(p[i])
  len2:=len(p[j])
  if len1!=len2 {
    return len1>len2
  }
  return p[i]<p[j]
}

