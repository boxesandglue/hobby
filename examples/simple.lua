-- Simple example: draw a curved path

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
local svg = h.svg()
    :padding(5)
    :add(path)
    :write("simple.svg")

print("Created simple.svg")

-- Also demonstrate predefined paths (already solved, no build() needed)
local circle = h.fullcircle():scaled(50):shifted(100, 50)
local square = h.unitsquare():scaled(30):shifted(160, 35)

local svg2 = h.svg()
    :padding(5)
    :add(circle)
    :add(square)
    :write("shapes.svg")

print("Created shapes.svg")
