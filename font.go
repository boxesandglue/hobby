package hobby

import (
	"os"

	"github.com/boxesandglue/mpgo/font"
	lua "github.com/speedata/go-lua"
)

// luaLoadFont loads a font file: h.loadfont("path/to/font.ttf")
func luaLoadFont(l *lua.State) int {
	path := lua.CheckString(l, 1)

	f, err := os.Open(path)
	if err != nil {
		lua.Errorf(l, "failed to open font %q: %s", path, err)
		return 0
	}
	defer f.Close()

	face, err := font.Load(f)
	if err != nil {
		lua.Errorf(l, "failed to load font %q: %s", path, err)
		return 0
	}

	pushFace(l, face)
	return 1
}

// registerFaceMeta registers the metatable for font faces
func registerFaceMeta(l *lua.State) {
	lua.NewMetaTable(l, "hobby.face")
	l.Pop(1)
}

// pushFace pushes a font.Face as userdata
func pushFace(l *lua.State, face *font.Face) {
	l.PushUserData(face)
	lua.SetMetaTableNamed(l, "hobby.face")
}

// checkFace checks if value at index is a font.Face
func checkFace(l *lua.State, index int) *font.Face {
	ud := l.ToUserData(index)
	if f, ok := ud.(*font.Face); ok {
		return f
	}
	lua.Errorf(l, "expected font face at argument %d", index)
	return nil
}
