// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package utils

import (
  "crypto/tls"
  "os"
  "github.com/rinusser/hopgoblin/log"
)


/*
  Loads a named certificate from the given certificate directory.
  Returns nil on error.

  You can use ResolveRelativePath() or GetResourceBaseDir() to easily create absolute directory paths.
 */
func LoadCertificate(directory string, name string) *tls.Certificate {
  if directory[len(directory):]!=string(os.PathSeparator) {
    directory+=string(os.PathSeparator)
  }
  cert,err:=tls.LoadX509KeyPair(directory+name+".pem",directory+name+".key")
  if err!=nil {
    log.Error("can't load cert/key files: ",err)
    return nil
  }
  return &cert
}
