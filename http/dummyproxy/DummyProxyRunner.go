// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package dummyproxy

import (
  "bufio"
  "fmt"
  "io"
  "io/ioutil"
  "os"
  "os/exec"
  "strings"
  "github.com/rinusser/hopgoblin/bootstrap"
  "github.com/rinusser/hopgoblin/log"
  "github.com/rinusser/hopgoblin/utils"
)


var proxyPath string="../../build/dummyproxy" //value will be changed by initHook()


func init() {
  bootstrap.AfterFlagParse(initHook)
}

func initHook() {
  proxyPath=utils.GetApplicationDir()+utils.GetConfigValue("test.proxy_executable_basename")
  executable,err:=os.Executable()
  if err!=nil {
    return
  }
  pos:=strings.LastIndex(executable,".")
  if pos>1 && executable[pos:]!=".test" {
    proxyPath+=executable[pos:]
  }
}


/*
  Wrapper for dummy proxy.
  Change port as required.
 */
type DummyProxyRunner struct {
  Port int
  proc *exec.Cmd
  stdout io.ReadCloser
  stderr io.ReadCloser
}

/*
  Create a new instance, will by default listen on port 64086.
 */
func NewDummyProxyRunner() DummyProxyRunner {
  return DummyProxyRunner {
    Port:64086,
  }
}


/*
  Start the dummy proxy.
  Will spawn a new process that waits and handles one connection, then exits.
 */
func (this *DummyProxyRunner) Start() error {
  proc:=exec.Command(proxyPath,fmt.Sprintf("--port=%d",this.Port),log.AssemblePassthroughArg())
  this.proc=proc

  this.stdout,_=proc.StdoutPipe()
  this.stderr,_=proc.StderrPipe()
  err:=proc.Start()

  if err==nil {
    reader:=bufio.NewReader(this.stderr)
    first_line,_:=reader.ReadString('\n')
    if strings.HasPrefix(first_line,"port:") {
      portstr:=strings.TrimSpace(first_line[5:])
      fmt.Sscanf(portstr,"%d",&this.Port)
    }
  }
  return err
}


/*
  Start the dummy proxy on a random port.
  Will panic if proxy didn't start.

  The port the proxy was started can either be read from this method's return value, or from the .Port field.
 */
func (this *DummyProxyRunner) StartRandom() int {
  this.Port=0
  err:=this.Start()
  if err!=nil {
    panic(err)
  }
  if this.Port==0 {
    panic("proxy didn't start")
  }
  return this.Port
}


/*
  Wait for the dummy proxy to finish and return stdout (requests) and stderr (log) data.
 */
func (this *DummyProxyRunner) ReadAndWait() ([]byte,[]byte,error) {
  outdata,outerr:=readFromPipe(this.stdout)
  errdata,errerr:=readFromPipe(this.stderr)

  err:=errerr
  if outerr!=nil {
    err=outerr
  }

  this.proc.Wait()
  return outdata,errdata,err
}

func readFromPipe(pipe io.ReadCloser) ([]byte,error) {
  data,err:=ioutil.ReadAll(pipe)
  if data!=nil && len(data)>0 {
    data=data[0:len(data)-1]
  }
  return data,err
}
