// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "fmt"
  "strings"
  "github.com/rinusser/hopgoblin/log"
)


/*
  Represents a single HTTP request.
 */
type Request struct {
  Method string              //e.g. "PUT"
  Url string                 //e.g. "/api/items/new"
  Is_ssl bool                //e.g. true
  message
}

/*
  Parses a string into an Request instance.
 */
func ParseRequest(input string) *Request {
  var rv Request
  message:=ParseMessage(input)
  if message==nil || len(message.firstLineParts)!=3 {
    log.Debug("could not parse HTTP request")
    return nil
  }

  rv.Method,rv.Url,rv.Protocol=message.firstLineParts[0],message.firstLineParts[1],message.firstLineParts[2]
  rv.Headers=message.Headers
  rv.Body=message.Body
  rv.Is_ssl=false //XXX should this check for "https" in URL?

  return &rv
}

/*
  Turns an Request instance into a string, ready for transmission to a server.
 */
func (request *Request) ToString() string {
  var rvs strings.Builder
  rvs.WriteString(fmt.Sprintf("%s %s HTTP/1.1\r\n",request.Method,request.Url))
  rvs.WriteString(request.Headers.ToString())
  rvs.WriteString("\r\n")
  rvs.Write(request.Body)
  return rvs.String()
}
