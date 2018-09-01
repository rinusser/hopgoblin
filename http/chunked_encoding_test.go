// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "testing"
  "github.com/stretchr/testify/assert"
)


/*
  Makes sure ChunkEncodeBody() works.
 */
func TestChunkEncodeBody(t *testing.T) {
  input:="this is a text with length 29"
  actual:=ChunkEncodeBody(input,900,1)
  expected:=[]byte("1d\r\nthis is a text with length 29\r\n0\r\n\r\n")
  assert.Equal(t,expected,actual,"chunk length should be in hex")
}

/*
  Makes sure ChunkDecodeBody() works.
 */
func TestChunkDecodeBody(t *testing.T) {
  cases:=[][]string {
    {"10\r\nABCDEFGHIJKLMNOP\r\n5\r\nQRSTU\r\n0\r\n\r\n","ABCDEFGHIJKLMNOPQRSTU"},
  }
  for _,c:=range cases {
    input:=[]byte(c[0])
    actual:=string(ChunkDecodeBody(input))
    assert.Equal(t,c[1],actual)
  }
}
