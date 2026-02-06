-- Equation solver: triangle with centroid
-- The centroid is found as intersection of two medians

local h = require("hobby")

local ctx = h.context()

-- Triangle vertices
local a = ctx:known(0, 0)
local b = ctx:known(120, 0)
local c = ctx:known(50, 90)

-- Midpoints of opposite sides
local mBC = ctx:midpointof(b, c)
local mAC = ctx:midpointof(a, c)

ctx:solve()

-- Centroid: intersection of medians a→mBC and b→mAC
-- (midpoints are now solved, so intersection can use them)
local centroid = ctx:intersectionof(a, mBC, b, mAC)

-- Draw
local pic = h.picture()

-- Triangle
local triangle = h.path()
    :moveto(a:point())
    :lineto(b:point())
    :lineto(c:point())
    :cycle()
    :stroke("black")
    :strokewidth(0.8)
    :build()
pic:add(triangle)

-- Medians (dashed)
local median1 = h.path()
    :moveto(a:point())
    :lineto(mBC:point())
    :stroke("steelblue")
    :strokewidth(0.5)
    :evenly()
    :build()
pic:add(median1)

local median2 = h.path()
    :moveto(b:point())
    :lineto(mAC:point())
    :stroke("steelblue")
    :strokewidth(0.5)
    :evenly()
    :build()
pic:add(median2)

-- Centroid dot
local dot = h.fullcircle()
    :scaled(6)
    :shifted(centroid.x, centroid.y)
    :fill("tomato")
pic:add(dot)

-- Labels
local black = h.color("black")
pic:dotlabel("A", a:point(), "llft", black)
pic:dotlabel("B", b:point(), "lrt", black)
pic:dotlabel("C", c:point(), "top", black)

h.svg()
    :padding(15)
    :addpicture(pic)
    :write("solver.svg")

print("Created solver.svg")
