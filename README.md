# go2dot

A utility to generate a Grapviz dot graph for a golang package structs.

## Install

Please use a later [version](https://gitlab.com/graphviz/graphviz/-/releases) of Graphviz/dot (>12) 

```shell
make

export PATH=/path/to/go2dot/bin:$PATH
```

## Usage

```shell
# Generate DOT file for a package (including private functions and structs)
go2dot dot -p ./pkg/golang/ > samples/graph.dot 

# Generate PNG
dot -Tpng samples/graph.dot -o samples/graph.png
```