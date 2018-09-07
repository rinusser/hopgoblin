// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "bufio"
  "bytes"
  "fmt"
  "strconv"
  "strings"
  "github.com/rinusser/hopgoblin/log"
)


/*
  Wraps body data in HTTP "chunked" transfer encoding.
  First chunk will be of size initial_chunk_size (if enough data is available), each chunk after that will be increased by
  size_delta in size. size_delta can be 0 or negative, but chunk size will never go below 1.
 */
func ChunkEncodeBody(input string, initial_chunk_size int, size_delta int) []byte {
  builder:=strings.Builder{}
  size:=initial_chunk_size
  if size<1 {
    size=1
  }
  remainder:=input
  for {
    if len(remainder)<=size {
      size=len(remainder)
    }
    builder.WriteString(fmt.Sprintf("%x\r\n",size))
    builder.WriteString(remainder[0:size])
    builder.WriteString("\r\n")
    if len(remainder)<=size {
      break
    }
    remainder=remainder[size:]
    size+=size_delta
    if size<1 {
      size=1
    }
  }
  builder.WriteString("0\r\n\r\n")
  return []byte(builder.String())
}

func chunkDecodeBody(in *bufio.Reader, data_out *bytes.Buffer, encoding_out *bytes.Buffer) error {
  for {
    chunk_size_text,err:=in.ReadString('\n')
    if err!=nil { panic(err) }
    if encoding_out!=nil {
      encoding_out.WriteString(chunk_size_text)
    }
    chunk_size_uint64,err:=strconv.ParseInt(strings.TrimSpace(chunk_size_text),16,32)
    if err!=nil { panic(fmt.Sprintf("got invalid chunk size text '%v'",chunk_size_text)) }
    chunk_size:=int(chunk_size_uint64)
    log.Trace("found chunk with size %d (%sh)",chunk_size,strings.TrimSpace(chunk_size_text))

    remaining:=chunk_size+2 //includes \r\n at end of chunk
    for {
      burst:=make([]byte,remaining)
      size,err:=in.Read(burst)
      if err!=nil { panic(err) }
      log.Debug("got chunk data: read %d of %d bytes left in chunk",size,remaining)
      remaining-=size
      payload:=[]byte{}
      encoding:=[]byte{}
      if remaining>=2 {
        payload=burst[0:size]
      } else {
        log.Trace("size=%d, remaining=%d",size,remaining)
        if size>=2 {
          payload=burst[0:size-2+remaining]
          encoding=burst[size-2+remaining:size]
        } else {
          encoding=burst[:size]
        }
      }
      data_out.Write(payload)
      if encoding_out!=nil {
        encoding_out.Write(encoding)
      }
      if remaining<=0 {
        break
      }
    }

    if chunk_size==0 {
      log.Trace("handled end chunk, stopping read")
      break
    }
  }

  return nil
}

/*
  Decodes HTTP body with "chunked" transfer encoding.
  Will panic if data isn't encoded properly.
 */
func ChunkDecodeBody(input []byte) []byte {
  buf:=&bytes.Buffer{}
  in:=bufio.NewReader(bytes.NewReader(input))
  /*err:=*/chunkDecodeBody(in,buf,nil)
  return buf.Bytes()
}
