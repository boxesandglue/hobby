[![Homepage](https://img.shields.io/badge/homepage-boxesandglue.dev-blue)](https://boxesandglue.dev/hobby)

# Hobby

Lua bindings for the [mpgo](https://github.com/boxesandglue/mpgo) MetaPost curve library. Named after John Hobby, the developer of MetaPost and the Hobby-Knuth smooth curve algorithm.

## Installation

### Homebrew (macOS/Linux)

```bash
brew install boxesandglue/tap/hobby
```

### Go

```bash
go install github.com/boxesandglue/hobby/cmd/hobby@latest
```

### Binary Downloads

Download binaries from the [releases page](https://github.com/boxesandglue/hobby/releases).

## Usage

```bash
hobby script.lua
```

## Examples

### Curved Path

A smooth curve through five points using the Hobby algorithm:

```lua
local h = require("hobby")

-- Create some points
local z0 = h.point(0, 0)
local z1 = h.point(60, 40)
local z2 = h.point(40, 90)
local z3 = h.point(10, 70)
local z4 = h.point(30, 50)

-- Build a closed curved path
local path = h.path()
    :moveto(z0)
    :curveto(z1)
    :curveto(z2)
    :curveto(z3)
    :curveto(z4)
    :cycle()
    :build()

-- Create SVG and add the path
h.svg()
    :padding(5)
    :add(path)
    :write("simple.svg")
```

![Curved path](doc/simple.svg)

### Transformations

Predefined shapes with various transformations:

```lua
local h = require("hobby")

-- Original circle
local circle = h.fullcircle():scaled(20):shifted(30, 80):stroke("black")

-- xscaled (stretched horizontally)
local xscaled = h.fullcircle():scaled(20):xscaled(2):shifted(90, 80):stroke("blue")

-- yscaled (stretched vertically)
local yscaled = h.fullcircle():scaled(20):yscaled(2):shifted(150, 80):stroke("red")

-- slanted (sheared)
local slanted = h.unitsquare():scaled(30):slanted(0.5):shifted(30, 20):stroke("green")

-- Combination
local combo = h.fullcircle()
    :scaled(20)
    :xscaled(1.5)
    :yscaled(0.7)
    :slanted(0.3)
    :shifted(130, 30)
    :stroke("purple")

-- Output to SVG
h.svg()
    :padding(10)
    :add(circle)
    :add(xscaled)
    :add(yscaled)
    :add(slanted)
    :add(combo)
    :write("transforms.svg")
```

![Transformations](doc/transforms.svg)

## Documentation

Full documentation is available at [boxesandglue.dev/hobby](https://boxesandglue.dev/hobby).

## License

BSD-3-Clause. See [LICENSE](LICENSE).
