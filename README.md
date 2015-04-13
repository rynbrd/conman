ConMan
======
ConMan is a container initializer which uses environment variables, command
line arguments, and a config file to render templates, set environment
variables, and execute a binary. It is meant to be used as a container's init
so as to avoid the tedious task of writing a bash script of some sort to do all
of this for you.

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

1. ConMan's environment.
2. Config file context.
3. Command line arguments.
4. Config file environment.

A consequence of this is that environment variables defined in the config file
are not available to be used as context to templated environment variables. In
addition, it is not possible to overwrite those environment variables by
setting a value on the command line.

Config File
-----------
The config file is YAML. The following working example describes the structure:

    # An initial context object to start with.
    context:
      GREETING: hello
      SUBJECT: world

    # Templates to render. The key is the destination and the value is the source.
    templates:
      /tmp/greeting: /tmp/greeting.tpl

    # Update the environment. Values are templated.
    env:
    - GREETING=greetings
    - SUBJECT=Mr. {{ .SUBJECT }}

    # Configure the executable to launch. Arguments may use Go template syntax. The
    # context is the same as for templates.
    exec:
    - /bin/echo
    - '{{ title .GREETING }}, {{ .SUBJECT }}!'

Command Line
------------
The following command line options are recognized:

    -help           Print the help.
    -var=NAME=VALUE Set a named value.
    -json=JSON      Set context from the provided JSON object.
    -config=FILE    Load configuration from this file. Defaults to
                    /etc/conman.yml.

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
`json` - Works just like `index` but takes a JSON string as the first argument.

License
-------
Copyright (c) 2015 Ryan Bourgeois. Licensed under BSD-Modified. See the
[LICENSE][1] file for a copy of the license.

[1]: https://raw.githubusercontent.com/BlueDragonX/conman/master/LICENSE "Sentinel License"
