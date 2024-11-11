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

## Future

- [ ] Add more parameters to the template in particular the size of the `d3` canvas, the `strength` of the various forces. Most certainly in a configuration file in the `.zk` directory?
- [ ] Keep the first `200` charcters of the content of each note to offer a preview when clicking on a node?
- [ ] If possible, being able to open a `note` from the Web launching the default editor for `markdown` content
- [ ] Add an option to chose another `tmpl` file to make this code more versatile
- [ ] Tests (zero tests for now!)
- [ ] Open the browser with the correct port when launching the webserver (+option to deactivate this behaviour)

# Typical commands

## Launching the webserver
The Webserver is watching the content of `out` and will serve a file directly (if there is a single file) 
or show the directory content (if there are multiple files)_
```shell
./zk-graph --server
```

You should launch the server and point a browser to this server

## Generating a file in two steps

```shell
zk graph -t daily --format=json > my_notes.json
zk-graph -json_file my_notes.json
```

You can now see the graph on your browser

## Generating a file with a pipe

```shell
zk graph -t daily --format=json | zk-graph -json_file my_notes.json
```

=> It will pipe the outcome of the `zk graph` command directly to the `zk-graph` utility and this file will be server by the webserver
