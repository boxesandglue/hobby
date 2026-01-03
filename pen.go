package hobby

import (
	"github.com/boxesandglue/mpgo/mp"
	lua "github.com/speedata/go-lua"
)

// Pen functions

// luaPenCircle creates a circular pen: hobby.pencircle(diameter)
func luaPenCircle(l *lua.State) int {
	d := float64(1) // default diameter
	if l.Top() >= 1 {
		d = lua.CheckNumber(l, 1)
	}
	pushPen(l, mp.PenCircle(mp.Number(d)))
	return 1
}

// luaPenSquare creates a square pen: hobby.pensquare(size)
func luaPenSquare(l *lua.State) int {
	size := float64(1) // default size
	if l.Top() >= 1 {
		size = lua.CheckNumber(l, 1)
	}
	pushPen(l, mp.PenSquare(mp.Number(size)))
	return 1
}

// luaPenRazor creates a razor pen: hobby.penrazor(size, angle?)
func luaPenRazor(l *lua.State) int {
	size := float64(1)
	if l.Top() >= 1 {
		size = lua.CheckNumber(l, 1)
	}
	if l.Top() >= 2 {
		angle := lua.CheckNumber(l, 2)
		pushPen(l, mp.PenRazorRotated(mp.Number(size), mp.Number(angle)))
	} else {
		pushPen(l, mp.PenRazor(mp.Number(size)))
	}
	return 1
}

// luaPenSpeck creates a speck pen (minimal size): hobby.penspeck()
func luaPenSpeck(l *lua.State) int {
	pushPen(l, mp.PenSpeck())
	return 1
}

// Pen wrapper and metatable

type penWrapper struct {
	pen *mp.Pen
}

func pushPen(l *lua.State, p *mp.Pen) {
	l.PushUserData(&penWrapper{pen: p})
	lua.SetMetaTableNamed(l, "hobby.pen")
}

func checkPen(l *lua.State, index int) *mp.Pen {
	if l.IsUserData(index) {
		if pw, ok := l.ToUserData(index).(*penWrapper); ok {
			return pw.pen
		}
	}
	lua.Errorf(l, "expected pen at argument %d", index)
	return nil
}

func registerPenMeta(l *lua.State) {
	lua.NewMetaTable(l, "hobby.pen")
	l.PushGoFunction(penIndex)
	l.SetField(-2, "__index")
	l.Pop(1)
}

func penIndex(l *lua.State) int {
	_ = l.ToUserData(1).(*penWrapper)
	key := lua.CheckString(l, 2)

	switch key {
	case "elliptical":
		pw := l.ToUserData(1).(*penWrapper)
		l.PushBoolean(pw.pen.Elliptical)
		return 1
	}

	return 0
}

// Dash pattern functions

// luaDashEvenly creates the standard "evenly" pattern: hobby.evenly()
func luaDashEvenly(l *lua.State) int {
	pushDash(l, mp.DashEvenly())
	return 1
}

// luaDashWithDots creates the "withdots" pattern: hobby.withdots()
func luaDashWithDots(l *lua.State) int {
	pushDash(l, mp.DashWithDots())
	return 1
}

// luaDashPattern creates a custom pattern: hobby.dashed(on, off, ...)
func luaDashPattern(l *lua.State) int {
	n := l.Top()
	if n == 0 {
		l.PushNil()
		return 1
	}
	onOff := make([]float64, n)
	for i := 1; i <= n; i++ {
		onOff[i-1] = lua.CheckNumber(l, i)
	}
	pushDash(l, mp.NewDashPattern(onOff...))
	return 1
}

// Dash wrapper and metatable

type dashWrapper struct {
	dash *mp.DashPattern
}

func pushDash(l *lua.State, d *mp.DashPattern) {
	l.PushUserData(&dashWrapper{dash: d})
	lua.SetMetaTableNamed(l, "hobby.dash")
}

func checkDash(l *lua.State, index int) *mp.DashPattern {
	if l.IsUserData(index) {
		if dw, ok := l.ToUserData(index).(*dashWrapper); ok {
			return dw.dash
		}
	}
	lua.Errorf(l, "expected dash pattern at argument %d", index)
	return nil
}

func registerDashMeta(l *lua.State) {
	lua.NewMetaTable(l, "hobby.dash")
	l.PushGoFunction(dashIndex)
	l.SetField(-2, "__index")
	l.Pop(1)
}

func dashIndex(l *lua.State) int {
	dw := l.ToUserData(1).(*dashWrapper)
	key := lua.CheckString(l, 2)

	switch key {
	case "scaled":
		l.PushGoFunction(func(l *lua.State) int {
			s := lua.CheckNumber(l, 2)
			scaled := dw.dash.Scaled(s)
			pushDash(l, scaled)
			return 1
		})
		return 1

	case "shifted":
		l.PushGoFunction(func(l *lua.State) int {
			offset := lua.CheckNumber(l, 2)
			shifted := dw.dash.Shifted(offset)
			pushDash(l, shifted)
			return 1
		})
		return 1
	}

	return 0
}
