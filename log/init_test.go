// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package log

import (
  "testing"
  "os"
  "github.com/rinusser/hopgoblin/bootstrap"
)


func TestMain(m *testing.M) {
  bootstrap.Init()
  os.Exit(m.Run())
}

