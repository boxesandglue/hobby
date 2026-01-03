-- Transformation example

local h = require("hobby")

-- Original circle
local circle = h.fullcircle():scaled(20):shifted(30, 80):stroke("black")

-- xscaled (stretched horizontally)
local xscaled = h.fullcircle():scaled(20):xscaled(2):shifted(90, 80):stroke("blue")

-- yscaled (stretched vertically)
local yscaled = h.fullcircle():scaled(20):yscaled(2):shifted(150, 80):stroke("red")

-- slanted (sheared)
local slanted = h.unitsquare():scaled(30):slanted(0.5):shifted(30, 20):stroke("green")

-- Combination
local combo = h.fullcircle()
    :scaled(20)
    :xscaled(1.5)
    :yscaled(0.7)
    :slanted(0.3)
    :shifted(130, 30)
    :stroke("purple")

-- Output to SVG
local svg = h.svg()
    :padding(10)
    :add(circle)
    :add(xscaled)
    :add(yscaled)
    :add(slanted)
    :add(combo)
    :write("transforms.svg")

print("Created transforms.svg")
