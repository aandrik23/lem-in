# lem-in

A Go implementation of the classic **Lem-in** ant farm path-finding and simulation challenge, with an optional Python-based visualizer.

## Features

* Parses an input file describing:

  * Number of ants
  * Definition of rooms with coordinates
  * Tunnels connecting rooms
  * Special `##start` and `##end` markers
* Finds all simple paths from start to end using DFS.
* Selects a set of disjoint shortest paths, maximizing throughput.
* Computes an optimal assignment of ants to each path to minimize total turns.
* Simulates ant movements, respecting occupancy rules:

  * Only one ant per non-`start`/non-`end` room per turn
  * Each tunnel used at most once per turn
* Outputs moves in the form `L<antID>-<roomName>` per turn.
* `--visualize` (`-v`) flag generates a `simulation.json` and invokes a Python visualizer.

## Installation

1. **Prerequisites**

   * Go 1.18+
   * Python 3.7+ (for visualization)
   * `networkx` and `matplotlib` Python packages:

     ```bash
     pip install networkx matplotlib
     ```

2. **Clone and build**

   ```bash
   git clone https://github.com/youruser/lem-in.git
   cd lem-in
   go build -o lemin main.go
   ```

## Usage

```bash
./lemin [options] <input-file>
```

### Options

* `-h`, `--help`
  Provide usage options.

* `-v`, `--visualize`
  Generate a `simulation.json` and run the Python visualizer (see below).

* `--debug`
  Provide debug logs.

### Examples

Basic run:

```bash
./lemin example1.txt
```

Run with visualization:

```bash
./lemin -v example1.txt
```

Run help:
```bash
./lemin -h
```

Run with debug logs:
```bash
./lemin --debug [-v] example1.txt
```

## Input Format

```
<ants>
##start
<name> <x> <y>
##end
<name> <x> <y>
<roomName> <x> <y>
...
<roomA>-<roomB>
...
```

* Room names must not start with `L` or `#` and contain no spaces.
* Coordinates are integers.

## Output

Each turn prints moves of the form:

```
L1-roomA L2-roomB …
```

## Visualization

When run with `-v`, the program creates `simulation.json` containing both:

* `rooms`: list of `{ name, x, y }`
* `moves`: list of `{ turn, ant, from, to }`

Then it calls the Python script:

```bash
python3 visualize.py --input simulation.json
```

which animates the ant movements using Matplotlib and NetworkX.

## Testing

Unit tests live alongside each package in `internal/`. Run:

```bash
go test ./internal/...
```

## Creators
aziagaki  
mtzemana  
aadriko  
vparik  
tdiridis

## License

MIT ©
