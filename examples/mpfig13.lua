-- MetaPost manual figure 13: line intersections
-- Translated from:
--   beginfig(13);
--   z1=-z2=(.2in,0);
--   x3=-x6=.3in;
--   x3+y3=x6+y6=1.1in;
--   z4=1/3[z3,z6]; z5=2/3[z3,z6];
--   z20=whatever[z1,z3]=whatever[z2,z4];
--   z30=whatever[z1,z4]=whatever[z2,z5];
--   z40=whatever[z1,z5]=whatever[z2,z6];
--   draw z1--z20--z2--z30--z1--z40--z2;
--   pickup pencircle scaled 1pt;
--   draw z1--z2; draw z3--z6;
--   endfig;

local h = require("hobby")

-- 1in = 72bp in MetaPost
local inch = 72

-- Phase 1: Solve z3, z6 via LinearXY
-- x3=-x6=.3in; x3+y3=x6+y6=1.1in
local eq = h.context()
local z3 = eq:unknown()
eq:eqx(z3, 0.3 * inch)
eq:linearxy(z3, 1, 1, 1.1 * inch)

local z6 = eq:unknown()
eq:eqx(z6, -0.3 * inch)
eq:linearxy(z6, 1, 1, 1.1 * inch)

eq:solve()

-- Phase 2: geometry with solved coordinates
local ctx = h.context()

-- z1=-z2=(.2in,0)
local z1 = ctx:known(0.2 * inch, 0)
local z2 = ctx:known(-0.2 * inch, 0)

local z3k = ctx:known(z3.x, z3.y)
local z6k = ctx:known(z6.x, z6.y)

-- z4=1/3[z3,z6]; z5=2/3[z3,z6]
local z4 = ctx:betweenat(z3k, z6k, 1/3)
local z5 = ctx:betweenat(z3k, z6k, 2/3)

ctx:solve()

-- z20=whatever[z1,z3]=whatever[z2,z4]  (intersection of two lines)
local z20 = ctx:intersectionof(z1, z3k, z2, z4)

-- z30=whatever[z1,z4]=whatever[z2,z5]
local z30 = ctx:intersectionof(z1, z4, z2, z5)

-- z40=whatever[z1,z5]=whatever[z2,z6]
local z40 = ctx:intersectionof(z1, z5, z2, z6k)

ctx:solve()

-- Build a picture for paths + labels
local pic = h.picture()

-- draw z1--z20--z2--z30--z1--z40--z2;
local main = h.path()
    :moveto(z1:point())
    :lineto(z20:point())
    :lineto(z2:point())
    :lineto(z30:point())
    :lineto(z1:point())
    :lineto(z40:point())
    :lineto(z2:point())
    :stroke("black")
    :strokewidth(0.5)
    :build()

-- pickup pencircle scaled 1pt; draw z1--z2;
local baseline = h.path()
    :moveto(z1:point())
    :lineto(z2:point())
    :stroke("black")
    :strokewidth(1)
    :build()

-- draw z3--z6;
local diagonal = h.path()
    :moveto(z3k:point())
    :lineto(z6k:point())
    :stroke("black")
    :strokewidth(1)
    :build()

pic:add(main)
pic:add(baseline)
pic:add(diagonal)

-- Dot labels
local black = h.color("black")
pic:dotlabel("1", z1:point(), "bot", black)
pic:dotlabel("2", z2:point(), "bot", black)
pic:dotlabel("3", z3k:point(), "rt", black)
pic:dotlabel("6", z6k:point(), "lft", black)
pic:dotlabel("20", z20:point(), "top", black)
pic:dotlabel("30", z30:point(), "top", black)
pic:dotlabel("40", z40:point(), "top", black)

h.svg()
    :padding(10)
    :addpicture(pic)
    :write("mpfig13.svg")

print("Created mpfig13.svg")
