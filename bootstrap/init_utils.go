// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package bootstrap

import (
  "flag"
)


type hookType func()

var after_flag_parse_hooks []hookType


/*
  Registers a hook to be called after flag.Parse() is done.
  Put initialization code that acts upon command-line arguments there.
 */
func AfterFlagParse(hook hookType) {
  after_flag_parse_hooks=append(after_flag_parse_hooks,hook)
}


/*
  Bootstraps the application.
  Calls flag.Parse() - applications should register any custom flags before calling this. Parsing those flags can be done in hooks
  registered with AfterFlagParse()
 */
func Init() {
  flag.Parse()
  for _,f:=range after_flag_parse_hooks {
    f()
  }
}
