// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package dummyproxy

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "fmt"
  "net"
  go_http "net/http"
  "time"
)


/*
  Makes sure the dummy proxy only listens on the port given in the command line.
 */
func TestDummyProxyRunnerPortArgument(t *testing.T) {
  proxyrunner:=NewDummyProxyRunner()
  proxyrunner.Port=64082
  err:=proxyrunner.Start()
  if err!=nil { panic(err) }

  time.Sleep(5e8)
  expectSuccess(t,64082)
  expectFailure(t,64083)
//  proxyrunner.ReadAndWait()

  proxyrunner=NewDummyProxyRunner()
  proxyrunner.Port=64083
  err=proxyrunner.Start()
  if err!=nil { panic(err) }

  time.Sleep(5e8)
  expectFailure(t,64082)
  expectSuccess(t,64083)
//  proxyrunner.ReadAndWait()
}


/*
  Makes sure the server can run in parallel on multiple, random ports
 */
func TestStartRandom(t *testing.T) {
  count:=5
  ports:=[]int{}

  for tc:=1;tc<=count;tc++ {
    ports=append(ports,setupStartRandomTest(t,tc))
  }

  for _,port:=range ports {
    expectSuccess(t,port)
  }
}

func setupStartRandomTest(t *testing.T, index int) int {
  runner:=NewDummyProxyRunner()
  port:=runner.StartRandom()
  assert.Equal(t,port,runner.Port,"return value should match .Port field, index "+string(index))
  return port
}


func performRequest(port int) (*go_http.Response,error) {
  transport:=&go_http.Transport {
    Dial: (&net.Dialer{
      Timeout: 2*time.Second,
    }).Dial,
  }
  client:=&go_http.Client {
    Transport: transport,
    Timeout: time.Duration(2*time.Second),
  }
  result,err:=client.Get(fmt.Sprintf("http://127.0.0.1:%d/proxytest",port))
  return result,err
}

func expectSuccess(t *testing.T, port int) {
  result,err:=performRequest(port)

  assert.Nil(t,err,"client.Get() failed")
  assert.NotNil(t,result)
  if result==nil {
    return
  }
  defer result.Body.Close()
  assert.Equal(t,418,result.StatusCode)
}

func expectFailure(t *testing.T, port int) {
  result,err:=performRequest(port)

  assert.NotNil(t,err,"client.Get() should have errored out")
  assert.Nil(t,result,"there shouldn't have been a result")
}
