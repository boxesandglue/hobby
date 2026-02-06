-- Smooth curve: same points, different methods
-- Shows the difference between straight lines and Hobby curves

local h = require("hobby")

local points = {
    h.point(0, 0),
    h.point(30, 50),
    h.point(80, 60),
    h.point(120, 30),
    h.point(160, 50),
    h.point(200, 0),
}

-- Straight lines (gray, dashed)
local lines = h.path():moveto(points[1])
for i = 2, #points do
    lines = lines:lineto(points[i])
end
lines = lines:stroke("gray"):evenly():strokewidth(0.5):build()

-- Smooth curve (blue)
local curve = h.path():moveto(points[1])
for i = 2, #points do
    curve = curve:curveto(points[i])
end
curve = curve:stroke("steelblue"):strokewidth(1):build()

-- Dot markers
local pic = h.picture()
pic:add(lines)
pic:add(curve)

local gray = h.color("gray")
for i = 1, #points do
    local circle = h.fullcircle()
        :scaled(4)
        :shifted(points[i].x, points[i].y)
        :fill("steelblue")
    pic:add(circle)
end

h.svg()
    :padding(8)
    :addpicture(pic)
    :write("smooth.svg")

print("Created smooth.svg")
