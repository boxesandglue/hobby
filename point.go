package hobby

import (
	"fmt"

	"github.com/boxesandglue/mpgo/mp"
	lua "github.com/speedata/go-lua"
)

// luaPoint creates a new point: hobby.point(x, y)
func luaPoint(l *lua.State) int {
	x := lua.CheckNumber(l, 1)
	y := lua.CheckNumber(l, 2)
	pushPoint(l, mp.P(x, y))
	return 1
}

// luaDir creates a unit vector at angle: hobby.dir(angle)
func luaDir(l *lua.State) int {
	angle := lua.CheckNumber(l, 1)
	pushPoint(l, mp.Dir(angle))
	return 1
}

// luaMidPoint returns the midpoint of two points: hobby.midpoint(a, b)
func luaMidPoint(l *lua.State) int {
	a := checkPoint(l, 1)
	b := checkPoint(l, 2)
	pushPoint(l, mp.MidPoint(a, b))
	return 1
}

// luaDistance returns the distance between two points: hobby.distance(a, b)
func luaDistance(l *lua.State) int {
	a := checkPoint(l, 1)
	b := checkPoint(l, 2)
	l.PushNumber(mp.Distance(a, b))
	return 1
}

// luaBetween returns point interpolated between a and b: hobby.between(a, b, t)
// Equivalent to MetaPost's t[a,b]
func luaBetween(l *lua.State) int {
	a := checkPoint(l, 1)
	b := checkPoint(l, 2)
	t := lua.CheckNumber(l, 3)
	pushPoint(l, mp.PointBetween(a, b, t))
	return 1
}

// registerPointMeta registers the metatable for points
func registerPointMeta(l *lua.State) {
	lua.NewMetaTable(l, "hobby.point")

	// __index: access x, y fields and methods
	l.PushGoFunction(pointIndex)
	l.SetField(-2, "__index")

	// __add: point + point
	l.PushGoFunction(pointAdd)
	l.SetField(-2, "__add")

	// __sub: point - point
	l.PushGoFunction(pointSub)
	l.SetField(-2, "__sub")

	// __mul: point * number or number * point
	l.PushGoFunction(pointMul)
	l.SetField(-2, "__mul")

	// __unm: -point
	l.PushGoFunction(pointUnm)
	l.SetField(-2, "__unm")

	// __tostring
	l.PushGoFunction(pointToString)
	l.SetField(-2, "__tostring")

	l.Pop(1)
}

func pointIndex(l *lua.State) int {
	p := checkPoint(l, 1)
	key := lua.CheckString(l, 2)

	switch key {
	case "x":
		l.PushNumber(p.X)
		return 1
	case "y":
		l.PushNumber(p.Y)
		return 1
	case "length":
		l.PushNumber(p.Length())
		return 1
	case "angle":
		l.PushNumber(p.Angle())
		return 1
	case "normalized":
		pushPoint(l, p.Normalized())
		return 1
	}

	return 0
}

func pointAdd(l *lua.State) int {
	a := checkPoint(l, 1)
	b := checkPoint(l, 2)
	pushPoint(l, a.Add(b))
	return 1
}

func pointSub(l *lua.State) int {
	a := checkPoint(l, 1)
	b := checkPoint(l, 2)
	pushPoint(l, a.Sub(b))
	return 1
}

func pointMul(l *lua.State) int {
	// Handle both point * number and number * point
	if l.IsNumber(1) {
		n := lua.CheckNumber(l, 1)
		p := checkPoint(l, 2)
		pushPoint(l, p.Mul(n))
	} else {
		p := checkPoint(l, 1)
		n := lua.CheckNumber(l, 2)
		pushPoint(l, p.Mul(n))
	}
	return 1
}

func pointUnm(l *lua.State) int {
	p := checkPoint(l, 1)
	pushPoint(l, p.Mul(-1))
	return 1
}

func pointToString(l *lua.State) int {
	p := checkPoint(l, 1)
	l.PushString(fmt.Sprintf("(%g, %g)", p.X, p.Y))
	return 1
}
