// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

/*
  Wrapper for running the DummyHTTPProxy server in a separate process.

  Running the server in the same thread as the test suite resulted in deadlocks ~5% of the time without any discernible reason.
  Moving parts of the code to a separate process fixed the issue - while keeping the execution order, thus most potential causes
  for deadlocks, the same.
 */
package dummyproxy
