// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "testing"
  "github.com/stretchr/testify/assert"
)


/*
  Makes sure .Get() fetches header names case-insensitively.
 */
func TestHeadersCaseInsensitivity(t *testing.T) {
  h:=NewHeaders()
  value1:="\tcAse\x00spaCe \n"
  value2:="other"
  h.Set("User-Agent",value1)
  h.Set("X-Something",value2)
  assertFoundAndEqual(t,h,"user-agent",value1)
  assertFoundAndEqual(t,h,"User-Agent",value1)
  assertFoundAndEqual(t,h,"usEr-agenT",value1)
  assertFoundAndEqual(t,h,"x-someTHING",value2)
}

func assertFoundAndEqual(t *testing.T, h *Headers, key string, expected string) {
  actual,found:=h.Get(key)
  assert.True(t,found)
  assert.Equal(t,expected,actual)
}


/*
  Makes sure .Keys() returns header keys in alphabetical order.
 */
func TestHeadersKeys(t *testing.T) {
  h:=NewHeaders()
  h.Set("User-Agent","asdf")
  h.Set("Host","fdsa")
  h.Set("X-Y","111")
  assert.Equal(t,[]string{"Host","User-Agent","X-Y"},h.Keys())
}


type headerAdderFunc func(h *Headers)

func headerAdder(key string, value string) headerAdderFunc {
  return func(h *Headers) {
    h.Set(key,value)
  }
}

/*
  Makes sure .Equals() works properly.
 */
func TestHeadersEquals(t *testing.T) {
  h1:=NewHeaders()
  h2:=NewHeaders()
  add1:=headerAdder("User-Agent","hopGoblin")
  add2:=headerAdder("Content-Type","application/json")
  add3:=headerAdder("Transfer-Encoding","chunked")

  assertEquality(t,h1,h2,"empty lists should be equal")
  add1(h1)
  assertInequality(t,h1,h2,"one list should be bigger")
  add1(h2)
  assertEquality(t,h1,h2,"two lists with same entry should be equal")

  add2(h1)
  add3(h2)
  assertInequality(t,h1,h2,"two lists with different entries shouldn't be equal")

  assert.Equal(t,2,len(h2.Keys()),"self-check, just making sure adder function is working")
  add3(h1)
  assert.Equal(t,3,len(h1.Keys()),"self-check, just making sure adder function is working")
  add2(h2)
  assertEquality(t,h1,h2,"entry order shouldn't matter")
}

func assertEquality(t *testing.T, h1 *Headers, h2 *Headers, message string) {
  assert.True(t,h1.Equals(h1),message)
  assert.True(t,h1.Equals(h2),message)
  assert.True(t,h2.Equals(h1),message)
  assert.True(t,h2.Equals(h2),message)
}

func assertInequality(t *testing.T, h1 *Headers, h2 *Headers, message string) {
  assert.False(t,h1.Equals(h2),message)
  assert.False(t,h2.Equals(h1),message)
}


/*
  Makes sure .ToString() works.
 */
func TestHeadersToString(t *testing.T) {
  h:=NewHeaders()
  h.Set("Content-Length","123")
  h.Set("DNT","1")
  expected:="Content-Length: 123\r\nDNT: 1\r\n"
  assert.Equal(t,expected,h.ToString())
}


/*
  Makes sure ParseHeaders() trims keys/values and preserves their cases
 */
func TestParseHeadersBasics(t *testing.T) {
  expected:=NewHeaders()
  expected.Set("hdr1","val1")
  expected.Set("headER 2","VaLue2")

  runParseHeadersTest(t,
    "hdr1:val1\nheadER 2\t :\t VaLue2\t \n\ndata:asdf\n",
    expected,
    "basic functionality")
}

/*
  Makes sure ParseHeaders() correctly handles CRLF newlines
 */
func TestParseHeadersCRLF(t *testing.T) {
  expected:=NewHeaders()
  expected.Set("Host","example.com")

  runParseHeadersTest(t,
    "Host: example.com\r\n\r\nData: yes\r\n",
    expected,
    "CRLF newlines should work")
}

/*
  Makes sure ParseHeaders() ignores the first line in HTTP requests/responses.
 */
func TestParseHeadersIgnoresHTTPRequestLine(t *testing.T) {
  expected:=NewHeaders()
  expected.Set("Host","example.com")

  runParseHeadersTest(t,
    "GET http://www.asdf.com/ HTTP/2.0\r\nHost: example.com\r\n\r\nData: yes\r\n",
    expected,
    "HTTP request line should be ignored")

    runParseHeadersTest(t,
    "HTTP/1.1 200 OK\r\nHost: example.com\r\n\r\nData: yes\r\n",
    expected,
    "HTTP response line should be ignored")
}

func runParseHeadersTest(t *testing.T, input string, expected *Headers, description string) {
  actual:=ParseHeaders(input)
  assert.True(t,expected.Equals(actual),description)
}

