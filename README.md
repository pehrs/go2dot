# go2dot

A utility to generate a [Graphviz](https://graphviz.org/) dot graph for a golang package structures and functions.

Example:

![graph example](./samples/graph.png)

## Install

Please use [resent version](https://gitlab.com/graphviz/graphviz/-/releases) of dot/graphviz.

```shell
brew install graphviz

git clone git@github.com:pehrs/go2dot.git
cd ./go2dot¡
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