// Copyright 2018 Richard Nusser
// Licensed under GPLv3 (see http://www.gnu.org/licenses/)

/*
  Site handlers: code managing incoming HTTP requests.

  Handlers must be registered with http.RegisterSiteHandler(), e.g. from an init() function. Handlers must implement the
  http.SiteHandler interface.

  Handlers' HandleRequest() methods can answer requests themselves, forward requests to another server/proxy, analyze or modify
  requests and responses, ... - anything goes!

  Note that intercepting HTTPS connections will trigger certificate warnings/errors in the connecting client (e.g. the browser).
  It's recommended that you create a self-signed certificate chain, load custom certificates (with appropriate hostnames entered)
  in the site handler and add your CA file to the browser (ideally in a separate profile just for this purpose, so you don't run a
  man-in-the-middle attack on yourself by accident).
 */
package sitehandlers
