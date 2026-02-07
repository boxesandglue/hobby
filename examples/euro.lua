-- Euro sign (€) — translated from MetaPost
--
-- Original MetaPost construction:
--   A circle with two slanted bars, clipped to create the € shape.
--   The bars are parallelograms (unitsquare + slant) aligned with the
--   right stroke angle (from o1 through o2).

local h = require("hobby")

-- tand(x) = sin(x)/cos(x) for degrees (MetaPost's sind/cosd)
local function tand(deg)
    return math.sin(math.rad(deg)) / math.cos(math.rad(deg))
end

-- Debug options (set to true to see construction)
local show_lines  = true
local show_dots   = true
local show_labels = true

----------------------------------------------------------------------
-- Parameters
----------------------------------------------------------------------
local x = 10                       -- base unit (MetaPost: box.width/14)
local radius = 5.5 * x
local topbarlength = 10 * x
local thickness = 1 * x

----------------------------------------------------------------------
-- Key points
----------------------------------------------------------------------
local origin = h.point(0, 0)
local o1 = h.point(0, -radius - 0.5 * thickness)       -- bottom of right stroke
local o2 = h.dir(40) * (radius - 0.5 * thickness)       -- 40° on inner circle
local o3 = h.dir(-40) * (radius - 0.5 * thickness)      -- -40° on inner circle
local alpha = (o2 - o1).angle                            -- angle of right stroke

----------------------------------------------------------------------
-- Auxiliary paths for intersections
----------------------------------------------------------------------
-- a1: the right slant (o1 → o2)
local a1 = h.path():moveto(o1):lineto(o2):build()

-- a2: origin → o3 (lower 40° angle)
local a2 = h.path():moveto(origin):lineto(o3):build()

----------------------------------------------------------------------
-- Circle (stroked outline)
----------------------------------------------------------------------
-- MetaPost: draw unitcircle scaled radius withpen pencircle scaled thickness
-- unitcircle = fullcircle scaled 2 → radius 1
-- hobby fullcircle has radius 0.5, so scaled(2*radius) gives radius=radius
local circle = h.fullcircle()
    :scaled(2 * radius)
    :strokewidth(thickness)
    :stroke("black")

----------------------------------------------------------------------
-- Horizontal bars (parallelograms via slant)
----------------------------------------------------------------------
-- MetaPost: unitsquare xscaled topbarlength yscaled thickness slanted (1/tand(alpha))
local hbar = h.unitsquare()
    :xscaled(topbarlength)
    :yscaled(thickness)
    :slanted(1 / tand(alpha))

-- Find c3: where horizontal line y = 0.5x crosses the right slant a1
-- MetaPost: ((-infinity,0.5x) -- (infinity,0.5x)) intersectionpoint a1
local horizLine = h.path()
    :moveto(h.point(-200, 0.5 * x))
    :lineto(h.point(200, 0.5 * x))
    :build()
local c3 = horizLine:intersectionpoint(a1)

-- Bar positions: bars extend leftward from c3
local topbarleft    = h.point(c3.x - topbarlength,  0.5 * x)
local bottombarleft = h.point(c3.x - topbarlength, -0.5 * x - thickness)

-- Filled bars (stroke "none" = MetaPost's fill without draw)
local topbar    = hbar:shifted(topbarleft.x, topbarleft.y):fill("black"):stroke("none")
local bottombar = hbar:shifted(bottombarleft.x, bottombarleft.y):fill("black"):stroke("none")

----------------------------------------------------------------------
-- Clipping path
----------------------------------------------------------------------
-- c2: intersection of right slant (a1) and lower angle (a2)
local c2 = a1:intersectionpoint(a2)

-- c4: a point beyond o2 along the stroke direction
local c4 = o2 + h.dir(alpha) * (2 * thickness)

-- Outer circle boundary (includes half the stroke width)
local outer = h.fullcircle():scaled(2 * (radius + 0.5 * thickness))

-- c1: where the extended slant (o1→c4) exits the outer circle
-- Because hobby iterates circle segments from 0° upward, the upper-right
-- exit point (≈35°) is found before the entry at o1 (270°).
local slantLine = h.path():moveto(o1):lineto(c4):build()
local c1 = slantLine:intersectionpoint(outer)

-- Bounding box of the drawn elements (via picture BBox)
local cp = h.picture():add(circle):add(topbar):add(bottombar)
local ll = cp.llcorner
local ur = cp.urcorner

-- c5: bottom-right clip anchor
local c5 = h.point(o3.x, ll.y)

-- Clip path: keeps left side of circle + bars, cuts the right opening
-- MetaPost: llcorner cp -- c5 -- o3 -- c2 -- c1 -- (xpart c1, top) -- ulcorner -- cycle
local clippath = h.path()
    :moveto(ll)                                  -- lower-left corner
    :lineto(c5)                                  -- bottom edge to c5
    :lineto(o3)                                  -- up to -40° on circle
    :lineto(c2)                                  -- to slant intersection
    :lineto(c1)                                  -- to outer circle exit
    :lineto(h.point(c1.x, ur.y))                -- straight up to top
    :lineto(h.point(ll.x, ur.y))                -- upper-left corner
    :cycle()
    :build()

----------------------------------------------------------------------
-- Assemble and clip
----------------------------------------------------------------------
local pic = h.picture()
    :add(circle)
    :add(topbar)
    :add(bottombar)
    :clip(clippath)

----------------------------------------------------------------------
-- Debug overlays
----------------------------------------------------------------------
local svg = h.svg():padding(5)
svg:addpicture(pic)

if show_lines then
    local inner = h.fullcircle():scaled(2 * (radius - 0.5 * thickness))
        :evenly():stroke("silver"):strokewidth(0.5)
    local outerVis = h.fullcircle():scaled(2 * (radius + 0.5 * thickness))
        :evenly():stroke("silver"):strokewidth(0.5)
    local a1vis = h.path():moveto(o1):lineto(o2):evenly():stroke("silver"):strokewidth(0.5):build()
    local a2vis = h.path():moveto(origin):lineto(o2):evenly():stroke("silver"):strokewidth(0.5):build()
    local a2vis2 = h.path():moveto(origin):lineto(o3):evenly():stroke("silver"):strokewidth(0.5):build()
    local topbarVis = hbar:shifted(topbarleft.x, topbarleft.y):stroke("silver"):strokewidth(0.5)
    local botbarVis = hbar:shifted(bottombarleft.x, bottombarleft.y):stroke("silver"):strokewidth(0.5)
    local clipVis = clippath:stroke("tomato"):strokewidth(0.5)

    svg:add(inner):add(outerVis):add(a1vis):add(a2vis):add(a2vis2)
    svg:add(topbarVis):add(botbarVis):add(clipVis)
end

if show_dots or show_labels then
    local dbgpic = h.picture()
    local gray = h.color("gray")
    local blue = h.color("navy")
    local col = show_labels and blue or gray
    local labeloptions = { color = col, fontsize = 5 }

    if show_dots then
        dbgpic:dotlabel("", origin,        "c",   gray)
        dbgpic:dotlabel("", c1,            "c",   gray)
        dbgpic:dotlabel("", c2,            "c",   gray)
        dbgpic:dotlabel("", c3,            "c",   gray)
        dbgpic:dotlabel("", c4,            "c",   gray)
        dbgpic:dotlabel("", c5,            "c",   gray)
        dbgpic:dotlabel("", o1,            "c",   gray)
        dbgpic:dotlabel("", o2,            "c",   gray)
        dbgpic:dotlabel("", topbarleft,    "c",   gray)
        dbgpic:dotlabel("", bottombarleft, "c",   gray)
    end
    if show_labels then
        dbgpic:label("origin",        origin,        "left", labeloptions)
        dbgpic:label("c1",            c1,            "rt",   labeloptions)
        dbgpic:label("c2",            c2,            "rt",   labeloptions)
        dbgpic:label("c3",            c3,            "rt",   labeloptions)
        dbgpic:label("c4",            c4,            "rt",   labeloptions)
        dbgpic:label("c5",            c5,            "bot",  labeloptions)
        dbgpic:label("o1",            o1,            "bot",  labeloptions)
        dbgpic:label("o2",            o2,            "bot",  labeloptions)
        dbgpic:label("topbarleft",    topbarleft,    "bot",  labeloptions)
        dbgpic:label("bottombarleft", bottombarleft, "bot",  labeloptions)
    end
    svg:addpicture(dbgpic)
end

svg:write("euro.svg")
print("Created euro.svg")
