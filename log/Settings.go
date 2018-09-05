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
  Prefixes map[string]Level
  TimestampFormat string
}

/*
  Create a new log Settings instance.
 */
func NewSettings() Settings {
  return Settings {
    Prefixes:make(map[string]Level,0),
    TimestampFormat:"",
  }
}

/*
  Returns the registered log prefixes, sorted by descending length.
 */
func (this *Settings) getSortedPrefixes() []string { //XXX: could cache this
  keys:=make([]string,0)
  for key,_:=range this.Prefixes {
    keys=append(keys,key)
  }
  sort.Sort(prefixesType(keys))
  return keys
}


/*
  Takes 2 sets of settings and combines them, with settings in the second parameter (s2) taking precedence.
 */
func MergeSettings(s1 Settings, s2 Settings) Settings {
  rv:=s1

  for key,value:=range s2.Prefixes {
    rv.Prefixes[key]=value
  }

  if s2.TimestampFormat!="" {
    rv.TimestampFormat=s2.TimestampFormat
  }

  return rv
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

