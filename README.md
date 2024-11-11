# Why

This repo enable a user of [zk](https://github.com/zk-org/zk) to export their graph and create a visual representations of the connections between their notes.

This repo borrow heavily from [zetteltools](https://github.com/joashxu/zetteltools) in term of concepts.
I did migrate the code to `go` so that it may one day be included with zk but also becase I am learning go ;-)

Some things I did add:
- ability to `pipe` directly the output of `zk graph` to the script so that it can generate a file
- using go `tmpl` to generate the content as this is standard in `go`
- use the latest `d3` version (version `7`) at the time of this writing (disclaimer: I have zero experience with `d3` but ... I also want to learn!)
- include an option  to serve the file directly using `htttp/server` as this is included in `go`

Becaue of this, I don't want to use any module unless the benefits would be huge.

## Future development

- [ ] Add more parameters to the template in particular the size of the `d3` canvas, the `strength` of the various forces. Most certainly in a configuration file in the `.zk` directory?
- [ ] Keep the first `200` charcters of the content of each note to offer a preview when clicking on a node?
- [ ] If possible, being able to open a `note` from the Web launching the default editor for `markdown` content
- [x] Embed the default template
- [x] Add an option to chose another `tmpl` file to make this code more versatile
- [ ] Tests (zero tests for now!)
- [ ] Open the browser with the correct port when launching the webserver (+option to deactivate this behaviour)

# Typical commands

## Installation

for now:
```shell
go install github.com/bligneri/zk-graph/cmd/zk-graph
```

this will give you the `zk-graph` binary and you can now start using it.


## Generating a file in two steps

```shell
zk graph -t daily --format=json > /tmp/zk-graph/my_notes.json
zk-graph -in /tmp/zk-graph/my_notes.json
```

You can now see the graph on your browser

## Generating a file with a pipe

```shell
zk graph -t daily --format=json | zk-graph -in -
```

=> It will pipe the outcome of the `zk graph` command directly to the `zk-graph` utility and this file will be server by the webserver

## Launching the webserver

The Webserver is watching the content of `/tmp/zk-graph` directory and will serve a file directly (if there is a single file)
or show the directory content (if there are multiple files)_
```shell
./zk-graph --server
```

The webserver will fail if there is no `*.html` to serve
