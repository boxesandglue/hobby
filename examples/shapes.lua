-- Shapes and styles example

local h = require("hobby")

-- Predefined shapes
local full = h.fullcircle():scaled(30):shifted(20, 80):stroke("black")
local half = h.halfcircle():scaled(30):shifted(70, 80):stroke("blue")
local quarter = h.quartercircle():scaled(30):shifted(120, 80):stroke("red")
local square = h.unitsquare():scaled(30):shifted(150, 65):stroke("green")

-- Using dir() for direction vectors
local right = h.dir(0)    -- (1, 0)
local up = h.dir(90)      -- (0, 1)
local diag = h.dir(45)    -- (0.707, 0.707)

-- Arrow using dir for positioning
local start = h.point(0, 40)
local arrow1 = h.path()
    :moveto(start)
    :lineto(start + right * 60)
    :arrow()
    :stroke("black")
    :build()

-- Line cap styles
local line_butt = h.path()
    :moveto(h.point(0, 0))
    :lineto(h.point(50, 0))
    :linecap("butt")
    :strokewidth(6)
    :stroke("navy")
    :build()

local line_round = h.path()
    :moveto(h.point(60, 0))
    :lineto(h.point(110, 0))
    :linecap("round")
    :strokewidth(6)
    :stroke("navy")
    :build()

local line_square = h.path()
    :moveto(h.point(120, 0))
    :lineto(h.point(170, 0))
    :linecap("square")
    :strokewidth(6)
    :stroke("navy")
    :build()

-- Line join styles
local join_miter = h.path()
    :moveto(h.point(0, 20))
    :lineto(h.point(20, 35))
    :lineto(h.point(40, 20))
    :linejoin("miter")
    :strokewidth(4)
    :stroke("darkred")
    :build()

local join_round = h.path()
    :moveto(h.point(60, 20))
    :lineto(h.point(80, 35))
    :lineto(h.point(100, 20))
    :linejoin("round")
    :strokewidth(4)
    :stroke("darkred")
    :build()

local join_bevel = h.path()
    :moveto(h.point(120, 20))
    :lineto(h.point(140, 35))
    :lineto(h.point(160, 20))
    :linejoin("bevel")
    :strokewidth(4)
    :stroke("darkred")
    :build()

-- Output to SVG
local svg = h.svg()
    :padding(10)
    :add(full)
    :add(half)
    :add(quarter)
    :add(square)
    :add(arrow1)
    :add(line_butt)
    :add(line_round)
    :add(line_square)
    :add(join_miter)
    :add(join_round)
    :add(join_bevel)
    :write("shapes.svg")

print("Created shapes.svg")

-- Print dir() values
print(string.format("dir(0) = %s", tostring(right)))
print(string.format("dir(90) = %s", tostring(up)))
print(string.format("dir(45) = %s", tostring(diag)))
