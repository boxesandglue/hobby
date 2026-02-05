-- labels.lua: Example demonstrating label support in hobby
local h = require("hobby")

-- Load a font for converting labels to glyph outlines (optional)
-- Without a font, labels are rendered as SVG <text> elements.
local face = h.loadfont("/System/Library/Fonts/Supplemental/Arial.ttf")

-- Create a picture
local pic = h.picture()

-- Draw a triangle
local z0 = h.point(0, 0)
local z1 = h.point(100, 0)
local z2 = h.point(50, 86.6)

local triangle = h.path()
    :moveto(z0)
    :lineto(z1)
    :lineto(z2)
    :close()
    :stroke(h.color("black"))
    :build()

pic:add(triangle)

-- Add labels at the vertices
pic:label("A", z0, "llft")   -- lower-left
pic:label("B", z1, "lrt")    -- lower-right
pic:label("C", z2, "top")    -- top

-- Add dot labels for midpoints
local m01 = h.midpoint(z0, z1)
local m12 = h.midpoint(z1, z2)
local m02 = h.midpoint(z0, z2)

pic:dotlabel("a", m12, "rt", h.color("blue"))
pic:dotlabel("b", m02, "lft", h.color("blue"))
pic:dotlabel("c", m01, "bot", h.color("blue"))

-- Add centroid with red label
local centroid = h.point(50, 28.87)
pic:dotlabel("S", centroid, "urt", h.color("red"))

-- Convert labels to glyph outline paths
pic:converttopaths(face)

-- Output SVG
local svg = h.svg()
    :padding(15)
    :addpicture(pic)
    :write("labels.svg")

print("Created labels.svg")
