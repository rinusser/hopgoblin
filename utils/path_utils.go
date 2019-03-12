// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package utils

import (
  "fmt"
  "os"
  "strings"
  "github.com/rinusser/hopgoblin/log"
)


var cachedBaseDir string=""


/*
  Assembles a path string relative to the application's base directory.
  Returned paths do not include a trailing path separator, so this will work with files and directories.

  Will panic if base directory can't be determined.
 */
func GetRelativePath(relative string) string {
  return GetBaseDir()+relative
}

/*
  Determines the application's base directory.
  This function will traverse the directory tree up until a "resources" subdirectory is found. This works with the main binary,
  test binaries and binaries in other packages.

  Will panic if number of parent directory traversals exceeds limit (currently 3 at most)
 */
func GetBaseDir() string {
  limit:=3

  if len(cachedBaseDir)>0 {
    return cachedBaseDir
  }

  traversal:=fmt.Sprintf("..%c",os.PathSeparator)
  for tc:=0;tc<=limit;tc++ {
    candidate:=strings.Repeat(traversal,tc)
    if lookForResourcesIn(candidate) {
      log.Debug("found basedir: \"%s\"",candidate)
      cachedBaseDir=candidate
      return cachedBaseDir
    }
  }
  panic("could not determine base directory")
}

func lookForResourcesIn(prefix string) bool {
  info,err:=os.Stat(prefix+"resources")
  return err==nil && info.IsDir()
}


/*
  Assembles a resource directory path. The returned path won't include a trailing path separator, e.g. "\" in Windows, so make sure
  you append it yourself if needed.

  Doesn't check if the target actually exists.
 */
func GetResourcePath(path string) string {
  return fmt.Sprintf("%s%c%s",GetRelativePath("resources"),os.PathSeparator,path)
}


/*
  Returns the application's executables directory. Binaries are expected to reside there.
 */
func GetApplicationDir() string {
  return fmt.Sprintf("%s%c",GetRelativePath("build"),os.PathSeparator)
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
