-- Advanced transformations example

local h = require("hobby")

-- Original square for reference
local origin = h.point(50, 50)
local square = h.unitsquare():scaled(30):shifted(35, 35):stroke("gray")

-- zscaled: complex multiplication (scale + rotate)
-- zscaled(1, 1) = scale by sqrt(2), rotate 45 degrees
local zscaled1 = h.unitsquare()
    :scaled(30)
    :shifted(35, 35)
    :zscaled(0.707, 0.707)  -- rotate 45 deg, keep size
    :stroke("blue")

-- zscaled with point
local zscaled2 = h.unitsquare()
    :scaled(20)
    :zscaled(h.dir(30))  -- rotate 30 degrees using dir
    :shifted(120, 50)
    :stroke("red")

-- rotatedaround: rotate around a specific point
local center = h.point(200, 50)
local rotated1 = h.unitsquare()
    :scaled(20)
    :shifted(190, 40)
    :stroke("lightgray")

local rotated2 = h.unitsquare()
    :scaled(20)
    :shifted(190, 40)
    :rotatedaround(center, 45)
    :stroke("green")

-- Mark the center of rotation
local centerMark = h.fullcircle():scaled(4):shifted(200, 50):fill("green")

-- reflectedabout: mirror across a line
local original = h.path()
    :moveto(h.point(250, 30))
    :curveto(h.point(280, 70))
    :curveto(h.point(300, 40))
    :stroke("purple")
    :build()

-- Mirror line: vertical at x=310
local p1 = h.point(310, 0)
local p2 = h.point(310, 100)

local reflected = h.path()
    :moveto(h.point(250, 30))
    :curveto(h.point(280, 70))
    :curveto(h.point(300, 40))
    :build()
    :reflectedabout(p1, p2)
    :stroke("orange")

-- Draw the mirror line
local mirrorLine = h.path()
    :moveto(h.point(310, 20))
    :lineto(h.point(310, 80))
    :evenly()
    :stroke("gray")
    :build()

-- Output to SVG
local svg = h.svg()
    :padding(10)
    :add(square)
    :add(zscaled1)
    :add(zscaled2)
    :add(rotated1)
    :add(rotated2)
    :add(centerMark)
    :add(original)
    :add(reflected)
    :add(mirrorLine)
    :write("transforms2.svg")

print("Created transforms2.svg")
