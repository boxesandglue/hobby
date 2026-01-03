package hobby

import (
	"github.com/boxesandglue/mpgo/draw"
	lua "github.com/speedata/go-lua"
)

// luaNewPicture creates a new Picture: h.picture()
func luaNewPicture(l *lua.State) int {
	pic := draw.NewPicture()
	pushPicture(l, pic)
	return 1
}

// registerPictureMeta registers the metatable for Pictures
func registerPictureMeta(l *lua.State) {
	lua.NewMetaTable(l, "hobby.picture")
	l.PushGoFunction(pictureIndex)
	l.SetField(-2, "__index")
	l.Pop(1)
}

// pushPicture pushes a Picture as userdata
func pushPicture(l *lua.State, p *draw.Picture) {
	l.PushUserData(p)
	lua.SetMetaTableNamed(l, "hobby.picture")
}

// checkPicture checks if value at index is a Picture
func checkPicture(l *lua.State, index int) *draw.Picture {
	ud := l.ToUserData(index)
	if p, ok := ud.(*draw.Picture); ok {
		return p
	}
	lua.Errorf(l, "expected picture at argument %d", index)
	return nil
}

func pictureIndex(l *lua.State) int {
	pic := checkPicture(l, 1)
	key := lua.CheckString(l, 2)

	switch key {
	case "add":
		// pic:add(path) - add a path to the picture
		l.PushGoFunction(func(l *lua.State) int {
			path := checkPath(l, 2)
			pic.AddPath(path)
			l.PushValue(1) // return self for chaining
			return 1
		})
		return 1

	case "addpicture":
		// pic:addpicture(other) - add all paths from another picture
		l.PushGoFunction(func(l *lua.State) int {
			other := checkPicture(l, 2)
			pic.AddPicture(other)
			l.PushValue(1) // return self for chaining
			return 1
		})
		return 1

	case "clip":
		// pic:clip(clipPath) - set clipping path
		l.PushGoFunction(func(l *lua.State) int {
			clipPath := checkPath(l, 2)
			pic.Clip(clipPath)
			l.PushValue(1) // return self for chaining
			return 1
		})
		return 1

	case "clippath":
		// pic:clippath() - get current clipping path
		l.PushGoFunction(func(l *lua.State) int {
			cp := pic.ClipPath()
			if cp == nil {
				l.PushNil()
			} else {
				pushPath(l, cp)
			}
			return 1
		})
		return 1

	case "paths":
		// pic:paths() - get all paths as a table
		l.PushGoFunction(func(l *lua.State) int {
			paths := pic.Paths()
			l.CreateTable(len(paths), 0)
			for i, p := range paths {
				pushPath(l, p)
				l.RawSetInt(-2, i+1)
			}
			return 1
		})
		return 1
	}

	return 0
}
