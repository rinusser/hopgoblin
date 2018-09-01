// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package bootstrap

import (
  "testing"
  "os"
)


func TestMain(m *testing.M) {
  Init()
  os.Exit(m.Run())
}

