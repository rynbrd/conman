ConMan
======
ConMan is a container initializer which uses environment variables, command
line arguments, and a config file to render templates, set environment
variables, and execute a binary. It is meant to be used as a container's init
so as to avoid the tedious task of writing a bash script of some sort to do all
of this for you.

[![Build Status](https://travis-ci.org/BlueDragonX/conman.svg?branch=master)](https://travis-ci.org/BlueDragonX/conman)

Operation
---------
ConMan constructs a context object which is a map of key/value pairs. The
initial context is loaded from the `context` value in the config file.
Environment variables are then added to the context followed by values set on
the command line.

The context is used to render the template files, environment values, and
command line arguments. The templates are rendered prior to calling `exec`.

Is it important to make note of the order in which values are added to the
context. The order is as follows:

1. Config file context.
2. ConMan's environment.
3. System context.
4. Command line arguments.
5. Config file environment.

A consequence of this is that environment variables defined in the config file
are not available to be used as context to templated environment variables. In
addition, it is not possible to overwrite those environment variables by
setting a value on the command line.

Config File
-----------
The config file is YAML. The following working example describes the structure:

	# An initial context object to start with.
	context:
	  nope: but yes
	  env:
		GREETING: Hello
		SUBJECT: world

	# Templates to render. The key is the destination and the value is the source.
	templates:
	  example.txt: example.tpl

	# Update the environment. Values are templated.
	env:
	- GREETING=Greetings
	- SUBJECT=Mr. {{ title .env.SUBJECT }}

	# Configure the executable to launch. Arguments may use Go template syntax. The
	# context is the same as for templates.
	exec:
	- /bin/echo
	- '{{ .env.GREETING }}, {{ .env.SUBJECT }}!'

The `context` section provides and initial context structure. This is most
useful for default values.

The `templates` section contains a map of templates. The keys are the
destination to write the rendered template to while the value is the source.
The keys and values themselves may be templated.

The `env` sections contains a list of environment variables to set for the
exec'd binary. These values may be templated. 

Lastly the `exec` section is a list of arguments representing the program to be
exec'd. These values may be templated.

Command Line
------------
The following command line options are recognized:

    -help           Print the help.
    -var=NAME=VALUE Set a named value.
    -json=JSON      Set context from the provided JSON object.
    -config=FILE    Load configuration from this file. Defaults to
                    /etc/conman.yml.

The `var` and `json` options load values into the context. They may be used to
initialize values in `sys` and `env` but will be overwritten if those values
are set by their respective modules.

Environment Variables
---------------------
Environment variables are available as a map under the context name `env`.
Undeclared environment variables will result in a `<no value>` in the template.

System Context
--------------
The system context resides under the context value `sys`. It contains a map
containing system related values. Currently these values include:

`ipaddress` - The IPv4 address.
`network` - The IPv4 network in CIDR format.
`gateway` - The IPv4 default gateway.
`hostname` - The system's hostname.

The `ipaddress` and `network` values are derived by first determining the
"primary" interface. This is the interface associated with the first default
gateway listed in /proc/net/route.

Template Functions
------------------
A number of template functions have been included to (hopefully) make your life
easier. They are:

`replace` - String replacement.
`join` - Join an array of strings.
`split` - Split a string.
`title` - Title case a string.
`upper` - Uppercase a string.
`lower` - Lowercase a string.
`json` - Unmarshal JSON into an object or array.
`addrHost` - Get the host part of a host:port formatted address.
`addrPort` - Get the port part of a host:port formatted address.
`urlScheme` - Get the scheme part of a URL.
`urlHost` - Get the host part of a URL. This include the :port if present.
`urlUsername` - Get the username part of a URL.
`urlPassword` - Get the password part of a URL.
`urlRawQuery` - Get the URL's query string.
`urlQuery` - Get the first value of a query key. Takes `name` as an additional parameter.
`urlFragment` - Get the fragment part of the URL.

License
-------
Copyright (c) 2015 Ryan Bourgeois. Licensed under BSD-Modified. See the
[LICENSE][1] file for a copy of the license.

[1]: https://raw.githubusercontent.com/BlueDragonX/conman/master/LICENSE "Sentinel License"
