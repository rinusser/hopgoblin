// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

/*
  Application configuration reader for logging facility.
  This package needed to be split from hopgoblin/log to avoid import loops.

  Import this package in your application, e.g. where your main() function is, to enable log configuration from the .ini file - 
  see e.g. main/main.go.

  Do not use this package in test suites.
 */
package appconfig
