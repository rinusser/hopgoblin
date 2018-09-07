// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "fmt"
)


/*
  Makes sure ParseResponse() parses all parts of the message correctly.
 */
func TestParseResponse(t *testing.T) {
  input:="HTTP/1.0 200 OK\r\n"+
         "X-Some-Header: yoyo: 1\r\n"+
         "Accept-Encoding: plain\r\n"+
         "\r\n"+
         "invalid body\rxx\nasdf\r\nfin"

  actual:=ParseResponse(input)
  assert.Equal(t,"HTTP/1.0",actual.Protocol,"HTTP version")
  assert.Equal(t,200,int(actual.Status),"HTTP status")
  expected_headers:=NewHeaders()
  expected_headers.Set("X-Some-Header","yoyo: 1")
  expected_headers.Set("Accept-Encoding","plain")
  assert.Equal(t,expected_headers,actual.Headers,"headers")
  assert.Equal(t,[]byte("invalid body\rxx\nasdf\r\nfin"),actual.Body,"response body")
}

/*
  Makes sure .ToString() formats the message correctly.
 */
func TestResponseToString(t *testing.T) {
  input:=Response {
    Status:403,
    message: message {
      Protocol: "HTTP/1.2",
      Headers:NewHeaders(),
      Body:[]byte("you\nhit\rpay\n\rdirt\t"),
    },
  }
  input.Headers.Set("Let-It-Be","atles")
  input.Headers.Set("Im-A-Header","true")

  actual:=input.ToString()
  expected:="HTTP/1.2 403 Forbidden\r\nIm-A-Header: true\r\nLet-It-Be: atles\r\n\r\nyou\nhit\rpay\n\rdirt\t"
  assert.Equal(t,expected,actual,"response body")
}

/*
  Makes sure CreateSimpleResponse() can create a response with known status codes and returns nil for unknown codes.
 */
func TestCreateSimpleResponse(t *testing.T) {
  response:=CreateSimpleResponse(200)
  assert.Equal(t,uint16(200),response.Status)

  response=CreateSimpleResponse(403)
  assert.Equal(t,uint16(403),response.Status)

  response=CreateSimpleResponse(999)
  assert.Nil(t,response)
}


/*
  Makes sure .GetPlainTextBodyString() can decode gzip encoding, chunked transfers and any combination thereof.
 */
func TestGetPlainTextBodyString(t *testing.T) {
  gzip_data:=[]byte {
    0x1f, 0x8b, 0x08, 0x00, 0x29, 0x13, 0x7f, 0x5b,
    0x00, 0x03, 0x2b, 0x2e, 0x4d, 0x4e, 0x4e, 0x2d,
    0x2e, 0x06, 0x00, 0xb2, 0xdf, 0x00, 0x6f, 0x07,
    0x00, 0x00, 0x00,
  }
  plain_data:=[]byte("success")

  for chunked:=0;chunked<2;chunked++ {
    for compress:=0;compress<2;compress++ {
      input:=NewResponse()
      data:=plain_data

      if compress>0 {
        data=gzip_data
        input.Headers.Set("Content-Encoding","gzip")
      }

      if chunked>0 {
        input.Headers.Set("Transfer-Encoding","chunked")
        chunk1:=data[0:3]
        chunk2:=data[3:]
        data=append(append(append(append(append(
          //chunk1 length
          []byte(fmt.Sprintf("%x\r\n",len(chunk1))),
          //chunk1 data
          chunk1...),
          //end of chunk1
          0x0d, 0x0a),
          //chunk2 length
          []byte(fmt.Sprintf("%x\r\n",len(chunk2)))...),
          //chunk2 data
          chunk2...),
          //end of chunk2
          0x0d, 0x0a,
          //chunk length 0
          0x30, 0x0d, 0x0a,
          //end of chunk
          0x0d, 0x0a,
        )
      }

      input.Body=data

      actual:=input.GetPlainTextBodyString()
      assert.Equal(t,"success",actual)
    }
  }
}
