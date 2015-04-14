package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"
)

const (
	DefaultConfigFile = "/etc/conman.yml"
)

func Fatal(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
	os.Exit(1)
}

func main() {
	// parse command line options
	vars := MapVar{}
	json := JsonVar{}
	configFile := ""
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flags.StringVar(&configFile, "config", DefaultConfigFile, "Load configuration from this file.")
	flags.Var(&vars, "var", "Add a value to the context. Formatted as `name=value`.")
	flags.Var(&json, "json", "Add the contents of the JsonVar object to the context.")
	flags.Parse(os.Args[1:])

	// retrieve configuration
	config := &Config{}
	if err := config.Read(configFile); err != nil {
		Fatal("%s\n", err)
	}

	if len(config.Exec) == 0 {
		Fatal("no exec args\n")
	}

	// build the environment and context
	environ := &Environ{}
	environ.Load(os.Environ())
	context := &Context{}
	context.Update(config.Context)
	context.Update(map[string]interface{}{"env": environ.Context()})
	context.Update(config.Context)
	context.Update(vars)
	context.Update(json)

	renderedEnv := make([]string, len(config.Env))
	for n, envVar := range config.Env {
		if renderedEnvVar, err := RenderString(envVar, context.Map()); err == nil {
			renderedEnv[n] = renderedEnvVar
		} else {
			Fatal("%s\n", err)
		}
	}
	environ.Load(renderedEnv)
	context.Update(map[string]interface{}{"env": environ.Context()})

	// render the exec args
	args := make([]string, len(config.Exec))
	for n, arg := range config.Exec {
		if renderedArg, err := RenderString(arg, context.Map()); err == nil {
			args[n] = renderedArg
		} else {
			Fatal("%s\n", err)
		}
	}

	// render the templates
	for dest, src := range config.Templates {
		if err := (&Template{src, dest}).Render(context.Map()); err != nil {
			Fatal("%s\n", err)
		}
	}

	// exec the command
	if err := syscall.Exec(args[0], args, environ.Values()); err != nil {
		Fatal("%s\n", err)
	}
}
