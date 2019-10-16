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
{"jsonrpc":"2.0","method":"response","params":{"Ok":{"name":"len","usage":"Return the length of a string","positional":[],"rest_positional":null,"named":{},"is_filter":true}}}
```

### Start and End Filter

```bash
$ ./nu_plugin_len
{"method":"begin_filter"}
{"jsonrpc":"2.0","method":"response","params":{"Ok":[]}}

{"method":"end_filter"}
{"jsonrpc":"2.0","method":"response","params":{"Ok":[]}}
```

### Calculate Length

```bash
$ ./nu_plugin_len
{"method":"begin_filter"}
{"jsonrpc":"2.0","method":"response","params":{"Ok":[]}}

{"method":"filter", "params": {"item": {"Primitive": {"String": "oogabooga"}}}}
{"jsonrpc":"2.0","method":"response","params":{"Ok":{"Value":{"item":{"Primitive":{"Int":9}}}}}}

{"method":"end_filter"}
{"jsonrpc":"2.0","method":"response","params":{"Ok":[]}}
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
$ docker run -it vanessa/nu-plugin-len
```

Once inside, you can use `nu -l trace` to confirm that nu found your plugin.
Here we see that it did!

```bash
/code(add/circleci)> nu -l trace
...
 TRACE nu::cli > Trying "/usr/local/bin/nu_plugin_len"
 TRACE nu::cli > processing response (176 bytes)
 TRACE nu::cli > response: {"jsonrpc":"2.0","method":"response","params":{"Ok":{"name":"len","usage":"Return the length of a string","positional":[],"rest_positional":null,"named":{},"is_filter":true}}}

 TRACE nu::cli > processing Signature { name: "len", usage: "Return the length of a string", positional: [], rest_positional: None, named: {}, is_filter: true }
 TRACE nu::data::config > config file = /root/.config/nu/config.toml
```

You can also (for newer versions of nu > 0.2.0) use help to see the command:

```bash
/code(master)> help len
Return the length of a string

Usage:
  > len 

/code(master)> 
```

Try out calculating the length of something! Here we are in a directory with
one file named "myname" that is empty.

```
/tmp/test> ls
━━━━━━━━┯━━━━━━┯━━━━━━━━━━┯━━━━━━┯━━━━━━━━━━━━━━━━┯━━━━━━━━━━━━━━━━
 name   │ type │ readonly │ size │ accessed       │ modified 
────────┼──────┼──────────┼──────┼────────────────┼────────────────
 myname │ File │          │  —   │ 41 seconds ago │ 41 seconds ago 
━━━━━━━━┷━━━━━━┷━━━━━━━━━━┷━━━━━━┷━━━━━━━━━━━━━━━━┷━━━━━━━━━━━━━━━━
```

Try listing, getting the name, and calculating the length.

```bash
/tmp/test> ls | get name | len
━━━━━━━━━━━
 <unknown> 
───────────
         6 
━━━━━━━━━━━
```

or test it out with debug.

```bash
ls | get name | len | debug

/tmp/test> ls | get name | len | debug
Tagged { tag: Tag { anchor: None, span: Span { start: 0, end: 2 } }, item: Primitive(Int(BigInt { sign: Plus, data: BigUint { data: [6] } })) }
━━━━━━━━━━━
 <unknown> 
───────────
         6 
━━━━━━━━━━━
```

Add another file to see the table get another row

```bash
touch four
```
```bash
/tmp/test> ls | get name | len 
━━━┯━━━━━━━━━━━
 # │ <unknown> 
───┼───────────
 0 │         4 
 1 │         6 
━━━┷━━━━━━━━━━━
```

Mind you, I'm not a wizard Go Programmer, but I'd like the community to 
at least have an example to start with! Please contribute to this plugin to make
it better!
