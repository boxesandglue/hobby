-- Arrow example

local h = require("hobby")

-- Simple arrow (drawarrow equivalent)
local arrow1 = h.path()
    :moveto(h.point(0, 0))
    :lineto(h.point(80, 0))
    :arrow()
    :stroke("black")
    :build()

-- Double arrow (drawdblarrow equivalent)
local arrow2 = h.path()
    :moveto(h.point(0, 20))
    :lineto(h.point(80, 20))
    :dblarrow()
    :stroke("blue")
    :build()

-- Curved arrow
local arrow3 = h.path()
    :moveto(h.point(0, 50))
    :curveto(h.point(40, 80))
    :curveto(h.point(80, 50))
    :arrow()
    :stroke("red")
    :build()

-- Arrow with custom style (longer, narrower)
local arrow4 = h.path()
    :moveto(h.point(0, 90))
    :lineto(h.point(80, 90))
    :arrow()
    :arrowstyle(8, 30)  -- length=8, angle=30
    :stroke("green")
    :build()

-- Arrow with custom style (shorter, wider)
local arrow5 = h.path()
    :moveto(h.point(0, 110))
    :lineto(h.point(80, 110))
    :arrow()
    :arrowstyle(6, 60)  -- length=6, angle=60
    :stroke("purple")
    :build()

-- Thick arrow with pen
local arrow6 = h.path()
    :moveto(h.point(100, 0))
    :curveto(h.point(140, 40))
    :curveto(h.point(180, 0))
    :arrow()
    :pen(h.pencircle(2))
    :stroke("orange")
    :build()

-- Dashed double arrow
local arrow7 = h.path()
    :moveto(h.point(100, 60))
    :lineto(h.point(180, 60))
    :dblarrow()
    :evenly()
    :stroke("navy")
    :build()

-- Output to SVG
local svg = h.svg()
    :padding(10)
    :add(arrow1)
    :add(arrow2)
    :add(arrow3)
    :add(arrow4)
    :add(arrow5)
    :add(arrow6)
    :add(arrow7)
    :write("arrows.svg")

print("Created arrows.svg")
