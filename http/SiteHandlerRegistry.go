// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

package http


var registered_site_handlers []SiteHandler

/*
  Register a site handler. Call this e.g. in other packages' init() functions to inject site handlers into the server.
 */
func RegisterSiteHandler(handler SiteHandler) {
  registered_site_handlers=append(registered_site_handlers,handler)
}

/*
  Returns all previously registered site handlers.
 */
func GetRegisteredSiteHandlers() []SiteHandler {
  return registered_site_handlers
}
