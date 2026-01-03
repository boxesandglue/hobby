-- Path operations example: arclength, arctime, intersection

local h = require("hobby")

-- Create a curved path
local curve = h.path()
    :moveto(h.point(0, 0))
    :curveto(h.point(50, 80))
    :curveto(h.point(100, 0))
    :stroke("blue")
    :build()

-- Get arc length
local len = curve.arclength
print(string.format("Arc length: %.2f", len))

-- Get point at half the arc length using arctime
local halfLen = len / 2
local t = curve:arctime(halfLen)
print(string.format("Time at half arc length: %.3f", t))

local midPoint = curve:pointat(t)
print(string.format("Mid point (by arc): %s", tostring(midPoint)))

-- Mark the midpoint
local midMark = h.fullcircle():scaled(4):shifted(midPoint.x, midPoint.y):fill("red")

-- Create two intersecting paths
local path1 = h.path()
    :moveto(h.point(150, 0))
    :lineto(h.point(250, 80))
    :stroke("green")
    :build()

local path2 = h.path()
    :moveto(h.point(150, 80))
    :lineto(h.point(250, 0))
    :stroke("purple")
    :build()

-- Find intersection
local t1, t2 = path1:intersectiontimes(path2)
print(string.format("Intersection times: t1=%.3f, t2=%.3f", t1, t2))

local intersect = path1:intersectionpoint(path2)
if intersect then
    print(string.format("Intersection point: %s", tostring(intersect)))
    -- Mark the intersection
    local intersectMark = h.fullcircle():scaled(6):shifted(intersect.x, intersect.y):fill("orange")

    -- Output to SVG
    local svg = h.svg()
        :padding(10)
        :add(curve)
        :add(midMark)
        :add(path1)
        :add(path2)
        :add(intersectMark)
        :write("pathops.svg")
else
    print("No intersection found")
    local svg = h.svg()
        :padding(10)
        :add(curve)
        :add(midMark)
        :add(path1)
        :add(path2)
        :write("pathops.svg")
end

print("Created pathops.svg")

-- Additional example: mark points at regular arc intervals
print("\nPoints at regular arc intervals:")
local circle = h.fullcircle():scaled(60):shifted(350, 40)
local circleLen = circle.arclength
print(string.format("Circle arc length: %.2f", circleLen))

local svg2 = h.svg():padding(10):add(circle:stroke("gray"))

for i = 0, 7 do
    local arcPos = (i / 8) * circleLen
    local ti = circle:arctime(arcPos)
    local pt = circle:pointat(ti)
    local mark = h.fullcircle():scaled(4):shifted(pt.x, pt.y):fill("blue")
    svg2:add(mark)
    print(string.format("  i=%d: arc=%.2f, t=%.3f, point=%s", i, arcPos, ti, tostring(pt)))
end

svg2:write("arclength_points.svg")
print("Created arclength_points.svg")
