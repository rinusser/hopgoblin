// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "fmt"
  "reflect"
  "regexp"
  "sort"
  "strings"
  "github.com/rinusser/hopgoblin/log"
)


/*
  Type for HTTP headers.

  Header name lookups are performed case-insensitively, but the original case will be preserved when generating the string
  representation.
 */
type Headers struct {
  data map[string][]string
}


/*
  Creates an empty header set.
 */
func NewHeaders() *Headers {
  return &Headers{data:make(map[string][]string)}
}

/*
  Parses HTTP headers from raw request data.
  Data needs to contain double newline to mark end of HTTP header block.
 */
func ParseHeaders(data string) *Headers {
  rv:=NewHeaders()
  separator:=strings.Index(data,"\n\n")
  if separator<0 {
    separator=strings.Index(data,"\r\n\r\n")
  }
  if separator<0 {
    log.Debug("could not find end of headers")
    return rv
  }
  request_line_matcher:=regexp.MustCompile(`^[^ ]+ [^ ]+ http/[0-9]\.[0-9]$`)
  for idx,value := range strings.Split(data[0:separator],"\n") {
    if idx==0 && request_line_matcher.MatchString(strings.ToLower(strings.TrimSpace(value))) {
      continue
    }
    colon:=strings.Index(value,":")
    if colon<0 {
      continue
    }
    rv.Set(strings.TrimSpace(value[0:colon]),strings.TrimSpace(value[colon+1:]))
  }
  return rv
}


/*
  Fetches a single header, lookup is case insensitive.
  The boolean return value is set to true if the header was found, false otherwise.
 */
func (this *Headers) Get(key string) (string, bool) {
  value:=""
  parts,found:=this.data[strings.ToLower(key)]
  if found {
    value=parts[1]
  }
  return value,found
}

/*
  Sets a header line. Key case will be preserved when calling ToString() later.
 */
func (this *Headers) Set(key string, value string) {
  this.data[strings.ToLower(key)]=[]string{key,value}
}

/*
  Returns an alphabetically sorted list of keys.
 */
func (this *Headers) Keys() []string {
  keys:=make([]string,0)
  for _,parts:=range this.data {
    keys=append(keys,parts[0])
  }
  sort.Strings(keys)
  return keys
}

/*
  Compares two Headers instances for equality. Instances are considered equal if they contain the same keys and values - both are compared case-sensitively.
 */
func (this *Headers) Equals(that *Headers) bool {
  keys1:=this.Keys()
  keys2:=that.Keys()
  if !reflect.DeepEqual(keys1,keys2) {
    return false
  }
  for _,key:=range keys1 {
    v1,_:=this.Get(key)
    v2,_:=that.Get(key)
    if v1!=v2 {
      return false
    }
  }
  return true
}

/*
  Renders HTTP headers into part of HTTP message ready for transmission in HTTP request/response data.
 */
func (this *Headers) ToString() string {
  var rvs strings.Builder
  for _,key:=range this.Keys() {
    value,_:=this.Get(key)
    rvs.WriteString(fmt.Sprintf("%s: %s\r\n",key,value))
  }
  return rvs.String()
}

