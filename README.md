# Chcat

Chcat is a simple filter that reads the standard input or a file, invokes a text editor and writes the result to standard output.

It is used to edit text files, usually templates of some kind, to change some values for a single operation
but preserve the original.

## Examples

Post a json template to a rest endpoint
```
chcat request.json | curl -d @- http://host/operations
```

Render a template
```
chcat data.yml | mustache - template.mustache
```

## Installation

Chcat is tested with go 1.18 on debian linux.

```
go install github.com/anastasop/chcat@latest
```

It can be implemented easily as an oneline script `cat && vi && cat` but using Go makes distribution, installation and evolution easier.

## Bugs

It cannot be used with process substitution like `cat <(chcat file)`. Most shells close the stdin for process substitution.
Failed with bash and ksh, worked only with [plan9 rc](https://9fans.github.io/plan9port/).

## License

Chcat is licensed under the [GPL](https://www.gnu.org/licenses/gpl-3.0.en.html).
