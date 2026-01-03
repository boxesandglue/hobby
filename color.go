package hobby

import (
	"github.com/boxesandglue/mpgo/mp"
	lua "github.com/speedata/go-lua"
)

// checkColor extracts a color from a Lua argument.
// Accepts:
//   - string: CSS color ("red", "#ff0000", "rgb(255,0,0)")
//   - table with r,g,b fields (0-1 range)
//   - table with 3 numbers {r, g, b} (0-1 range)
//   - 3 numbers on the stack starting at index (r, g, b in 0-1 range)
func checkColor(l *lua.State, index int) mp.Color {
	// Check for colorWrapper userdata first
	if l.IsUserData(index) {
		if cw, ok := l.ToUserData(index).(*colorWrapper); ok {
			return cw.color
		}
	}

	switch {
	case l.IsString(index):
		// CSS color string
		css := lua.CheckString(l, index)
		return mp.ColorCSS(css)

	case l.IsTable(index):
		// Try {r=, g=, b=} first
		l.Field(index, "r")
		if !l.IsNil(-1) {
			r, _ := l.ToNumber(-1)
			l.Pop(1)
			l.Field(index, "g")
			g, _ := l.ToNumber(-1)
			l.Pop(1)
			l.Field(index, "b")
			b, _ := l.ToNumber(-1)
			l.Pop(1)
			return mp.ColorRGB(r, g, b)
		}
		l.Pop(1)

		// Try array-style {r, g, b}
		l.RawGetInt(index, 1)
		r, _ := l.ToNumber(-1)
		l.Pop(1)
		l.RawGetInt(index, 2)
		g, _ := l.ToNumber(-1)
		l.Pop(1)
		l.RawGetInt(index, 3)
		b, _ := l.ToNumber(-1)
		l.Pop(1)
		return mp.ColorRGB(r, g, b)

	case l.IsNumber(index):
		// Three separate numbers: r, g, b
		r := lua.CheckNumber(l, index)
		g := lua.CheckNumber(l, index+1)
		b := lua.CheckNumber(l, index+2)
		return mp.ColorRGB(r, g, b)
	}

	lua.Errorf(l, "expected color at argument %d", index)
	return mp.Color{}
}

// luaColorRGB creates an RGB color: hobby.rgb(r, g, b)
func luaColorRGB(l *lua.State) int {
	r := lua.CheckNumber(l, 1)
	g := lua.CheckNumber(l, 2)
	b := lua.CheckNumber(l, 3)
	pushColor(l, mp.ColorRGB(r, g, b))
	return 1
}

// luaColorGray creates a grayscale color: hobby.gray(g)
func luaColorGray(l *lua.State) int {
	g := lua.CheckNumber(l, 1)
	pushColor(l, mp.ColorGray(g))
	return 1
}

// luaColorCSS creates a color from CSS string: hobby.color("red")
func luaColorCSS(l *lua.State) int {
	css := lua.CheckString(l, 1)
	pushColor(l, mp.ColorCSS(css))
	return 1
}

// Color userdata wrapper
type colorWrapper struct {
	color mp.Color
}

func pushColor(l *lua.State, c mp.Color) {
	l.PushUserData(&colorWrapper{color: c})
	lua.SetMetaTableNamed(l, "hobby.color")
}

func registerColorMeta(l *lua.State) {
	lua.NewMetaTable(l, "hobby.color")
	l.PushGoFunction(colorToString)
	l.SetField(-2, "__tostring")
	l.Pop(1)
}

func colorToString(l *lua.State) int {
	cw := l.ToUserData(1).(*colorWrapper)
	l.PushString(cw.color.CSS())
	return 1
}
