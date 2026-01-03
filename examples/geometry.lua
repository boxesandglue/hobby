-- Geometry helpers example: midpoint, distance, between

local h = require("hobby")

local p1 = h.point(0, 0)
local p2 = h.point(100, 60)

-- Midpoint
local mid = h.midpoint(p1, p2)
print(string.format("Midpoint of %s and %s: %s", tostring(p1), tostring(p2), tostring(mid)))

-- Distance
local dist = h.distance(p1, p2)
print(string.format("Distance: %.2f", dist))

-- Between (interpolation) - equivalent to MetaPost's t[a,b]
local quarter = h.between(p1, p2, 0.25)  -- 0.25[p1,p2]
local third = h.between(p1, p2, 0.333)   -- 1/3[p1,p2]
local half = h.between(p1, p2, 0.5)      -- 0.5[p1,p2] = midpoint
local threequarter = h.between(p1, p2, 0.75)  -- 0.75[p1,p2]

print(string.format("0.25[p1,p2]: %s", tostring(quarter)))
print(string.format("0.5[p1,p2]: %s", tostring(half)))
print(string.format("0.75[p1,p2]: %s", tostring(threequarter)))

-- Draw the line
local line = h.path()
    :moveto(p1)
    :lineto(p2)
    :stroke("gray")
    :build()

-- Mark the points
local mark1 = h.fullcircle():scaled(6):shifted(p1.x, p1.y):fill("black")
local mark2 = h.fullcircle():scaled(6):shifted(p2.x, p2.y):fill("black")
local markMid = h.fullcircle():scaled(6):shifted(mid.x, mid.y):fill("red")
local markQuarter = h.fullcircle():scaled(4):shifted(quarter.x, quarter.y):fill("blue")
local markThreeQuarter = h.fullcircle():scaled(4):shifted(threequarter.x, threequarter.y):fill("green")

-- Practical example: construct a triangle with midpoints
local a = h.point(150, 0)
local b = h.point(250, 0)
local c = h.point(200, 70)

local triangle = h.path()
    :moveto(a):lineto(b):lineto(c):cycle()
    :stroke("navy")
    :build()

-- Midpoints of each side
local mab = h.midpoint(a, b)
local mbc = h.midpoint(b, c)
local mca = h.midpoint(c, a)

-- Medial triangle (connecting midpoints)
local medial = h.path()
    :moveto(mab):lineto(mbc):lineto(mca):cycle()
    :stroke("orange")
    :build()

-- Mark midpoints
local markMab = h.fullcircle():scaled(4):shifted(mab.x, mab.y):fill("orange")
local markMbc = h.fullcircle():scaled(4):shifted(mbc.x, mbc.y):fill("orange")
local markMca = h.fullcircle():scaled(4):shifted(mca.x, mca.y):fill("orange")

-- Centroid (intersection of medians) = 1/3 from each vertex to opposite midpoint
-- Or simply the average of all three vertices
local centroid = h.between(a, h.midpoint(b, c), 0.6667)
local markCentroid = h.fullcircle():scaled(5):shifted(centroid.x, centroid.y):fill("red")

-- Output to SVG
local svg = h.svg()
    :padding(10)
    :add(line)
    :add(mark1)
    :add(mark2)
    :add(markMid)
    :add(markQuarter)
    :add(markThreeQuarter)
    :add(triangle)
    :add(medial)
    :add(markMab)
    :add(markMbc)
    :add(markMca)
    :add(markCentroid)
    :write("geometry.svg")

print("Created geometry.svg")
