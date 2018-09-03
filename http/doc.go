// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

/*
  Package handles HTTP messages and network I/O.
  Handles HTTP requests/responses, communication with proxy servers, TLS encryption and so on.

  The Client code will use a CA certificate pool to validate remote certificates consisting of the system's list of CA
  certificates, and any additional certificate files from the resources/certs/ directory that start with "CA-". Currently there's
  no need to add additional CA certificates, as remote certificate checks are disabled.
 */
package http
