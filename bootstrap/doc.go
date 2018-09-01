// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

/*
  Bootstrap helper, used mostly to register/parse command-line arguments in the correct order.

  To allow individual packages to parse command line arguments into global settings (e.g. log levels) without interfering with
  each other a custom mechanism was required. Additionally, Go's test runner doesn't support adding a custom bootstrapper to an
  entire test suite (asking developers to copy the same test runner code over and over into each package instead).

  This bootstrap mechanism works around the first issue and alleviates the second.

  To run package initialization code once command-line arguments have been parsed, register a hook in your init() function:

    import "flag"
    import "github.com/rinusser/hopgoblin/bootstrap"

    var argument=flag.Int("argument",1234,"some integer argument")

    func init() {
      bootstrap.AfterFlagParse(initHook)
    }

    func initHook() {
      //handle argument value here
    }

  Then, assuming the application's main() function somehow calls bootstrap.Init(), your hook will be called once your command-line
  parameters are available.

  Keep in mind Go loads packages only as required: a package's init() methods, and thus any hooks attempted to register, won't be
  invoked unless the package is imported somewhere. See for example the main/main.go file where the "sitehandlers" package is
  imported with a dummy alias.
 */
package bootstrap
