// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "strings"
)

/*
  Represents common parts of HTTP requests and responses.
 */
type message struct {
  firstLineParts []string  //first line of request/response split by whitespaces, e.g. {"GET","/","HTTP/1.1"}
  Protocol string          //e.g. "HTTP/1.1"
  Headers *Headers         //e.g. {"Content-Type":"application/json"}
  Body []byte              //e.g. {0x31,0x32,0x33}
}


/*
  Reads parts common to HTTP requests/responses from string.
  Returns nil on error.
 */
func ParseMessage(input string) *message {
  first_newline_pos:=strings.Index(input,"\n")
  if first_newline_pos<0 {
    return nil
  }
  first_line:=strings.TrimSpace(input[0:first_newline_pos])
  first_parts:=strings.Split(first_line," ")

  headers:=ParseHeaders(input)

  body:=[]byte(input[strings.Index(input,"\r\n\r\n")+4:])

  return &message{
    firstLineParts:first_parts,
    Headers:headers,
    Body:body,
  }
}
