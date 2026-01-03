-- Clipping example

local h = require("hobby")

-- Create a clipping path (a circle)
-- Predefined paths (fullcircle, unitsquare, etc.) are already solved
local clipCircle = h.fullcircle()
    :scaled(60)
    :shifted(50, 50)

-- Create some paths to be clipped
local line1 = h.path()
    :moveto(h.point(0, 0))
    :lineto(h.point(100, 100))
    :stroke("red")
    :strokewidth(3)
    :build()

local line2 = h.path()
    :moveto(h.point(0, 100))
    :lineto(h.point(100, 0))
    :stroke("blue")
    :strokewidth(3)
    :build()

-- Predefined paths don't need :build()
local square = h.unitsquare()
    :scaled(80)
    :shifted(10, 10)
    :stroke("green")
    :strokewidth(2)

-- Create a picture and add paths
local pic = h.picture()
    :add(line1)
    :add(line2)
    :add(square)
    :clip(clipCircle)

-- Draw the clip boundary as reference (dashed)
local clipOutline = h.fullcircle()
    :scaled(60)
    :shifted(50, 50)
    :evenly()
    :stroke("gray")

-- Output to SVG
local svg = h.svg()
    :padding(10)
    :addpicture(pic)
    :add(clipOutline)
    :write("clip.svg")

print("Created clip.svg")
