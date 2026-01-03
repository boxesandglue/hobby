package hobby

import (
	"bytes"
	"os"

	"github.com/boxesandglue/mpgo/svg"
	lua "github.com/speedata/go-lua"
)

// luaNewSVG creates a new SVG builder: hobby.svg()
func luaNewSVG(l *lua.State) int {
	builder := svg.NewBuilder()
	pushSVG(l, builder)
	return 1
}

// registerSVGMeta registers the metatable for SVG builders
func registerSVGMeta(l *lua.State) {
	lua.NewMetaTable(l, "hobby.svg")
	l.PushGoFunction(svgIndex)
	l.SetField(-2, "__index")
	l.Pop(1)
}

func svgIndex(l *lua.State) int {
	s := checkSVG(l, 1)
	key := lua.CheckString(l, 2)

	switch key {
	case "add":
		l.PushGoFunction(func(l *lua.State) int {
			path := checkPath(l, 2)
			s.AddPathFromPath(path)
			l.PushValue(1) // return self for chaining
			return 1
		})
		return 1

	case "padding":
		l.PushGoFunction(func(l *lua.State) int {
			p := lua.CheckNumber(l, 2)
			s.Padding(p)
			l.PushValue(1)
			return 1
		})
		return 1

	case "write":
		l.PushGoFunction(func(l *lua.State) int {
			filename := lua.CheckString(l, 2)
			f, err := os.Create(filename)
			if err != nil {
				lua.Errorf(l, "cannot create file: %s", err.Error())
				return 0
			}
			defer f.Close()
			s.WriteTo(f)
			return 0
		})
		return 1

	case "tostring":
		l.PushGoFunction(func(l *lua.State) int {
			var buf bytes.Buffer
			s.WriteTo(&buf)
			l.PushString(buf.String())
			return 1
		})
		return 1

	case "addpicture":
		l.PushGoFunction(func(l *lua.State) int {
			pic := checkPicture(l, 2)
			s.AddPicture(pic)
			l.PushValue(1) // return self for chaining
			return 1
		})
		return 1
	}

	return 0
}
