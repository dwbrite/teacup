# teacup
A small web server, mostly for hosting dwbrite.com
It is designed to run without super user permissions, 
as a supplement to an existing web server configuration, such as nginx or apache.
It can also be run on "well-known" ports given special permissions.

`setcap` doesn't always work, so you may have to run as root.

## Roll your own server
For the purposes of creating simple and concrete examples, we'll be using nginx.

* Set up a PostgreSQL database and make sure the service is running.
  * I recommend using the `~/.pgpass` file, or environment variables so that no passwords are exposed.
* Add an nginx configuration routing to the port this program runs on.
  * See `examples/nginx.conf` if you intend on enabling SSL/TLS access.
  * Make sure the ports used by nginx are open and being forwarded to the server machine.
* Create and deploy your SSL/TLS certificates. I use Let's Encrypt + `certbot` for my certs.
  * You'll want to create deploy scripts which copy the required files with proper permissions 
  to the user running the go server.
* Install the `pq` PostgreSQL driver with `go get github.com/lib/pq`* (is this necessary?)
* Install `teacup` with `go get -v github.com/dwbrite/teacup/...`
  * Alternatively, clone this repository.
* Start writing your own little website today :)

## License
Copyright (C) 2018 Devin Brite  
SPDX-License-Identifier: GPL-3.0-or-later OR  CC0-1.0 OR MIT

Full texts can be found in `LICENSE.md`