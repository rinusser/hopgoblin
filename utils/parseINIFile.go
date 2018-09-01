// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package utils

import (
  "bytes"
  "io/ioutil"
  "strings"
  "github.com/rinusser/hopgoblin/log"
)


/*
  Parses an INI file, returns contents as a key=>value map.

  Supports sections. Keys within sections will be prefixed by the section name followed by a dot (.). An empty section name removes
  the prefix from following keys.

  Section names and keys will be converted to lowercase. Section names, keys and values will have any leading and trailing
  whitespaces removed. Lines starting with either a semicolon (;) or pound sign (#) are considered comments. Comments and empty
  lines will be ignored.

  For example, parsing this file:

    Setting 1 = Value 1
    #this is a comment
    ;this too!

    [section1]
    setting2=value2

    [section 2]
      setting 3=

    []
    setting4 = root

  will result in this map:

    map[string]string {
      "Setting 1":"Value 1",
      "section1.setting2":"value2",
      "section 2.setting 3":"",
      "setting4":"root",
    }
 */
func ParseINIFile(filename string) *map[string]string {
  data,err:=ioutil.ReadFile(filename)
  if err!=nil {
    log.Warn("could not read INI file %s",filename)
    return nil
  }
  return parseINIData(data)
}

/*
  Parses INI data. Parsing details see ParseINIFile().
 */
func parseINIData(data []byte) *map[string]string {
  rv:=make(map[string]string)
  lines_data:=bytes.Split(data,[]byte("\n"))

  current_section:=""
  for idx,line_data:=range lines_data {
    line:=strings.TrimSpace(string(line_data))

    if len(line)<1 {
      continue
    }
    first_char:=line[0]

    if first_char==';' || first_char=='#' {
      continue
    }

    if first_char=='[' {
      if line[len(line)-1]!=']' {
        log.Error("invalid section in line %d: \"%s\"",idx+1,line)
        return nil
      }
      current_section=strings.ToLower(strings.TrimSpace(line[1:len(line)-1]))
      if len(current_section)>0 {
        current_section+="."
      }
      continue
    }

    split_pos:=strings.Index(line,"=")
    if split_pos<1 {
      log.Error("invalid setting in line %d: \"%s\"",idx+1,line)
      return nil
    }

    key:=strings.ToLower(strings.TrimSpace(line[0:split_pos]))
    if len(key)<1 {
      log.Error("empty key in line %d: \"%s\"",idx+1,line)
      return nil
    }
    value:=strings.TrimSpace(line[split_pos+1:])

    rv[current_section+key]=value
  }

  return &rv
}
