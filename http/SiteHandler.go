// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http

import (
  "bufio"
  "crypto/tls"
)


/*
  All site handlers must implement this interface.

  HandlesHost() should return true if the site handler instance is responsible for handling requests to this host.

  HandleRequest() gets called for any incoming requests the site handler is responsible for.

  GetCertificateMap() should return a mapping of hostnames to certificates. Supports wildcards, e.g. "*.example.com". Make sure to
  include a mapping for the base domain (e.g. "example.com") if you want to match that as well.
 */
type SiteHandler interface {
  HandlesHost(host string) bool
  HandleRequest(server *Server, buf *bufio.ReadWriter, request *Request) //TODO: clear up, probably change buf to Writer
  GetCertificateMap() map[string]*tls.Certificate
}
