// Package hobby provides Lua bindings for the mpgo MetaPost curve library.
package hobby

import (
	"github.com/boxesandglue/mpgo/mp"
	"github.com/boxesandglue/mpgo/svg"
	lua "github.com/speedata/go-lua"
)

// engine is the shared curve-solving engine.
var engine = mp.NewEngine()

// Open registers the hobby module for require('hobby').
// Call this after lua.OpenLibraries(l).
func Open(l *lua.State) {
	// Register metatables first (needed before any hobby objects are created)
	registerPointMeta(l)
	registerPathMeta(l)
	registerSVGMeta(l)
	registerColorMeta(l)
	registerPenMeta(l)
	registerDashMeta(l)
	registerPictureMeta(l)
	registerLabelMeta(l)
	registerContextMeta(l)
	registerVarMeta(l)
	registerCtxPathBuilderMeta(l)
	registerFaceMeta(l)

	// Register the module loader in package.preload
	l.Field(lua.RegistryIndex, "_PRELOAD")
	l.PushGoFunction(loader)
	l.SetField(-2, "hobby")
	l.Pop(1)
}

// loader is the module loader function called by require('hobby')
func loader(l *lua.State) int {
	// Create the "hobby" table
	l.NewTable()

	// Point functions
	l.PushGoFunction(luaPoint)
	l.SetField(-2, "point")

	l.PushGoFunction(luaDir)
	l.SetField(-2, "dir")

	l.PushGoFunction(luaMidPoint)
	l.SetField(-2, "midpoint")

	l.PushGoFunction(luaDistance)
	l.SetField(-2, "distance")

	l.PushGoFunction(luaBetween)
	l.SetField(-2, "between")

	// Path builder
	l.PushGoFunction(luaNewPath)
	l.SetField(-2, "path")

	// Predefined paths
	l.PushGoFunction(luaFullCircle)
	l.SetField(-2, "fullcircle")

	l.PushGoFunction(luaHalfCircle)
	l.SetField(-2, "halfcircle")

	l.PushGoFunction(luaQuarterCircle)
	l.SetField(-2, "quartercircle")

	l.PushGoFunction(luaUnitSquare)
	l.SetField(-2, "unitsquare")

	// SVG output
	l.PushGoFunction(luaNewSVG)
	l.SetField(-2, "svg")

	// Color constructors
	l.PushGoFunction(luaColorRGB)
	l.SetField(-2, "rgb")

	l.PushGoFunction(luaColorGray)
	l.SetField(-2, "gray")

	l.PushGoFunction(luaColorCSS)
	l.SetField(-2, "color")

	// Pen constructors
	l.PushGoFunction(luaPenCircle)
	l.SetField(-2, "pencircle")

	l.PushGoFunction(luaPenSquare)
	l.SetField(-2, "pensquare")

	l.PushGoFunction(luaPenRazor)
	l.SetField(-2, "penrazor")

	l.PushGoFunction(luaPenSpeck)
	l.SetField(-2, "penspeck")

	// Dash pattern constructors
	l.PushGoFunction(luaDashEvenly)
	l.SetField(-2, "evenly")

	l.PushGoFunction(luaDashWithDots)
	l.SetField(-2, "withdots")

	l.PushGoFunction(luaDashPattern)
	l.SetField(-2, "dashed")

	// Path operations
	l.PushGoFunction(luaBuildCycle)
	l.SetField(-2, "buildcycle")

	// Picture
	l.PushGoFunction(luaNewPicture)
	l.SetField(-2, "picture")

	// Context (equation solver)
	l.PushGoFunction(luaNewContext)
	l.SetField(-2, "context")

	// Font loading
	l.PushGoFunction(luaLoadFont)
	l.SetField(-2, "loadfont")

	// Return the table (it's on top of the stack)
	return 1
}

// Helper to get a Point from a Lua value (table with x, y or userdata)
func checkPoint(l *lua.State, index int) mp.Point {
	if l.IsUserData(index) {
		if p, ok := l.ToUserData(index).(*mp.Point); ok {
			return *p
		}
	}
	// Try table with x, y fields
	l.Field(index, "x")
	x := lua.CheckNumber(l, -1)
	l.Pop(1)
	l.Field(index, "y")
	y := lua.CheckNumber(l, -1)
	l.Pop(1)
	return mp.P(x, y)
}

// Helper to push a Point as userdata
func pushPoint(l *lua.State, p mp.Point) {
	ptr := new(mp.Point)
	*ptr = p
	l.PushUserData(ptr)
	lua.SetMetaTableNamed(l, "hobby.point")
}

// Helper to push a Path as userdata
func pushPath(l *lua.State, p *mp.Path) {
	l.PushUserData(p)
	lua.SetMetaTableNamed(l, "hobby.path")
}

// Helper to check if value at index is a Path
func checkPath(l *lua.State, index int) *mp.Path {
	ud := l.ToUserData(index)
	if p, ok := ud.(*mp.Path); ok {
		return p
	}
	lua.Errorf(l, "expected path at argument %d", index)
	return nil
}

// Helper to push an SVG builder as userdata
func pushSVG(l *lua.State, s *svg.Builder) {
	l.PushUserData(s)
	lua.SetMetaTableNamed(l, "hobby.svg")
}

// Helper to check if value at index is an SVG builder
func checkSVG(l *lua.State, index int) *svg.Builder {
	ud := l.ToUserData(index)
	if s, ok := ud.(*svg.Builder); ok {
		return s
	}
	lua.Errorf(l, "expected svg at argument %d", index)
	return nil
}
