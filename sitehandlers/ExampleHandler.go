// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package sitehandlers

import (
  "bufio"
  "crypto/tls"
  "github.com/rinusser/hopgoblin/http"
  "github.com/rinusser/hopgoblin/utils"
)


func init() {
  http.RegisterSiteHandler(NewExampleHandler())
}


/*
  Example site handler for asdf.com.
 */
type ExampleHandler struct {
  utils.MultiRegexMatcher
}

/*
  required by http.SiteHandler interface
 */
func (h ExampleHandler) HandlesHost(host string) bool {
  return h.MatchesAnyRegex(host)
}

/*
  required by http.SiteHandler interface
 */
func (h ExampleHandler) HandleRequest(server *http.Server, browserio *bufio.ReadWriter, request *http.Request) {
  var response *http.Response
  var err error
  client:=http.NewClient()
  client.CopyProxySettings(server)
  response,err=client.ForwardRequest(*request)
  if err!=nil {
    return
  }
  server.WriteAndFlush(browserio,response.ToString())
}

/*
  required by http.SiteHandler interface
 */
func (this ExampleHandler) GetCertificateMap() map[string]*tls.Certificate {
  //commented out lines work, but require sitehandlers/certs/asdf.com.pem and sitehandlers/certs/asdf.com.key files.

  //cert:=utils.LoadCertificate(utils.ResolveRelativePath("certs"),"asdf.com")
  return map[string]*tls.Certificate {
  //  "asdf.com":cert,
  //  "*.asdf.com":cert,
  }
}


/*
  Creates a new ExampleHandler instance.
 */
func NewExampleHandler() http.SiteHandler {
  var h=ExampleHandler {
    MultiRegexMatcher: utils.NewMultiRegexMatcher([]string {
      `(^|\.)asdf\.com$`,
    }),
  }

  return h
}

