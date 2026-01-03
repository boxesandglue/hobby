package hobby

import (
	"github.com/boxesandglue/mpgo/draw"
	"github.com/boxesandglue/mpgo/mp"
	lua "github.com/speedata/go-lua"
)

// luaNewContext creates a new equation-solving context: h.context()
func luaNewContext(l *lua.State) int {
	ctx := draw.NewContext()
	pushContext(l, ctx)
	return 1
}

// registerContextMeta registers the metatable for Context
func registerContextMeta(l *lua.State) {
	lua.NewMetaTable(l, "hobby.context")
	l.PushGoFunction(contextIndex)
	l.SetField(-2, "__index")
	l.Pop(1)
}

// registerVarMeta registers the metatable for Var (point variables)
func registerVarMeta(l *lua.State) {
	lua.NewMetaTable(l, "hobby.var")
	l.PushGoFunction(varIndex)
	l.SetField(-2, "__index")
	l.Pop(1)
}

// pushContext pushes a Context as userdata
func pushContext(l *lua.State, ctx *draw.Context) {
	l.PushUserData(ctx)
	lua.SetMetaTableNamed(l, "hobby.context")
}

// checkContext checks if value at index is a Context
func checkContext(l *lua.State, index int) *draw.Context {
	ud := l.ToUserData(index)
	if ctx, ok := ud.(*draw.Context); ok {
		return ctx
	}
	lua.Errorf(l, "expected context at argument %d", index)
	return nil
}

// pushVar pushes a Var as userdata
func pushVar(l *lua.State, v *draw.Var) {
	l.PushUserData(v)
	lua.SetMetaTableNamed(l, "hobby.var")
}

// checkVar checks if value at index is a Var
func checkVar(l *lua.State, index int) *draw.Var {
	ud := l.ToUserData(index)
	if v, ok := ud.(*draw.Var); ok {
		return v
	}
	lua.Errorf(l, "expected var at argument %d", index)
	return nil
}

func contextIndex(l *lua.State) int {
	ctx := checkContext(l, 1)
	key := lua.CheckString(l, 2)

	switch key {
	case "unknown":
		// ctx:unknown() - create unknown point variable
		l.PushGoFunction(func(l *lua.State) int {
			v := ctx.Unknown()
			pushVar(l, v)
			return 1
		})
		return 1

	case "known":
		// ctx:known(x, y) - create known point variable
		l.PushGoFunction(func(l *lua.State) int {
			x := lua.CheckNumber(l, 2)
			y := lua.CheckNumber(l, 3)
			v := ctx.Known(x, y)
			pushVar(l, v)
			return 1
		})
		return 1

	case "point":
		// ctx:point() - alias for unknown
		l.PushGoFunction(func(l *lua.State) int {
			v := ctx.Point()
			pushVar(l, v)
			return 1
		})
		return 1

	case "points":
		// ctx:points(n) - create n unknown points, returns table
		l.PushGoFunction(func(l *lua.State) int {
			n := lua.CheckInteger(l, 2)
			vars := ctx.Points(n)
			l.CreateTable(len(vars), 0)
			for i, v := range vars {
				pushVar(l, v)
				l.RawSetInt(-2, i+1)
			}
			return 1
		})
		return 1

	case "eq":
		// ctx:eq(v, point) - constrain v to equal point
		l.PushGoFunction(func(l *lua.State) int {
			v := checkVar(l, 2)
			p := checkPoint(l, 3)
			ctx.Eq(v, p)
			l.PushValue(1)
			return 1
		})
		return 1

	case "eqx":
		// ctx:eqx(v, x) - constrain v.x to value
		l.PushGoFunction(func(l *lua.State) int {
			v := checkVar(l, 2)
			x := lua.CheckNumber(l, 3)
			ctx.EqX(v, x)
			l.PushValue(1)
			return 1
		})
		return 1

	case "eqy":
		// ctx:eqy(v, y) - constrain v.y to value
		l.PushGoFunction(func(l *lua.State) int {
			v := checkVar(l, 2)
			y := lua.CheckNumber(l, 3)
			ctx.EqY(v, y)
			l.PushValue(1)
			return 1
		})
		return 1

	case "eqvar":
		// ctx:eqvar(a, b) - constrain a = b
		l.PushGoFunction(func(l *lua.State) int {
			a := checkVar(l, 2)
			b := checkVar(l, 3)
			ctx.EqVar(a, b)
			l.PushValue(1)
			return 1
		})
		return 1

	case "eqvarx":
		// ctx:eqvarx(a, b) - constrain a.x = b.x
		l.PushGoFunction(func(l *lua.State) int {
			a := checkVar(l, 2)
			b := checkVar(l, 3)
			ctx.EqVarX(a, b)
			l.PushValue(1)
			return 1
		})
		return 1

	case "eqvary":
		// ctx:eqvary(a, b) - constrain a.y = b.y
		l.PushGoFunction(func(l *lua.State) int {
			a := checkVar(l, 2)
			b := checkVar(l, 3)
			ctx.EqVarY(a, b)
			l.PushValue(1)
			return 1
		})
		return 1

	case "midpoint":
		// ctx:midpoint(m, a, b) - constrain m = midpoint of a and b
		l.PushGoFunction(func(l *lua.State) int {
			m := checkVar(l, 2)
			a := checkVar(l, 3)
			b := checkVar(l, 4)
			ctx.MidPoint(m, a, b)
			l.PushValue(1)
			return 1
		})
		return 1

	case "midpointof":
		// ctx:midpointof(a, b) - returns new var at midpoint
		l.PushGoFunction(func(l *lua.State) int {
			a := checkVar(l, 2)
			b := checkVar(l, 3)
			m := ctx.MidPointOf(a, b)
			pushVar(l, m)
			return 1
		})
		return 1

	case "between":
		// ctx:between(p, a, b, t) - constrain p = t[a,b]
		l.PushGoFunction(func(l *lua.State) int {
			p := checkVar(l, 2)
			a := checkVar(l, 3)
			b := checkVar(l, 4)
			t := lua.CheckNumber(l, 5)
			ctx.Between(p, a, b, t)
			l.PushValue(1)
			return 1
		})
		return 1

	case "betweenat":
		// ctx:betweenat(a, b, t) - returns new var at t[a,b]
		l.PushGoFunction(func(l *lua.State) int {
			a := checkVar(l, 2)
			b := checkVar(l, 3)
			t := lua.CheckNumber(l, 4)
			p := ctx.BetweenAt(a, b, t)
			pushVar(l, p)
			return 1
		})
		return 1

	case "collinear":
		// ctx:collinear(p, a, b) - constrain p on line through a and b
		l.PushGoFunction(func(l *lua.State) int {
			p := checkVar(l, 2)
			a := checkVar(l, 3)
			b := checkVar(l, 4)
			ctx.Collinear(p, a, b)
			l.PushValue(1)
			return 1
		})
		return 1

	case "intersection":
		// ctx:intersection(p, a1, a2, b1, b2) - constrain p = intersection of lines
		l.PushGoFunction(func(l *lua.State) int {
			p := checkVar(l, 2)
			a1 := checkVar(l, 3)
			a2 := checkVar(l, 4)
			b1 := checkVar(l, 5)
			b2 := checkVar(l, 6)
			err := ctx.Intersection(p, a1, a2, b1, b2)
			if err != nil {
				lua.Errorf(l, "intersection error: %s", err.Error())
				return 0
			}
			l.PushValue(1)
			return 1
		})
		return 1

	case "intersectionof":
		// ctx:intersectionof(a1, a2, b1, b2) - returns new var at intersection
		l.PushGoFunction(func(l *lua.State) int {
			a1 := checkVar(l, 2)
			a2 := checkVar(l, 3)
			b1 := checkVar(l, 4)
			b2 := checkVar(l, 5)
			p, err := ctx.IntersectionOf(a1, a2, b1, b2)
			if err != nil {
				lua.Errorf(l, "intersection error: %s", err.Error())
				return 0
			}
			pushVar(l, p)
			return 1
		})
		return 1

	case "sum":
		// ctx:sum(result, a, b) - constrain result = a + b
		l.PushGoFunction(func(l *lua.State) int {
			result := checkVar(l, 2)
			a := checkVar(l, 3)
			b := checkVar(l, 4)
			ctx.Sum(result, a, b)
			l.PushValue(1)
			return 1
		})
		return 1

	case "diff":
		// ctx:diff(result, a, b) - constrain result = a - b
		l.PushGoFunction(func(l *lua.State) int {
			result := checkVar(l, 2)
			a := checkVar(l, 3)
			b := checkVar(l, 4)
			ctx.Diff(result, a, b)
			l.PushValue(1)
			return 1
		})
		return 1

	case "scaled":
		// ctx:scaled(result, v, t) - constrain result = t * v
		l.PushGoFunction(func(l *lua.State) int {
			result := checkVar(l, 2)
			v := checkVar(l, 3)
			t := lua.CheckNumber(l, 4)
			ctx.Scaled(result, v, t)
			l.PushValue(1)
			return 1
		})
		return 1

	case "solve":
		// ctx:solve() - solve the equation system
		l.PushGoFunction(func(l *lua.State) int {
			err := ctx.Solve()
			if err != nil {
				lua.Errorf(l, "solve error: %s", err.Error())
				return 0
			}
			l.PushValue(1)
			return 1
		})
		return 1

	case "path":
		// ctx:path() - create a PathBuilder linked to this context
		l.PushGoFunction(func(l *lua.State) int {
			pb := &ContextPathBuilder{
				builder: ctx.NewPath(),
				ctx:     ctx,
			}
			l.PushUserData(pb)
			lua.SetMetaTableNamed(l, "hobby.ctxpathbuilder")
			return 1
		})
		return 1
	}

	return 0
}

func varIndex(l *lua.State) int {
	v := checkVar(l, 1)
	key := lua.CheckString(l, 2)

	switch key {
	case "x":
		// v.x or v:x() - get x coordinate
		l.PushNumber(v.X())
		return 1

	case "y":
		// v.y or v:y() - get y coordinate
		l.PushNumber(v.Y())
		return 1

	case "xy":
		// v:xy() - returns x, y
		l.PushGoFunction(func(l *lua.State) int {
			x, y := v.XY()
			l.PushNumber(x)
			l.PushNumber(y)
			return 2
		})
		return 1

	case "point":
		// v:point() - returns as hobby point
		l.PushGoFunction(func(l *lua.State) int {
			pushPoint(l, v.Point())
			return 1
		})
		return 1

	case "setx":
		// v:setx(x) - set x coordinate
		l.PushGoFunction(func(l *lua.State) int {
			x := lua.CheckNumber(l, 2)
			v.SetX(x)
			l.PushValue(1)
			return 1
		})
		return 1

	case "sety":
		// v:sety(y) - set y coordinate
		l.PushGoFunction(func(l *lua.State) int {
			y := lua.CheckNumber(l, 2)
			v.SetY(y)
			l.PushValue(1)
			return 1
		})
		return 1

	case "setxy":
		// v:setxy(x, y) - set both coordinates
		l.PushGoFunction(func(l *lua.State) int {
			x := lua.CheckNumber(l, 2)
			y := lua.CheckNumber(l, 3)
			v.SetXY(x, y)
			l.PushValue(1)
			return 1
		})
		return 1
	}

	return 0
}

// ContextPathBuilder wraps a PathBuilder with context for variable references
type ContextPathBuilder struct {
	builder *draw.PathBuilder
	ctx     *draw.Context
}

// registerCtxPathBuilderMeta registers the metatable for context-linked path builders
func registerCtxPathBuilderMeta(l *lua.State) {
	lua.NewMetaTable(l, "hobby.ctxpathbuilder")
	l.PushGoFunction(ctxPathBuilderIndex)
	l.SetField(-2, "__index")
	l.Pop(1)
}

func ctxPathBuilderIndex(l *lua.State) int {
	cpb := l.ToUserData(1).(*ContextPathBuilder)
	pb := cpb.builder
	key := lua.CheckString(l, 2)

	switch key {
	case "movetovar":
		// pb:movetovar(var) - move to a context variable
		l.PushGoFunction(func(l *lua.State) int {
			v := checkVar(l, 2)
			pb.MoveToVar(v)
			l.PushValue(1)
			return 1
		})
		return 1

	case "linetovar":
		// pb:linetovar(var) - line to a context variable
		l.PushGoFunction(func(l *lua.State) int {
			v := checkVar(l, 2)
			pb.LineToVar(v)
			l.PushValue(1)
			return 1
		})
		return 1

	case "curvetovar":
		// pb:curvetovar(var) - curve to a context variable
		l.PushGoFunction(func(l *lua.State) int {
			v := checkVar(l, 2)
			pb.CurveToVar(v)
			l.PushValue(1)
			return 1
		})
		return 1

	case "moveto":
		l.PushGoFunction(func(l *lua.State) int {
			p := checkPoint(l, 2)
			pb.MoveTo(p)
			l.PushValue(1)
			return 1
		})
		return 1

	case "lineto":
		l.PushGoFunction(func(l *lua.State) int {
			p := checkPoint(l, 2)
			pb.LineTo(p)
			l.PushValue(1)
			return 1
		})
		return 1

	case "curveto":
		l.PushGoFunction(func(l *lua.State) int {
			p := checkPoint(l, 2)
			pb.CurveTo(p)
			l.PushValue(1)
			return 1
		})
		return 1

	case "curvetowithcontrols":
		l.PushGoFunction(func(l *lua.State) int {
			pt := checkPoint(l, 2)
			c1 := checkPoint(l, 3)
			c2 := checkPoint(l, 4)
			pb.CurveToWithControls(pt, c1, c2)
			l.PushValue(1)
			return 1
		})
		return 1

	case "dir":
		l.PushGoFunction(func(l *lua.State) int {
			angle := lua.CheckNumber(l, 2)
			pb.WithDirection(angle)
			l.PushValue(1)
			return 1
		})
		return 1

	case "indir":
		l.PushGoFunction(func(l *lua.State) int {
			angle := lua.CheckNumber(l, 2)
			pb.WithIncomingDirection(angle)
			l.PushValue(1)
			return 1
		})
		return 1

	case "curl":
		l.PushGoFunction(func(l *lua.State) int {
			c := lua.CheckNumber(l, 2)
			pb.WithCurl(c)
			l.PushValue(1)
			return 1
		})
		return 1

	case "outcurl":
		l.PushGoFunction(func(l *lua.State) int {
			c := lua.CheckNumber(l, 2)
			pb.WithOutgoingCurl(c)
			l.PushValue(1)
			return 1
		})
		return 1

	case "incurl":
		l.PushGoFunction(func(l *lua.State) int {
			c := lua.CheckNumber(l, 2)
			pb.WithIncomingCurl(c)
			l.PushValue(1)
			return 1
		})
		return 1

	case "tension":
		l.PushGoFunction(func(l *lua.State) int {
			t := lua.CheckNumber(l, 2)
			pb.WithTension(t)
			l.PushValue(1)
			return 1
		})
		return 1

	case "outtension":
		l.PushGoFunction(func(l *lua.State) int {
			t := lua.CheckNumber(l, 2)
			pb.WithOutgoingTension(t)
			l.PushValue(1)
			return 1
		})
		return 1

	case "intension":
		l.PushGoFunction(func(l *lua.State) int {
			t := lua.CheckNumber(l, 2)
			pb.WithIncomingTension(t)
			l.PushValue(1)
			return 1
		})
		return 1

	case "tensionatleast":
		l.PushGoFunction(func(l *lua.State) int {
			t := lua.CheckNumber(l, 2)
			pb.WithTensionAtLeast(t)
			l.PushValue(1)
			return 1
		})
		return 1

	case "tensioninfinity":
		l.PushGoFunction(func(l *lua.State) int {
			pb.WithTensionInfinity()
			l.PushValue(1)
			return 1
		})
		return 1

	case "stroke":
		l.PushGoFunction(func(l *lua.State) int {
			c := checkColor(l, 2)
			pb.WithStrokeColor(c)
			l.PushValue(1)
			return 1
		})
		return 1

	case "strokewidth":
		l.PushGoFunction(func(l *lua.State) int {
			w := lua.CheckNumber(l, 2)
			pb.WithStrokeWidth(w)
			l.PushValue(1)
			return 1
		})
		return 1

	case "fill":
		l.PushGoFunction(func(l *lua.State) int {
			c := checkColor(l, 2)
			pb.WithFill(c)
			l.PushValue(1)
			return 1
		})
		return 1

	case "pen":
		l.PushGoFunction(func(l *lua.State) int {
			pen := checkPen(l, 2)
			pb.WithPen(pen)
			l.PushValue(1)
			return 1
		})
		return 1

	case "dash":
		l.PushGoFunction(func(l *lua.State) int {
			dash := checkDash(l, 2)
			pb.WithDashPattern(dash)
			l.PushValue(1)
			return 1
		})
		return 1

	case "evenly":
		l.PushGoFunction(func(l *lua.State) int {
			pb.DashedEvenly()
			l.PushValue(1)
			return 1
		})
		return 1

	case "withdots":
		l.PushGoFunction(func(l *lua.State) int {
			pb.DashedWithDots()
			l.PushValue(1)
			return 1
		})
		return 1

	case "arrow":
		l.PushGoFunction(func(l *lua.State) int {
			pb.WithArrow()
			l.PushValue(1)
			return 1
		})
		return 1

	case "dblarrow":
		l.PushGoFunction(func(l *lua.State) int {
			pb.WithDoubleArrow()
			l.PushValue(1)
			return 1
		})
		return 1

	case "arrowstyle":
		l.PushGoFunction(func(l *lua.State) int {
			length := lua.CheckNumber(l, 2)
			angle := lua.CheckNumber(l, 3)
			pb.WithArrowStyle(length, angle)
			l.PushValue(1)
			return 1
		})
		return 1

	case "linejoin":
		l.PushGoFunction(func(l *lua.State) int {
			join := lua.CheckString(l, 2)
			switch join {
			case "miter":
				pb.WithLineJoin(mp.LineJoinMiter)
			case "round":
				pb.WithLineJoin(mp.LineJoinRound)
			case "bevel":
				pb.WithLineJoin(mp.LineJoinBevel)
			default:
				lua.Errorf(l, "unknown linejoin: %s", join)
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
				pb.WithLineCap(mp.LineCapButt)
			case "round":
				pb.WithLineCap(mp.LineCapRounded)
			case "square":
				pb.WithLineCap(mp.LineCapSquared)
			default:
				lua.Errorf(l, "unknown linecap: %s", cap)
			}
			l.PushValue(1)
			return 1
		})
		return 1

	case "close", "cycle":
		l.PushGoFunction(func(l *lua.State) int {
			pb.Close()
			l.PushValue(1)
			return 1
		})
		return 1

	case "build":
		l.PushGoFunction(func(l *lua.State) int {
			path, err := pb.Solve()
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
