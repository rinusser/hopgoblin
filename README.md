# Synopsis

A custom HTTP/HTTPS proxy supporting per-site plugins. This application is NOT intended as a general-purpose proxy server for day-to-day browsing.


# General

### Use Responsibly

This proxy server purposefully violates the HTTP standards and breaks SSL certificate chains. Individual site handlers (i.e. code responsible for specific target hosts) may block, redirect or even modify individual requests and their responses.

Use this software responsibly: use it only on your own computers/networks or after getting permission to do so.

### Why Go?

I needed an application like this and had a short list of languages to write it in. A few weeks prior I had read the basics of Go and wanted to use it an actual project to make sure my initial impression was fair - so I took the opportunity and wrote this application in Go instead.

Please don't use this code as reference. I didn't particularly enjoy writing it and feel like this is reflected in the resulting code quality.

### You spelled...

It's a pun.


# Requirements

You'll need at least:

* a Go development environment (I'm using 1.10)
* an upstream HTTP/HTTPS proxy
* at least one SSL certificate keypair (can be self-signed) if you'll be handling HTTPS requests

If you want to use additional development/build tools provided:
* Linux, or Windows with Linux utilities installed (make, grep, find etc.)


# Installation

First you'll need to fetch the sources:

    go get github.com/rinusser/hopgoblin

This will add the application to src/github.com/rinusser/hopgoblin/ in your GOPATH. All following paths in these installation
instructions are relative to this directory.

Regardless of whether you use the included Makefile or not, you'll need to install a dependency to be able to run tests:

    go get github.com/stretchr/testify

You will need an SSL keypair to handle HTTPS connections: add the certificate as resources/certs/test.pem and the keyfile as
resources/certs/test.key. The certificate needs to be valid for at least `direct.local` and `proxied.local` - to easily identify it
it's recommended to use a distinct DN (or first alt name), e.g. `hopgoblin.localhost`. You'll need to add the CA you signed your
certificate(s) with into the resources/certs/ directory, its filename must start with "CA-", otherwise the automated tests will
fail.

Alternatively you can disable HTTPS support by editing resources/application.ini and commenting out the
server.default\_certificate\_file setting.

### GNU Make

There's a Makefile included, just run make:

    make

This will call the `clean`, `build` and `test` targets.

There are additional targets available:

* `make run` starts the application
* `make check` runs various code checks, will e.g. look for undocumented functions and types. Some checks will take a few minutes under Cygwin.
* `make doc` starts godoc, documentation will be available on http://127.0.0.1:64079/pkg/github.com/rinusser/hopgoblin/?m=all
* `make todo` checks source files for TODO comments

The `run` and `test` targets will pass on parameters set in the `ARGS` environment variable, for example call

    ARGS=-log=trace make test

(or `set ARGS=...` then `make test` in Windows) to run the test suite with logging set to the highest verbosity.

### Manual

Instead you can also compile the application manually - at least the main binary needs to be built:

    go build -o build/hopgoblin[.exe] github.com/rinusser/hopgoblin/main

(set the OS-appropriate executable extension, if any)

# Usage

The server is started by running the main executable in the build/ directory:

    hopgoblin

By default the server will listen on TCP port 64080 - enter localhost:64080 as your http/https proxy in your browser/client.


# Tests

Contains automated unit and integration tests. The test code contains a lot of intentional sleep time - the entire test suite can easily take 30+ seconds to finish.

### GNU Make

The tests can be run with:

    make test

### Manual

Alternatively, you can invoke `go test` directly:

    go test ./...

Some tests require the dummyproxy binary to be built - if that didn't happen yet, do so manually (include executable extension as
required):

    go build -o build/dummyproxy[.exe] github.com/rinusser/hopgoblin/http/dummyproxy/main


# Legal

### Copyright

Copyright (C) 2018 Richard Nusser

### License

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

