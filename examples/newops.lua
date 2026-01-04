-- newops.lua: Example demonstrating new path operations
local h = require("hobby")

-- Create a curved path
local curve = h.path()
    :moveto(h.point(0, 0))
    :curveto(h.point(50, 80))
    :curveto(h.point(100, 0))
    :stroke(h.color("black"))
    :build()

-- Find where the tangent is horizontal (pointing right)
local t_horiz = curve:directiontime(1, 0)
print(string.format("Horizontal tangent at t = %.3f", t_horiz))

-- Find the point at that time
local pt_horiz = curve:directionpoint(1, 0)
if pt_horiz then
    print(string.format("Point: (%.1f, %.1f)", pt_horiz.x, pt_horiz.y))
end

-- Create two intersecting paths for cutbefore/cutafter demo
local path1 = h.path()
    :moveto(h.point(0, 50))
    :lineto(h.point(100, 50))
    :stroke(h.color("blue"))
    :build()

local path2 = h.path()
    :moveto(h.point(50, 0))
    :lineto(h.point(50, 100))
    :stroke(h.color("red"))
    :build()

-- Cut path1 before and after intersection with path2
local before = path1:cutafter(path2)  -- from start to intersection
local after = path1:cutbefore(path2)  -- from intersection to end

print(string.format("path1 length: %d segments", path1.length))
print(string.format("cutafter length: %d segments", before.length))
print(string.format("cutbefore length: %d segments", after.length))

-- Create SVG with all paths
local svg = h.svg()
    :padding(10)
    :add(curve)
    :add(path1:evenly())
    :add(path2:evenly())

-- Mark the horizontal tangent point
if pt_horiz then
    local dot = h.fullcircle():scaled(6):shifted(pt_horiz.x, pt_horiz.y):fill(h.color("green"))
    svg:add(dot)
end

svg:write("newops.svg")
print("Created newops.svg")
