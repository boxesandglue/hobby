-- Context equation solver example
-- Demonstrates solving geometric constraints

local h = require("hobby")

-- Create a context for equation solving
local ctx = h.context()

-- Define known points (corners of a rectangle)
local z0 = ctx:known(0, 0)      -- bottom-left
local z1 = ctx:known(100, 0)    -- bottom-right
local z2 = ctx:known(100, 60)   -- top-right
local z3 = ctx:known(0, 60)     -- top-left

-- Create derived points using constraints
local mid_bottom = ctx:midpointof(z0, z1)   -- midpoint of bottom edge
local mid_top = ctx:midpointof(z3, z2)      -- midpoint of top edge
local center = ctx:midpointof(mid_bottom, mid_top)  -- center of rectangle

-- Point at 1/3 along bottom edge
local third = ctx:betweenat(z0, z1, 1/3)

-- Intersection of diagonals
local diag_center = ctx:intersectionof(z0, z2, z1, z3)

-- Solve the system
ctx:solve()

-- Print computed values
print(string.format("Center: (%.1f, %.1f)", center.x, center.y))
print(string.format("Diagonal intersection: (%.1f, %.1f)", diag_center.x, diag_center.y))
print(string.format("1/3 point: (%.1f, %.1f)", third.x, third.y))

-- Draw the rectangle and computed points
local rect = h.path()
    :moveto(z0:point())
    :lineto(z1:point())
    :lineto(z2:point())
    :lineto(z3:point())
    :cycle()
    :stroke("black")
    :build()

-- Diagonals (dashed)
local diag1 = h.path()
    :moveto(z0:point())
    :lineto(z2:point())
    :evenly()
    :stroke("gray")
    :build()

local diag2 = h.path()
    :moveto(z1:point())
    :lineto(z3:point())
    :evenly()
    :stroke("gray")
    :build()

-- Mark the center point
local centerMark = h.fullcircle():scaled(6):shifted(center.x, center.y):fill("red")

-- Mark the third point
local thirdMark = h.fullcircle():scaled(6):shifted(third.x, third.y):fill("blue")

-- Output to SVG
h.svg()
    :padding(10)
    :add(rect)
    :add(diag1)
    :add(diag2)
    :add(centerMark)
    :add(thirdMark)
    :write("context.svg")

print("Created context.svg")

-- More complex example: Triangle with medians
print("\n-- Triangle with medians --")

local ctx2 = h.context()

-- Triangle vertices
local a = ctx2:known(0, 0)
local b = ctx2:known(80, 0)
local c = ctx2:known(40, 70)

-- Midpoints of sides (computed as known points for intersection)
local mAB = ctx2:known(40, 0)     -- midpoint of a-b
local mBC = ctx2:known(60, 35)    -- midpoint of b-c
local mCA = ctx2:known(20, 35)    -- midpoint of c-a

-- Centroid (intersection of medians: a-mBC and b-mCA)
local centroid = ctx2:intersectionof(a, mBC, b, mCA)

ctx2:solve()

print(string.format("Centroid: (%.1f, %.1f)", centroid.x, centroid.y))
print(string.format("Expected: (%.1f, %.1f)", (0+80+40)/3, (0+0+70)/3))

-- Draw triangle
local triangle = h.path()
    :moveto(a:point())
    :lineto(b:point())
    :lineto(c:point())
    :cycle()
    :stroke("black")
    :build()

-- Draw medians
local median1 = h.path()
    :moveto(a:point())
    :lineto(mBC:point())
    :stroke("blue")
    :build()

local median2 = h.path()
    :moveto(b:point())
    :lineto(mCA:point())
    :stroke("blue")
    :build()

local median3 = h.path()
    :moveto(c:point())
    :lineto(mAB:point())
    :stroke("blue")
    :build()

-- Mark centroid
local centroidMark = h.fullcircle():scaled(8):shifted(centroid.x, centroid.y):fill("red")

h.svg()
    :padding(10)
    :add(triangle)
    :add(median1)
    :add(median2)
    :add(median3)
    :add(centroidMark)
    :write("triangle.svg")

print("Created triangle.svg")

-- Example using ctx:path() with variable references
print("\n-- Path with variable references --")

local ctx3 = h.context()

local p0 = ctx3:known(0, 0)
local p1 = ctx3:known(100, 0)
local p2 = ctx3:unknown()
local p3 = ctx3:unknown()

-- p2 is at the midpoint of p0-p1, shifted up by 50
local offset = ctx3:known(0, 50)
local mid = ctx3:midpointof(p0, p1)
ctx3:sum(p2, mid, offset)

-- p3 is 1/4 along p0-p1, shifted up by 30
local quarter = ctx3:betweenat(p0, p1, 0.25)
local offset2 = ctx3:known(0, 30)
ctx3:sum(p3, quarter, offset2)

ctx3:solve()

print(string.format("p2: (%.1f, %.1f)", p2.x, p2.y))
print(string.format("p3: (%.1f, %.1f)", p3.x, p3.y))

-- Build path using variable references
local path = ctx3:path()
    :movetovar(p0)
    :curvetovar(p3)
    :curvetovar(p2)
    :linetovar(p1)
    :stroke("purple")
    :strokewidth(2)
    :build()

h.svg()
    :padding(10)
    :add(path)
    :write("varpath.svg")

print("Created varpath.svg")
