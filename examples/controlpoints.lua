-- Control points visualization
-- Draws a fullcircle and shows its 8 knots with their control points

local h = require("hobby")

local cx, cy = 50, 50   -- center
local scale = 80         -- diameter

-- The main circle
local circle = h.fullcircle():scaled(scale):shifted(cx, cy):stroke("black")

-- Helpers
local function dot(x, y, color, size)
    return h.fullcircle():scaled(size or 4):shifted(x, y):fill(color)
end

local function line(x1, y1, x2, y2, color)
    return h.path()
        :moveto(h.point(x1, y1))
        :lineto(h.point(x2, y2))
        :stroke(color)
        :strokewidth(0.3)
        :build()
end

local svg = h.svg()
    :padding(8)
    :add(circle)

-- Draw knots, control points, and handles for each of the 8 segments
local n = circle.length  -- 8 for fullcircle
for i = 0, n - 1 do
    local pt = circle:pointat(i)
    local post = circle:postcontrol(i)      -- right control point
    local pre = circle:precontrol(i)        -- left control point

    -- Control handle lines (knot -> control point)
    svg:add(line(pt.x, pt.y, post.x, post.y, "cornflowerblue"))
    svg:add(line(pt.x, pt.y, pre.x, pre.y, "tomato"))

    -- Control points as small dots
    svg:add(dot(post.x, post.y, "cornflowerblue", 2.5))
    svg:add(dot(pre.x, pre.y, "tomato", 2.5))

    -- Knot on the circle
    svg:add(dot(pt.x, pt.y, "black"))
end

svg:write("controlpoints.svg")
print("Created controlpoints.svg")
