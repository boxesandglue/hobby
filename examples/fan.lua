-- Fan of curves (MetaPost manual figure 7)
-- 10 curves from the same start point, all leaving at 45°
-- but arriving at different angles: 0°, -10°, -20°, ..., -90°

local h = require("hobby")

local cm = 28.35
local svg = h.svg():padding(5)

for a = 0, 9 do
    local path = h.path()
        :moveto(h.point(0, 0))
        :dir(45)
        :curveto(h.point(6 * cm, 0))
        :indir(-10 * a)
        :stroke("black")
        :strokewidth(0.5)
        :build()
    svg:add(path)
end

svg:write("fan.svg")
print("Created fan.svg")
