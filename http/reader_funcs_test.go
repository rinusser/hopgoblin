// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "bufio"
  "io"
  "time"
)


type ReadHTTPMessageAsStringTestcase struct {
  input string
  expectation string
  burst_sizes []byte
  message string
}

func CreateRHMASTestcase(input string, expectation string, burst_sizes []byte, message string) ReadHTTPMessageAsStringTestcase {
  //TODO: change burst_sizes arg to []int: even though it's currently not necessary, []byte is misleading
  //TODO: make private
  rv:=ReadHTTPMessageAsStringTestcase {
    input:input,
    burst_sizes:burst_sizes,
    message:message,
  }
  if len(rv.expectation)<1 {
    rv.expectation=input
  } else {
    rv.expectation=expectation
  }

  return rv
}

/*
  Makes sure ReadHTTPMessageAsString() minds the Content-Length header and the "chunked" transfer encoding correctly.
  This test will take ~30s to run: it contains intentional delays to simulate network data coming in multiple bursts.
 */
func TestReadHTTPMessageAsString(t *testing.T) {
  chunked_text:="HTTP/1.1 200 OK\r\nTransfer-Encoding: chunked\r\n\r\n3\r\nyo\n\r\n4\r\nmama\r\n10\r\nis a nice lady!!\r\n0\r\n\r\n"
  tiny_chunked:="HTTP/1.1 200 OK\r\nTransfer-Encoding: chunked\r\n\r\n1\r\nx\r\n0\r\n\r\n"

  cases:=[]ReadHTTPMessageAsStringTestcase {
    CreateRHMASTestcase("GET / HTTP/1.1\r\n\r\n","",[]byte{},"plain GET request"),
    CreateRHMASTestcase("GET / HTTP/1.1\r\nContent-Length:1\r\n\r\n","",[]byte{},"GET request with content length"),
    CreateRHMASTestcase("GET / HTTP/1.1\r\nContent-Length:4\r\n\r\n\x00\x01\x02\x03","",[]byte{},"GET request with binary 0"),
    CreateRHMASTestcase("POST /asdf HTTP/1.1\r\nContent-Length:2\r\n\r\nyo","",[]byte{},"plain POST request"),
    CreateRHMASTestcase(chunked_text,"",[]byte{},"chunked transfer encoding; chunk lengths need to be parsed as hex numbers!"),
    CreateRHMASTestcase(chunked_text,"",[]byte{47,13,11,2,13},"reader should wait for delayed chunks"),
    CreateRHMASTestcase("GET / HTTP/1.1\r\n\r\n","GET / HTTP/1.1\r\n\r\nasdf",[]byte{},"HTTP body without content-length should be dropped"),

    //the following lines test whether the server gracefully handles partially received chunked encoding - for example when the
    // 2 bytes of an expected CRLF newline are transmitted with a pause in between
    CreateRHMASTestcase(tiny_chunked,"",[]byte{50},      "burst just before chunk trailer, full rest"),
    CreateRHMASTestcase(tiny_chunked,"",[]byte{50,1,1,1},"burst just before chunk trailer, step over"),
    CreateRHMASTestcase(tiny_chunked,"",[]byte{50,1,2},  "burst just before chunk trailer, step, jump out"),
    CreateRHMASTestcase(tiny_chunked,"",[]byte{50,2,1},  "burst just before chunk trailer, jump into, step out"),
    CreateRHMASTestcase(tiny_chunked,"",[]byte{50,2,2},  "burst just before chunk trailer, jump into, jump out"),
    CreateRHMASTestcase(tiny_chunked,"",[]byte{51},      "burst at chunk trailer, full rest"),
    CreateRHMASTestcase(tiny_chunked,"",[]byte{51,1,1,1},"burst at chunk trailer, step over"),
    CreateRHMASTestcase(tiny_chunked,"",[]byte{51,2,1},  "burst at chunk trailer, jump across"),
    CreateRHMASTestcase(tiny_chunked,"",[]byte{51,3,1},  "burst at chunk trailer, jump over"),
    CreateRHMASTestcase(tiny_chunked,"",[]byte{52},      "burst in the middle of chunk trailer, full rest"),
    CreateRHMASTestcase(tiny_chunked,"",[]byte{52,1,1,1},"burst in the middle of chunk trailer, step over"),
    CreateRHMASTestcase(tiny_chunked,"",[]byte{52,2},    "burst in the middle of chunk trailer, jump out"),
    CreateRHMASTestcase(tiny_chunked,"",[]byte{53},      "burst after chunk trailer, full rest"),
    CreateRHMASTestcase(tiny_chunked,"",[]byte{53,1,1,1},"burst after chunk trailer, step over"),
    CreateRHMASTestcase(tiny_chunked,"",[]byte{56,1,1},"burst at last chunk trailer, step over"),
    CreateRHMASTestcase(tiny_chunked,"",[]byte{56,2},"burst at last chunk trailer, jump to end"),
    CreateRHMASTestcase(tiny_chunked,"",[]byte{57,1},"burst in the middle of last chunk trailer, step to end"),
    CreateRHMASTestcase(tiny_chunked,"",[]byte{58},"burst after last chunk trailer"),
  }
  for _,c:=range cases {
    runReadHTTPMessageAsStringTest(t,c)
  }
}

func runReadHTTPMessageAsStringTest(t *testing.T, c ReadHTTPMessageAsStringTestcase) {
  if len(c.expectation)<1 {
    c.expectation=c.input
  }

  reader,pusher:=io.Pipe()
  buf:=bufio.NewReadWriter(bufio.NewReader(reader),bufio.NewWriter(&io.PipeWriter{}))
  go func() {
    data:=[]byte(c.input)
    for _,size:=range c.burst_sizes {
      pusher.Write(data[0:size])
      data=data[size:]
      time.Sleep(5e8)
    }
    pusher.Write(data)
    pusher.Close()
  }()

  result,err:=ReadHTTPMessageAsString(buf)
  assert.Nil(t,err,c.message)
  assert.Equal(t,c.expectation,result,c.message)
}
