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
  burstSizes []int
  message string
}

func createRHMASTestcase(input string, expectation string, burst_sizes []int, message string) ReadHTTPMessageAsStringTestcase {
  rv:=ReadHTTPMessageAsStringTestcase {
    input:input,
    burstSizes:burst_sizes,
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
    createRHMASTestcase("GET / HTTP/1.1\r\n\r\n","",[]int{},"plain GET request"),
    createRHMASTestcase("GET / HTTP/1.1\r\nContent-Length:1\r\n\r\n","",[]int{},"GET request with content length"),
    createRHMASTestcase("GET / HTTP/1.1\r\nContent-Length:4\r\n\r\n\x00\x01\x02\x03","",[]int{},"GET request with binary 0"),
    createRHMASTestcase("POST /asdf HTTP/1.1\r\nContent-Length:2\r\n\r\nyo","",[]int{},"plain POST request"),
    createRHMASTestcase(chunked_text,"",[]int{},"chunked transfer encoding; chunk lengths need to be parsed as hex numbers!"),
    createRHMASTestcase(chunked_text,"",[]int{47,13,11,2,13},"reader should wait for delayed chunks"),
    createRHMASTestcase("GET / HTTP/1.1\r\n\r\n","GET / HTTP/1.1\r\n\r\nasdf",[]int{},"HTTP body without content-length should be dropped"),

    //the following lines test whether the server gracefully handles partially received chunked encoding - for example when the
    // 2 bytes of an expected CRLF newline are transmitted with a pause in between
    createRHMASTestcase(tiny_chunked,"",[]int{50},      "burst just before chunk trailer, full rest"),
    createRHMASTestcase(tiny_chunked,"",[]int{50,1,1,1},"burst just before chunk trailer, step over"),
    createRHMASTestcase(tiny_chunked,"",[]int{50,1,2},  "burst just before chunk trailer, step, jump out"),
    createRHMASTestcase(tiny_chunked,"",[]int{50,2,1},  "burst just before chunk trailer, jump into, step out"),
    createRHMASTestcase(tiny_chunked,"",[]int{50,2,2},  "burst just before chunk trailer, jump into, jump out"),
    createRHMASTestcase(tiny_chunked,"",[]int{51},      "burst at chunk trailer, full rest"),
    createRHMASTestcase(tiny_chunked,"",[]int{51,1,1,1},"burst at chunk trailer, step over"),
    createRHMASTestcase(tiny_chunked,"",[]int{51,2,1},  "burst at chunk trailer, jump across"),
    createRHMASTestcase(tiny_chunked,"",[]int{51,3,1},  "burst at chunk trailer, jump over"),
    createRHMASTestcase(tiny_chunked,"",[]int{52},      "burst in the middle of chunk trailer, full rest"),
    createRHMASTestcase(tiny_chunked,"",[]int{52,1,1,1},"burst in the middle of chunk trailer, step over"),
    createRHMASTestcase(tiny_chunked,"",[]int{52,2},    "burst in the middle of chunk trailer, jump out"),
    createRHMASTestcase(tiny_chunked,"",[]int{53},      "burst after chunk trailer, full rest"),
    createRHMASTestcase(tiny_chunked,"",[]int{53,1,1,1},"burst after chunk trailer, step over"),
    createRHMASTestcase(tiny_chunked,"",[]int{56,1,1},"burst at last chunk trailer, step over"),
    createRHMASTestcase(tiny_chunked,"",[]int{56,2},"burst at last chunk trailer, jump to end"),
    createRHMASTestcase(tiny_chunked,"",[]int{57,1},"burst in the middle of last chunk trailer, step to end"),
    createRHMASTestcase(tiny_chunked,"",[]int{58},"burst after last chunk trailer"),
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
    for _,size:=range c.burstSizes {
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
