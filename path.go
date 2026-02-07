package hobby

import (
	"github.com/boxesandglue/mpgo/draw"
	"github.com/boxesandglue/mpgo/mp"
	"github.com/boxesandglue/mpgo/svg"
	lua "github.com/speedata/go-lua"
)

// PathBuilder wraps draw.PathBuilder for Lua
type PathBuilder struct {
	builder *draw.PathBuilder
}

// luaBuildCycle creates a closed region from multiple paths: hobby.buildcycle(p1, p2, ...)
func luaBuildCycle(l *lua.State) int {
	n := l.Top()
	if n == 0 {
		l.PushNil()
		return 1
	}
	paths := make([]*mp.Path, n)
	for i := 1; i <= n; i++ {
		paths[i-1] = checkPath(l, i)
	}
	result := mp.BuildCycle(paths...)
	if result == nil {
		l.PushNil()
		return 1
	}
	pushPath(l, result)
	return 1
}

// luaNewPath creates a new path builder: hobby.path()
func luaNewPath(l *lua.State) int {
	pb := &PathBuilder{
		builder: draw.NewPath(),
	}
	l.PushUserData(pb)
	lua.SetMetaTableNamed(l, "hobby.pathbuilder")
	return 1
}

// luaFullCircle returns the predefined fullcircle path
func luaFullCircle(l *lua.State) int {
	path := mp.FullCircle()
	e := mp.NewEngine()
	e.AddPath(path)
	e.Solve()
	pushPath(l, path)
	return 1
}

// luaHalfCircle returns the predefined halfcircle path
func luaHalfCircle(l *lua.State) int {
	path := mp.HalfCircle()
	e := mp.NewEngine()
	e.AddPath(path)
	e.Solve()
	pushPath(l, path)
	return 1
}

// luaQuarterCircle returns the predefined quartercircle path
func luaQuarterCircle(l *lua.State) int {
	path := mp.QuarterCircle()
	e := mp.NewEngine()
	e.AddPath(path)
	e.Solve()
	pushPath(l, path)
	return 1
}

// luaUnitSquare returns the predefined unitsquare path
func luaUnitSquare(l *lua.State) int {
	path := mp.UnitSquare()
	e := mp.NewEngine()
	e.AddPath(path)
	e.Solve()
	pushPath(l, path)
	return 1
}

// registerPathMeta registers metatables for paths and path builders
func registerPathMeta(l *lua.State) {
	// PathBuilder metatable
	lua.NewMetaTable(l, "hobby.pathbuilder")
	l.PushGoFunction(pathBuilderIndex)
	l.SetField(-2, "__index")
	l.Pop(1)

	// Path metatable
	lua.NewMetaTable(l, "hobby.path")
	l.PushGoFunction(pathIndex)
	l.SetField(-2, "__index")
	l.Pop(1)
}

func pathBuilderIndex(l *lua.State) int {
	pb := l.ToUserData(1).(*PathBuilder)
	key := lua.CheckString(l, 2)

	switch key {
	case "moveto":
		l.PushGoFunction(func(l *lua.State) int {
			p := checkPoint(l, 2)
			pb.builder.MoveTo(p)
			l.PushValue(1) // return self for chaining
			return 1
		})
		return 1

	case "lineto":
		l.PushGoFunction(func(l *lua.State) int {
			p := checkPoint(l, 2)
			pb.builder.LineTo(p)
			l.PushValue(1)
			return 1
		})
		return 1

	case "curveto":
		l.PushGoFunction(func(l *lua.State) int {
			p := checkPoint(l, 2)
			pb.builder.CurveTo(p)
			l.PushValue(1)
			return 1
		})
		return 1

	case "curvetowithcontrols":
		// curvetowithcontrols(pt, c1, c2) - curve with explicit control points
		l.PushGoFunction(func(l *lua.State) int {
			pt := checkPoint(l, 2)
			c1 := checkPoint(l, 3)
			c2 := checkPoint(l, 4)
			pb.builder.CurveToWithControls(pt, c1, c2)
			l.PushValue(1)
			return 1
		})
		return 1

	case "dir":
		l.PushGoFunction(func(l *lua.State) int {
			angle := lua.CheckNumber(l, 2)
			pb.builder.WithDirection(angle)
			l.PushValue(1)
			return 1
		})
		return 1

	case "indir":
		// Incoming direction (direction at arrival)
		l.PushGoFunction(func(l *lua.State) int {
			angle := lua.CheckNumber(l, 2)
			pb.builder.WithIncomingDirection(angle)
			l.PushValue(1)
			return 1
		})
		return 1

	case "curl":
		l.PushGoFunction(func(l *lua.State) int {
			c := lua.CheckNumber(l, 2)
			pb.builder.WithCurl(c)
			l.PushValue(1)
			return 1
		})
		return 1

	case "outcurl":
		// Outgoing curl (curl at departure)
		l.PushGoFunction(func(l *lua.State) int {
			c := lua.CheckNumber(l, 2)
			pb.builder.WithOutgoingCurl(c)
			l.PushValue(1)
			return 1
		})
		return 1

	case "incurl":
		// Incoming curl (curl at arrival)
		l.PushGoFunction(func(l *lua.State) int {
			c := lua.CheckNumber(l, 2)
			pb.builder.WithIncomingCurl(c)
			l.PushValue(1)
			return 1
		})
		return 1

	case "tension":
		l.PushGoFunction(func(l *lua.State) int {
			t := lua.CheckNumber(l, 2)
			pb.builder.WithTension(t)
			l.PushValue(1)
			return 1
		})
		return 1

	case "outtension":
		// Outgoing tension
		l.PushGoFunction(func(l *lua.State) int {
			t := lua.CheckNumber(l, 2)
			pb.builder.WithOutgoingTension(t)
			l.PushValue(1)
			return 1
		})
		return 1

	case "intension":
		// Incoming tension
		l.PushGoFunction(func(l *lua.State) int {
			t := lua.CheckNumber(l, 2)
			pb.builder.WithIncomingTension(t)
			l.PushValue(1)
			return 1
		})
		return 1

	case "tensionatleast":
		// Tension atleast (minimum tension)
		l.PushGoFunction(func(l *lua.State) int {
			t := lua.CheckNumber(l, 2)
			pb.builder.WithTensionAtLeast(t)
			l.PushValue(1)
			return 1
		})
		return 1

	case "tensioninfinity":
		// Infinite tension (straight line segment with smooth connections)
		l.PushGoFunction(func(l *lua.State) int {
			pb.builder.WithTensionInfinity()
			l.PushValue(1)
			return 1
		})
		return 1

	case "stroke":
		l.PushGoFunction(func(l *lua.State) int {
			c := checkColor(l, 2)
			pb.builder.WithStrokeColor(c)
			l.PushValue(1)
			return 1
		})
		return 1

	case "strokewidth":
		l.PushGoFunction(func(l *lua.State) int {
			w := lua.CheckNumber(l, 2)
			pb.builder.WithStrokeWidth(w)
			l.PushValue(1)
			return 1
		})
		return 1

	case "fill":
		l.PushGoFunction(func(l *lua.State) int {
			c := checkColor(l, 2)
			pb.builder.WithFill(c)
			l.PushValue(1)
			return 1
		})
		return 1

	case "pen":
		l.PushGoFunction(func(l *lua.State) int {
			pen := checkPen(l, 2)
			pb.builder.WithPen(pen)
			l.PushValue(1)
			return 1
		})
		return 1

	case "dash":
		l.PushGoFunction(func(l *lua.State) int {
			dash := checkDash(l, 2)
			pb.builder.WithDashPattern(dash)
			l.PushValue(1)
			return 1
		})
		return 1

	case "evenly":
		l.PushGoFunction(func(l *lua.State) int {
			pb.builder.DashedEvenly()
			l.PushValue(1)
			return 1
		})
		return 1

	case "withdots":
		l.PushGoFunction(func(l *lua.State) int {
			pb.builder.DashedWithDots()
			l.PushValue(1)
			return 1
		})
		return 1

	case "arrow":
		l.PushGoFunction(func(l *lua.State) int {
			pb.builder.WithArrow()
			l.PushValue(1)
			return 1
		})
		return 1

	case "dblarrow":
		l.PushGoFunction(func(l *lua.State) int {
			pb.builder.WithDoubleArrow()
			l.PushValue(1)
			return 1
		})
		return 1

	case "arrowstyle":
		l.PushGoFunction(func(l *lua.State) int {
			length := lua.CheckNumber(l, 2)
			angle := lua.CheckNumber(l, 3)
			pb.builder.WithArrowStyle(length, angle)
			l.PushValue(1)
			return 1
		})
		return 1

	case "linejoin":
		l.PushGoFunction(func(l *lua.State) int {
			join := lua.CheckString(l, 2)
			switch join {
			case "miter":
				pb.builder.WithLineJoin(mp.LineJoinMiter)
			case "round":
				pb.builder.WithLineJoin(mp.LineJoinRound)
			case "bevel":
				pb.builder.WithLineJoin(mp.LineJoinBevel)
			default:
				lua.Errorf(l, "unknown linejoin: %s (use miter, round, bevel)", join)
			}
			l.PushValue(1)
			return 1
		})
		return 1

	case "linecap":
		l.PushGoFunction(func(l *lua.State) int {
			cap := lua.CheckString(l, 2)
			switch cap {
			case "butt":
				pb.builder.WithLineCap(mp.LineCapButt)
			case "round":
				pb.builder.WithLineCap(mp.LineCapRounded)
			case "square":
				pb.builder.WithLineCap(mp.LineCapSquared)
			default:
				lua.Errorf(l, "unknown linecap: %s (use butt, round, square)", cap)
			}
			l.PushValue(1)
			return 1
		})
		return 1

	case "close":
		l.PushGoFunction(func(l *lua.State) int {
			pb.builder.Close()
			l.PushValue(1)
			return 1
		})
		return 1

	case "cycle":
		// Alias for close
		l.PushGoFunction(func(l *lua.State) int {
			pb.builder.Close()
			l.PushValue(1)
			return 1
		})
		return 1

	case "build":
		l.PushGoFunction(func(l *lua.State) int {
			path, err := pb.builder.Solve()
			if err != nil {
				lua.Errorf(l, "path build error: %s", err.Error())
				return 0
			}
			pushPath(l, path)
			return 1
		})
		return 1
	}

	return 0
}

func pathIndex(l *lua.State) int {
	path := checkPath(l, 1)
	key := lua.CheckString(l, 2)

	switch key {
	case "length":
		// Number of segments in the path (MetaPost's "length p")
		l.PushInteger(path.PathLength())
		return 1

	case "arclength":
		l.PushNumber(path.ArcLength())
		return 1

	case "arctime":
		l.PushGoFunction(func(l *lua.State) int {
			arcLen := lua.CheckNumber(l, 2)
			t := path.ArcTime(mp.Number(arcLen))
			l.PushNumber(float64(t))
			return 1
		})
		return 1

	case "intersectiontimes":
		l.PushGoFunction(func(l *lua.State) int {
			other := checkPath(l, 2)
			t1, t2 := path.IntersectionTimes(other)
			l.PushNumber(float64(t1))
			l.PushNumber(float64(t2))
			return 2
		})
		return 1

	case "intersectionpoint":
		l.PushGoFunction(func(l *lua.State) int {
			other := checkPath(l, 2)
			x, y, found := path.IntersectionPoint(other)
			if !found {
				l.PushNil()
				return 1
			}
			pushPoint(l, mp.P(float64(x), float64(y)))
			return 1
		})
		return 1

	case "pointat":
		l.PushGoFunction(func(l *lua.State) int {
			t := lua.CheckNumber(l, 2)
			x, y := path.PointOf(mp.Number(t))
			pushPoint(l, mp.P(float64(x), float64(y)))
			return 1
		})
		return 1

	case "directionat":
		l.PushGoFunction(func(l *lua.State) int {
			t := lua.CheckNumber(l, 2)
			dx, dy := path.DirectionOf(mp.Number(t))
			pushPoint(l, mp.P(float64(dx), float64(dy)))
			return 1
		})
		return 1

	case "directiontime":
		// path:directiontime(dx, dy) - find t where tangent equals (dx, dy)
		l.PushGoFunction(func(l *lua.State) int {
			dx := lua.CheckNumber(l, 2)
			dy := lua.CheckNumber(l, 3)
			t := path.DirectionTimeOf(mp.Number(dx), mp.Number(dy))
			l.PushNumber(float64(t))
			return 1
		})
		return 1

	case "directionpoint":
		// path:directionpoint(dx, dy) - find point where tangent equals (dx, dy)
		l.PushGoFunction(func(l *lua.State) int {
			dx := lua.CheckNumber(l, 2)
			dy := lua.CheckNumber(l, 3)
			x, y, found := path.DirectionPointOf(mp.Number(dx), mp.Number(dy))
			if !found {
				l.PushNil()
				return 1
			}
			pushPoint(l, mp.P(float64(x), float64(y)))
			return 1
		})
		return 1

	case "precontrol":
		l.PushGoFunction(func(l *lua.State) int {
			t := lua.CheckNumber(l, 2)
			x, y := path.PrecontrolOf(mp.Number(t))
			pushPoint(l, mp.P(float64(x), float64(y)))
			return 1
		})
		return 1

	case "postcontrol":
		l.PushGoFunction(func(l *lua.State) int {
			t := lua.CheckNumber(l, 2)
			x, y := path.PostcontrolOf(mp.Number(t))
			pushPoint(l, mp.P(float64(x), float64(y)))
			return 1
		})
		return 1

	case "subpath":
		l.PushGoFunction(func(l *lua.State) int {
			t1 := lua.CheckNumber(l, 2)
			t2 := lua.CheckNumber(l, 3)
			sub := path.Subpath(mp.Number(t1), mp.Number(t2))
			pushPath(l, sub)
			return 1
		})
		return 1

	case "cutbefore":
		// path:cutbefore(other) - cut self before intersection with other
		l.PushGoFunction(func(l *lua.State) int {
			other := checkPath(l, 2)
			cut := path.CutBefore(other)
			pushPath(l, cut)
			return 1
		})
		return 1

	case "cutafter":
		// path:cutafter(other) - cut self after intersection with other
		l.PushGoFunction(func(l *lua.State) int {
			other := checkPath(l, 2)
			cut := path.CutAfter(other)
			pushPath(l, cut)
			return 1
		})
		return 1

	case "reversed":
		l.PushGoFunction(func(l *lua.State) int {
			rev := path.Reversed()
			pushPath(l, rev)
			return 1
		})
		return 1

	case "scaled":
		l.PushGoFunction(func(l *lua.State) int {
			s := lua.CheckNumber(l, 2)
			scaled := path.Scaled(s)
			pushPath(l, scaled)
			return 1
		})
		return 1

	case "xscaled":
		l.PushGoFunction(func(l *lua.State) int {
			s := lua.CheckNumber(l, 2)
			scaled := path.XScaled(mp.Number(s))
			pushPath(l, scaled)
			return 1
		})
		return 1

	case "yscaled":
		l.PushGoFunction(func(l *lua.State) int {
			s := lua.CheckNumber(l, 2)
			scaled := path.YScaled(mp.Number(s))
			pushPath(l, scaled)
			return 1
		})
		return 1

	case "slanted":
		l.PushGoFunction(func(l *lua.State) int {
			s := lua.CheckNumber(l, 2)
			slanted := path.Slanted(mp.Number(s))
			pushPath(l, slanted)
			return 1
		})
		return 1

	case "shifted":
		l.PushGoFunction(func(l *lua.State) int {
			dx := lua.CheckNumber(l, 2)
			dy := lua.CheckNumber(l, 3)
			shifted := path.Shifted(dx, dy)
			pushPath(l, shifted)
			return 1
		})
		return 1

	case "rotated":
		l.PushGoFunction(func(l *lua.State) int {
			angle := lua.CheckNumber(l, 2)
			rotated := path.Rotated(angle)
			pushPath(l, rotated)
			return 1
		})
		return 1

	case "zscaled":
		l.PushGoFunction(func(l *lua.State) int {
			// zscaled(a, b) or zscaled(point)
			if l.IsUserData(2) {
				p := checkPoint(l, 2)
				scaled := path.ZScaled(mp.Number(p.X), mp.Number(p.Y))
				pushPath(l, scaled)
			} else {
				a := lua.CheckNumber(l, 2)
				b := lua.CheckNumber(l, 3)
				scaled := path.ZScaled(mp.Number(a), mp.Number(b))
				pushPath(l, scaled)
			}
			return 1
		})
		return 1

	case "rotatedaround":
		l.PushGoFunction(func(l *lua.State) int {
			// rotatedaround(point, angle) or rotatedaround(cx, cy, angle)
			if l.IsUserData(2) {
				p := checkPoint(l, 2)
				angle := lua.CheckNumber(l, 3)
				rotated := path.RotatedAround(mp.Number(p.X), mp.Number(p.Y), mp.Number(angle))
				pushPath(l, rotated)
			} else {
				cx := lua.CheckNumber(l, 2)
				cy := lua.CheckNumber(l, 3)
				angle := lua.CheckNumber(l, 4)
				rotated := path.RotatedAround(mp.Number(cx), mp.Number(cy), mp.Number(angle))
				pushPath(l, rotated)
			}
			return 1
		})
		return 1

	case "scaledaround":
		l.PushGoFunction(func(l *lua.State) int {
			// scaledaround(point, scale) or scaledaround(cx, cy, scale)
			if l.IsUserData(2) {
				p := checkPoint(l, 2)
				s := lua.CheckNumber(l, 3)
				scaled := path.ScaledAround(mp.Number(p.X), mp.Number(p.Y), mp.Number(s))
				pushPath(l, scaled)
			} else {
				cx := lua.CheckNumber(l, 2)
				cy := lua.CheckNumber(l, 3)
				s := lua.CheckNumber(l, 4)
				scaled := path.ScaledAround(mp.Number(cx), mp.Number(cy), mp.Number(s))
				pushPath(l, scaled)
			}
			return 1
		})
		return 1

	case "reflectedabout":
		l.PushGoFunction(func(l *lua.State) int {
			// reflectedabout(p1, p2) or reflectedabout(x1, y1, x2, y2)
			if l.IsUserData(2) {
				p1 := checkPoint(l, 2)
				p2 := checkPoint(l, 3)
				reflected := path.ReflectedAbout(mp.Number(p1.X), mp.Number(p1.Y), mp.Number(p2.X), mp.Number(p2.Y))
				pushPath(l, reflected)
			} else {
				x1 := lua.CheckNumber(l, 2)
				y1 := lua.CheckNumber(l, 3)
				x2 := lua.CheckNumber(l, 4)
				y2 := lua.CheckNumber(l, 5)
				reflected := path.ReflectedAbout(mp.Number(x1), mp.Number(y1), mp.Number(x2), mp.Number(y2))
				pushPath(l, reflected)
			}
			return 1
		})
		return 1

	case "stroke":
		l.PushGoFunction(func(l *lua.State) int {
			c := checkColor(l, 2)
			path.Style.Stroke = c
			l.PushValue(1)
			return 1
		})
		return 1

	case "strokewidth":
		l.PushGoFunction(func(l *lua.State) int {
			w := lua.CheckNumber(l, 2)
			path.Style.StrokeWidth = w
			l.PushValue(1)
			return 1
		})
		return 1

	case "fill":
		l.PushGoFunction(func(l *lua.State) int {
			c := checkColor(l, 2)
			path.Style.Fill = c
			l.PushValue(1)
			return 1
		})
		return 1

	case "pen":
		l.PushGoFunction(func(l *lua.State) int {
			pen := checkPen(l, 2)
			path.Style.Pen = pen
			l.PushValue(1)
			return 1
		})
		return 1

	case "dash":
		l.PushGoFunction(func(l *lua.State) int {
			dash := checkDash(l, 2)
			path.Style.Dash = dash
			l.PushValue(1)
			return 1
		})
		return 1

	case "evenly":
		l.PushGoFunction(func(l *lua.State) int {
			path.Style.Dash = mp.DashEvenly()
			l.PushValue(1)
			return 1
		})
		return 1

	case "withdots":
		l.PushGoFunction(func(l *lua.State) int {
			path.Style.Dash = mp.DashWithDots()
			l.PushValue(1)
			return 1
		})
		return 1

	case "arrow":
		l.PushGoFunction(func(l *lua.State) int {
			path.Style.Arrow.End = true
			if path.Style.Arrow.Length == 0 {
				path.Style.Arrow.Length = mp.DefaultAHLength
			}
			if path.Style.Arrow.Angle == 0 {
				path.Style.Arrow.Angle = mp.DefaultAHAngle
			}
			l.PushValue(1)
			return 1
		})
		return 1

	case "dblarrow":
		l.PushGoFunction(func(l *lua.State) int {
			path.Style.Arrow.Start = true
			path.Style.Arrow.End = true
			if path.Style.Arrow.Length == 0 {
				path.Style.Arrow.Length = mp.DefaultAHLength
			}
			if path.Style.Arrow.Angle == 0 {
				path.Style.Arrow.Angle = mp.DefaultAHAngle
			}
			l.PushValue(1)
			return 1
		})
		return 1

	case "arrowstyle":
		l.PushGoFunction(func(l *lua.State) int {
			length := lua.CheckNumber(l, 2)
			angle := lua.CheckNumber(l, 3)
			path.Style.Arrow.Length = mp.Number(length)
			path.Style.Arrow.Angle = mp.Number(angle)
			l.PushValue(1)
			return 1
		})
		return 1

	case "linejoin":
		l.PushGoFunction(func(l *lua.State) int {
			join := lua.CheckString(l, 2)
			switch join {
			case "miter":
				path.Style.LineJoin = mp.LineJoinMiter
			case "round":
				path.Style.LineJoin = mp.LineJoinRound
			case "bevel":
				path.Style.LineJoin = mp.LineJoinBevel
			default:
				lua.Errorf(l, "unknown linejoin: %s (use miter, round, bevel)", join)
			}
			l.PushValue(1)
			return 1
		})
		return 1

	case "linecap":
		l.PushGoFunction(func(l *lua.State) int {
			cap := lua.CheckString(l, 2)
			switch cap {
			case "butt":
				path.Style.LineCap = mp.LineCapButt
			case "round":
				path.Style.LineCap = mp.LineCapRounded
			case "square":
				path.Style.LineCap = mp.LineCapSquared
			default:
				lua.Errorf(l, "unknown linecap: %s (use butt, round, square)", cap)
			}
			l.PushValue(1)
			return 1
		})
		return 1

	case "llcorner":
		minX, minY, _, _ := svg.PathBBox(path)
		pushPoint(l, mp.P(minX, minY))
		return 1

	case "lrcorner":
		_, minY, maxX, _ := svg.PathBBox(path)
		pushPoint(l, mp.P(maxX, minY))
		return 1

	case "ulcorner":
		minX, _, _, maxY := svg.PathBBox(path)
		pushPoint(l, mp.P(minX, maxY))
		return 1

	case "urcorner":
		_, _, maxX, maxY := svg.PathBBox(path)
		pushPoint(l, mp.P(maxX, maxY))
		return 1

	case "center":
		minX, minY, maxX, maxY := svg.PathBBox(path)
		pushPoint(l, mp.P((minX+maxX)/2, (minY+maxY)/2))
		return 1

	case "bbox":
		l.PushGoFunction(func(l *lua.State) int {
			minX, minY, maxX, maxY := svg.PathBBox(path)
			pushPath(l, bboxPath(minX, minY, maxX, maxY))
			return 1
		})
		return 1
	}

	return 0
}

// bboxPath creates a closed rectangular path from bounding box coordinates.
func bboxPath(minX, minY, maxX, maxY float64) *mp.Path {
	coords := [][2]float64{
		{minX, minY},
		{maxX, minY},
		{maxX, maxY},
		{minX, maxY},
	}
	p := mp.NewPath()
	var knots []*mp.Knot
	for _, c := range coords {
		k := mp.NewKnot()
		k.XCoord = mp.Number(c[0])
		k.YCoord = mp.Number(c[1])
		k.LeftX = k.XCoord
		k.LeftY = k.YCoord
		k.RightX = k.XCoord
		k.RightY = k.YCoord
		k.LType = mp.KnotExplicit
		k.RType = mp.KnotExplicit
		knots = append(knots, k)
	}
	for i, k := range knots {
		p.Append(k)
		if i > 0 {
			knots[i-1].Next = k
			k.Prev = knots[i-1]
		}
	}
	knots[len(knots)-1].Next = knots[0]
	knots[0].Prev = knots[len(knots)-1]
	return p
}
