-- Color example: demonstrate stroke and fill colors

local h = require("hobby")

-- Create a curved path with red stroke
local curve = h.path()
    :moveto(h.point(0, 0))
    :curveto(h.point(50, 80))
    :curveto(h.point(100, 0))
    :stroke("red")
    :strokewidth(2)
    :build()

-- Create a filled circle (blue fill, no stroke)
local circle = h.fullcircle()
    :scaled(40)
    :shifted(50, 50)

-- Create a square with green fill and black stroke
local square = h.path()
    :moveto(h.point(120, 10))
    :lineto(h.point(170, 10))
    :lineto(h.point(170, 60))
    :lineto(h.point(120, 60))
    :cycle()
    :fill(h.rgb(0.2, 0.8, 0.3))  -- RGB values 0-1
    :stroke("#000000")           -- CSS hex color
    :strokewidth(1.5)
    :build()

-- Create path with grayscale
local gray_path = h.path()
    :moveto(h.point(0, -20))
    :curveto(h.point(100, -20))
    :stroke(h.gray(0.5))
    :strokewidth(3)
    :build()

-- Output to SVG
local svg = h.svg()
    :padding(10)
    :add(curve)
    :add(circle)
    :add(square)
    :add(gray_path)
    :write("colors.svg")

print("Created colors.svg")
