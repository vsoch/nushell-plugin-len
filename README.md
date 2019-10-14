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
