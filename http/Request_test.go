// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "testing"
  "github.com/stretchr/testify/assert"
)


/*
  Makes sure ParseRequest() parses all parts of the message correctly.
 */
func TestParseRequest(t *testing.T) {
  input:="GET /asdf HTTP/1.0\r\n"+
         "X-Some-Header: yoyo: 1\r\n"+
         "Accept-Encoding: plain\r\n"+
         "\r\n"+
         "invalid body\rxx\nasdf\r\nfin"

  actual:=ParseRequest(input)
  assert.Equal(t,"GET",actual.Method,"HTTP method")
  assert.Equal(t,"/asdf",actual.Url,"request URL")
  assert.Equal(t,"HTTP/1.0",actual.Protocol,"HTTP version")

  expected_headers:=NewHeaders()
  expected_headers.Set("Accept-Encoding","plain")
  expected_headers.Set("X-Some-Header","yoyo: 1")
  assert.Equal(t,expected_headers,actual.Headers,"headers")
  assert.Equal(t,[]byte("invalid body\rxx\nasdf\r\nfin"),actual.Body,"request body")
}

/*
  Makes sure ParseRequest() returns nil when parsing invalid/incomplete messages.
 */
func TestParseRequestIncomplete(t *testing.T) {
  cases:=[]string {
    "",
    "GET / HTTP/1.1",
  }
  for _,input:=range cases {
    actual:=ParseRequest(input)
    assert.Nil(t,actual,"parsing invalid request should have failed gracefully")
  }
}


/*
  Makes sure .ToString() formats an HTTP request correctly.
 */
func TestRequestToString(t *testing.T) {
  input:=Request {
    Method:"DOALREADY",
    Url:"uri://some/crap",
    message: message {
      Headers:NewHeaders(),
      Body:[]byte("you\nhit\rpay\n\rdirt\t"),
    },
  }
  input.Headers.Set("Let-It-Be","atles")
  input.Headers.Set("Im-A-Header","true")

  actual:=input.ToString()
  expected:="DOALREADY uri://some/crap HTTP/1.1\r\nIm-A-Header: true\r\nLet-It-Be: atles\r\n\r\nyou\nhit\rpay\n\rdirt\t"
  assert.Equal(t,expected,actual,"request body")
}

