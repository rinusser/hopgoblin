// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "bufio"
  "bytes"
  "fmt"
  "io"
  "strings"
  "github.com/rinusser/hopgoblin/log"
)


func readHTTPMessageHeader(buf *bufio.ReadWriter, builder *strings.Builder) error {
  for {
    line,err:=buf.ReadString('\n')
    log.Trace("got line: %s",line)
    if err==io.EOF {
      log.Trace("got EOF, stopping read")
      break
    } else if err!=nil {
      log.Error("can't read from buffer: %v",err)
      return err
    }
//    fmt.Println("ReadHTTPMessageAsString: got line: ",strings.TrimSpace(line))
    builder.WriteString(line)
    if line=="\r\n" {
      log.Trace("found empty line, stopping read")
      break
    }
  }
  return nil
}

func readHTTPMessageBodyWithLength(buf *bufio.ReadWriter, builder *strings.Builder, length_text string) error {
  length:=0
  fmt.Sscanf(strings.TrimSpace(length_text),"%d",&length)
  log.Trace("got Content-Length header: raw=%q => value=%d",length_text,length)
  for {
    if length==0 {
      break
    }
    bufsize:=4096
    if length<bufsize {
      bufsize=length
    }
    chunk:=make([]byte,bufsize)
    length-=len(chunk)
    log.Trace("built buffer with size %d, %d remaining",len(chunk),length)
    size,err:=buf.Read(chunk)
    if size==0 || err==io.EOF {
      break
    } else if err!=nil {
      return err
    }
    builder.Write(chunk[0:size])
  }
  return nil
}

func readHTTPMessageChunkedBody(in *bufio.ReadWriter, builder *strings.Builder) error { //TODO: change in to bufio.Reader
  buf:=&bytes.Buffer{}
  err:=chunkDecodeBody(in.Reader,buf,buf)
  builder.Write(buf.Bytes())
  return err
}

/*
  Reads an entire HTTP request/response from the input stream.

  This function currently requires either the Content-Length header to be included, or the Transfer-Encoding set to "chunked".
  If neither condition is satisfied any message body after the HTTP headers will be ignored.
 */
func ReadHTTPMessageAsString(buf *bufio.ReadWriter) (string,error) { //TODO: change to IO reader
  log.Trace("starting to read http message from buffer")
  var rv strings.Builder
  err:=readHTTPMessageHeader(buf,&rv)
  if err!=nil {
    return "",err
  }
  log.Trace("finished reading headers")
  headers:=ParseHeaders(rv.String())
  log.Trace("got headers")
  length_text,found_content_length:=headers.Get("Content-Length")
  xfer_encoding_text,found_xfer_encoding:=headers.Get("Transfer-Encoding")
  if found_content_length {
    readHTTPMessageBodyWithLength(buf,&rv,length_text)
  } else if (found_xfer_encoding && xfer_encoding_text=="chunked") {
    readHTTPMessageChunkedBody(buf,&rv)
  }
  log.Trace("finished reading message, returning..")
  return rv.String(),nil
}

