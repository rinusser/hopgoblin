// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

/*
  Logging facility supporting log levels.

  Log messages are printed to stderr. Only messages at or above the active log level threshold are output. Thresholds can be set
  via command line. Individual settings can be passed in separate --log arguments, or comma-separated:

    app --log=pkg1/=TRACE --log=pkg2/=DEBUG,pkg3/=OFF

  this will show messages prefixed "pkg1/" at level TRACE or above, prefixed "pkg2/" at DEBUG or above, and no messages prefixed
  "pkg3/".

  If either no prefix:

    app --log=WARN

  or the prefix "*" is specified:

    app --log=*=WARN

  the default level is set.

  Prefixes are case sensitive, levels are not.

  This logging facility routes logs message through Go's built-in "log" package, although none of its features are used. You can
  redirect the log output by setting a custom log writer in Go's "log" package.

  Go's built-in log facility doesn't allow changing the date format in log messages. This package does, and will by default show
  dates in yyyy-mm-dd format (the international standard), and times with millisecond precision.
 */
package log
