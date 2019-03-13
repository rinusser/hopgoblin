// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "bytes"
  "compress/gzip"
  "fmt"
  "io/ioutil"
  "strings"
  "github.com/rinusser/hopgoblin/log"
)


/*
  Represents an individual HTTP response.
 */
type Response struct {
  Status uint16  //e.g. 301
  message
}

var statusMessages = map[uint16]string {
  200:"OK",
  206:"Partial Content",
  301:"Moved Permanently",
  302:"Found",
  303:"See Other",
  304:"Not Modified",
  307:"Temporary Redirect",
  308:"Permanent Redirect",
  400:"Bad Request",
  401:"Unauthorized",
  403:"Forbidden",
  404:"Not Found",
  500:"Internal Server Error",
  502:"Bad Gateway",
  503:"Service Unavailable",
  504:"Gateway Timeout",
}

/*
  Creates a new HTTP instance, defaulting to status 200 (OK).
 */
func NewResponse() *Response {
  return &Response {
    Status: 200,
    message: message {
      Protocol: "HTTP/1.1",
      Headers: NewHeaders(),
    },
  }
}

/*
  Creates a new HTTP instance with the given status code and a rudimentary content body.
 */
func CreateSimpleResponse(code uint16) *Response {
  status_text,found:=statusMessages[code]
  if !found {
    log.Error("response code %d unknown",code)
    return nil
  }

  rv:=NewResponse()
  rv.Status=uint16(code)
  rv.Headers=NewHeaders()
  rv.Headers.Set("Content-Type","text/plain")
  rv.Body=[]byte(status_text)
  return rv
}

/*
  Parses a string into a Response instance.
 */
func ParseResponse(input string) Response {
  var rv Response
  message:=ParseMessage(input)  //TODO: null check

  fmt.Sscanf(message.firstLineParts[1],"%d",&rv.Status)
  rv.Protocol=message.firstLineParts[0]
  rv.Headers=message.Headers
  rv.Body=message.Body

  return rv
}

/*
  Generates a string representation of a Response instance into a string ready for transmission.
 */
func (response *Response) ToString() string {
  var rvs strings.Builder
  status_text,found:=statusMessages[response.Status]
  if !found {
    status_text="UNHANDLED"
  }
  rvs.WriteString(fmt.Sprintf("%s %03d %s\r\n",response.Protocol,response.Status,status_text))
  rvs.WriteString(response.Headers.ToString())
  rvs.WriteString("\r\n")
  rvs.Write(response.Body)
  return rvs.String()
}

/*
  Gets the response body, with any transfer encoding and compression stripped off.
 */
func (this *Response) GetPlainTextBodyString() string {
  body:=this.Body

  transfer_decoders:=map[string]bodyDecoder {
    "chunked":ChunkDecodeBody,
  }

  content_decoders:=map[string]bodyDecoder {
    "gzip":decompressGzip,
  }

  body=decodeBody(body,this.Headers,"transfer-encoding",transfer_decoders)
  body=decodeBody(body,this.Headers,"content-encoding",content_decoders)

  return string(body)
}

type bodyDecoder func([]byte)[]byte

func decodeBody(body []byte, headers *Headers, header string, funcs map[string]bodyDecoder) []byte {
  encoding,_:=headers.Get(header)
  log.Trace("%s: '%s'",header,encoding)
  f,found:=funcs[encoding]
  if found {
    body=f(body)
  } else if encoding!="" {
    log.Warn("unsupported %s \"%s\"",header,encoding)
  }
  return body
}

func decompressGzip(input []byte) []byte {
  bytes_reader:=bytes.NewReader(input)
  gzip_reader,err:=gzip.NewReader(bytes_reader)
  if err!=nil {
    log.Error("can't create gzip reader")
    return input
  }
  defer gzip_reader.Close()
  plain,err:=ioutil.ReadAll(gzip_reader)
  if err!=nil {
    log.Warn("invalid gzip encoding")
    return input
  }
  return plain
}
