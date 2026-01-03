-- BuildCycle and control points example

local h = require("hobby")

-- Example 1: Show control points of a curve
local curve = h.path()
    :moveto(h.point(0, 0))
    :curveto(h.point(50, 60))
    :curveto(h.point(100, 0))
    :stroke("blue")
    :build()

-- Get control points at t=1 (the middle knot)
local pt = curve:pointat(1)
local pre = curve:precontrol(1)   -- incoming control point
local post = curve:postcontrol(1) -- outgoing control point

print(string.format("Point at t=1: %s", tostring(pt)))
print(string.format("Precontrol at t=1: %s", tostring(pre)))
print(string.format("Postcontrol at t=1: %s", tostring(post)))

-- Draw control point handles
local handle1 = h.path()
    :moveto(pre)
    :lineto(pt)
    :stroke("red")
    :evenly()
    :build()

local handle2 = h.path()
    :moveto(pt)
    :lineto(post)
    :stroke("green")
    :evenly()
    :build()

-- Mark the points
local ptMark = h.fullcircle():scaled(4):shifted(pt.x, pt.y):fill("blue")
local preMark = h.fullcircle():scaled(4):shifted(pre.x, pre.y):fill("red")
local postMark = h.fullcircle():scaled(4):shifted(post.x, post.y):fill("green")

local svg1 = h.svg()
    :padding(10)
    :add(curve)
    :add(handle1)
    :add(handle2)
    :add(ptMark)
    :add(preMark)
    :add(postMark)
    :write("controlpoints.svg")

print("Created controlpoints.svg")

-- Example 2: BuildCycle - create a filled region from intersecting paths
-- Create four lines forming a diamond shape
local top = h.point(150, 80)
local right = h.point(200, 40)
local bottom = h.point(150, 0)
local left = h.point(100, 40)

local line1 = h.path():moveto(left):lineto(top):stroke("gray"):build()
local line2 = h.path():moveto(top):lineto(right):stroke("gray"):build()
local line3 = h.path():moveto(right):lineto(bottom):stroke("gray"):build()
local line4 = h.path():moveto(bottom):lineto(left):stroke("gray"):build()

-- Build cycle from the four lines
local diamond = h.buildcycle(line1, line2, line3, line4)
if diamond then
    diamond:fill(h.rgb(0.8, 0.9, 1.0)):stroke("navy")
    print("BuildCycle succeeded!")
else
    print("BuildCycle returned nil")
end

-- Example 3: BuildCycle with curves - lens shape
local arc1 = h.path()
    :moveto(h.point(250, 40))
    :curveto(h.point(300, 80))
    :curveto(h.point(350, 40))
    :stroke("gray")
    :build()

local arc2 = h.path()
    :moveto(h.point(350, 40))
    :curveto(h.point(300, 0))
    :curveto(h.point(250, 40))
    :stroke("gray")
    :build()

local lens = h.buildcycle(arc1, arc2)
if lens then
    lens:fill(h.rgb(1.0, 0.9, 0.8)):stroke("orange")
    print("Lens BuildCycle succeeded!")
end

local svg2 = h.svg()
    :padding(10)
    :add(line1)
    :add(line2)
    :add(line3)
    :add(line4)
if diamond then svg2:add(diamond) end
svg2:add(arc1)
    :add(arc2)
if lens then svg2:add(lens) end
svg2:write("buildcycle.svg")

print("Created buildcycle.svg")
