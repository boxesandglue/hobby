-- Pen and dash pattern example

local h = require("hobby")

-- Line with default pen (pencircle)
local line1 = h.path()
    :moveto(h.point(0, 0))
    :lineto(h.point(100, 0))
    :stroke("black")
    :build()

-- Line with thick pencircle
local line2 = h.path()
    :moveto(h.point(0, 20))
    :lineto(h.point(100, 20))
    :pen(h.pencircle(4))
    :stroke("blue")
    :build()

-- Line with pensquare
local line3 = h.path()
    :moveto(h.point(0, 40))
    :lineto(h.point(100, 40))
    :pen(h.pensquare(3))
    :stroke("red")
    :build()

-- Dashed line (evenly)
local line4 = h.path()
    :moveto(h.point(0, 60))
    :lineto(h.point(100, 60))
    :evenly()
    :stroke("green")
    :build()

-- Dotted line (withdots)
local line5 = h.path()
    :moveto(h.point(0, 80))
    :lineto(h.point(100, 80))
    :withdots()
    :stroke("purple")
    :build()

-- Custom dash pattern (long dash, short gap)
local line6 = h.path()
    :moveto(h.point(0, 100))
    :lineto(h.point(100, 100))
    :dash(h.dashed(8, 2))
    :stroke("orange")
    :build()

-- Scaled dash pattern
local line7 = h.path()
    :moveto(h.point(0, 120))
    :lineto(h.point(100, 120))
    :dash(h.evenly():scaled(2))
    :stroke("brown")
    :build()

-- Curve with pen and dash
local curve = h.path()
    :moveto(h.point(120, 0))
    :curveto(h.point(170, 60))
    :curveto(h.point(220, 0))
    :pen(h.pencircle(2))
    :evenly()
    :stroke("navy")
    :build()

-- Output to SVG
local svg = h.svg()
    :padding(10)
    :add(line1)
    :add(line2)
    :add(line3)
    :add(line4)
    :add(line5)
    :add(line6)
    :add(line7)
    :add(curve)
    :write("pens.svg")

print("Created pens.svg")
