# Nushell Plugin in GoLang

This is an example implementation of a length (len) plugin for Nushell in GoLang. While
there are many ways to go about this, I created this particular implementation,
of course wanting to practice a bit of Go and learn about nushell. This is
incredibly challenging because GoLang doesn't make it easy to parse nested
levels of Json without knowing the structure in advance.

**under development**

I haven't tested yet with nushell, or finished the examples below. Stay tuned!

## Build

You can use the Makefile to build the plugin

```bash
$ make
go build -o nu_plugin_len
```

## Test Without Nu

It's possible to test without nushell by giving json objects to the binary.
After you run it and press enter, here are some functions to interact:

### Config

```bash
$ ./nu_plugin_len
{"method":"config"}
len
```

### Start and End Filter

```bash
$ ./nu_plugin_len
{"method":"begin_filter"}
[]
{"method":"end_filter"}
[]
```

### Calculate Length

```bash
$ ./nu_plugin_len
{"method":"begin_filter"}
[]
{"method":"filter", "params": {"item": {"Primitive": {"String": "oogabooga"}}}}
{"jsonrpc":"2.0","method":"response","params":{"Ok":{"Value":9}}}
{"method":"end_filter"}
[]
```

## Logging

Note that since I'm going to be using the plugin in a container, I don't
mind logging to a temporary file at `/tmp/nu-plugin-len.log`. If you want
to remove this, remove the logger.* snippets from [main.go](main.go) along
with the "log" import.


## Test With Nu

Once you are happy, you can install the plugin with Nushell easily via Docker.
Here we build the container using first GoLang to compile, and then
copying the binary into quay.io/nushell/nu-base in /usr/local/bin.
We do this so that the plugin is discovered. So first, build the container:

```bash
$ docker build -t vanessa/nu-plugin-len .
```

Then shell inside - the default entrypoint is already the nushell.

```bash
$ docker exec -it vanessa/nu-plugin-len
```

Once inside, you can use `nu -l trace` to confirm that nu found your plugin.
Here we see that it did!

```bash
/code(add/circleci)> nu -l trace
 TRACE nu::cli > Looking for plugins in "/usr/local/cargo/bin"
 TRACE nu::cli > Looking for plugins in "/usr/local/sbin"
 TRACE nu::cli > Looking for plugins in "/usr/local/bin"
 TRACE nu::cli > Found "nu_plugin_len"
 TRACE nu::cli > processing response (4 bytes)
```

**under development**
