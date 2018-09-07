// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package bootstrap

import (
  "flag"
)


type hookType func()

var afterFlagParseHooks []hookType


/*
  Registers a hook to be called after flag.Parse() is done.
  Put initialization code that acts upon command-line arguments there.
 */
func AfterFlagParse(hook hookType) {
  afterFlagParseHooks=append(afterFlagParseHooks,hook)
}


/*
  Bootstraps the application.
  Calls flag.Parse() - applications should register any custom flags before calling this. Parsing those flags can be done in hooks
  registered with AfterFlagParse()
 */
func Init() {
  flag.Parse()
  for _,f:=range afterFlagParseHooks {
    f()
  }
}
