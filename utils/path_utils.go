// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package utils

import (
  "fmt"
  "os"
  "runtime"
  "strings"
  "github.com/rinusser/hopgoblin/log"
)


var basedir_cache string=""

/*
  Assembles a resource directory path. The returned path will include a trailing path separator, e.g. "\" in Windows.

  Doesn't check if the directory actually exists.
 */
func GetResourceDir(subdir string) string {
  return fmt.Sprintf("%s%s%c",GetResourceBaseDir(),subdir,os.PathSeparator)
}

/*
  Finds the resource base directory.
  Test binaries are executed in their package directory, so this traverses the directory tree up until "resources/" is found.

  Will panic if number of parent directory traversals exceeds limit (currently 3 at most)
 */
func GetResourceBaseDir() string {
  limit:=3

  if len(basedir_cache)>0 {
    return basedir_cache
  }
  traversal:=fmt.Sprintf("..%c",os.PathSeparator)
  for tc:=0;tc<=limit;tc++ {
    prefix:=strings.Repeat(traversal,tc)
    if lookForResourcesIn(prefix) {
      basedir_cache=fmt.Sprintf("%s%s%c",prefix,"resources",os.PathSeparator)
      log.Debug("found resource basedir: \"%s\"",basedir_cache)
      return basedir_cache
    }
  }
  panic("could not locate resources/ directory")
}

func lookForResourcesIn(prefix string) bool {
  info,err:=os.Stat(prefix+"resources")
  return err==nil && info.IsDir()
}


/*
  Returns the application's main directory. Binaries are expected to reside there.
 */
func GetApplicationDir() string {
  return GetResourceDir("../build")
}


/*
  Turns a sourcefile-relative path into an absolute path for file I/O.
  If the resulting path should be used as a directory, make sure to append the path separator to the resulting string before
  appending adding a file name.

  This function won't check if the relative path points to an existing file system entry, it just creates the new path string.
 */
func ResolveRelativePath(path string) string {
  _,file,_,_:=runtime.Caller(1)
  dir:=dirname(file)
  return dir+path
}

func dirname(path string) string {
  end:=strings.LastIndex(path,string(os.PathSeparator))
  if end<0 {
    end=strings.LastIndex(path,"/")
  }
  if end<0 {
    log.Warn("function argument doesn't seem to be a path")
    return path+string(os.PathSeparator)
  }
  return path[0:end+1]
}
