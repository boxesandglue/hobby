package hobby

import (
	"github.com/boxesandglue/mpgo/draw"
	"github.com/boxesandglue/mpgo/mp"
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

// parseAnchor converts an anchor string to mp.Anchor
func parseAnchor(s string) mp.Anchor {
	switch s {
	case "center", "c":
		return mp.AnchorCenter
	case "left", "lft":
		return mp.AnchorLeft
	case "right", "rt":
		return mp.AnchorRight
	case "top":
		return mp.AnchorTop
	case "bottom", "bot":
		return mp.AnchorBottom
	case "upperleft", "ulft":
		return mp.AnchorUpperLeft
	case "upperright", "urt":
		return mp.AnchorUpperRight
	case "lowerleft", "llft":
		return mp.AnchorLowerLeft
	case "lowerright", "lrt":
		return mp.AnchorLowerRight
	default:
		return mp.AnchorCenter
	}
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

	case "label":
		// pic:label(text, point, anchor) - add a text label
		// anchor is a string: "center", "top", "bot", "lft", "rt", "ulft", "urt", "llft", "lrt"
		l.PushGoFunction(func(l *lua.State) int {
			text := lua.CheckString(l, 2)
			pos := checkPoint(l, 3)
			anchorStr := lua.OptString(l, 4, "center")
			anchor := parseAnchor(anchorStr)
			pic.Label(text, pos, anchor)
			l.PushValue(1) // return self for chaining
			return 1
		})
		return 1

	case "dotlabel":
		// pic:dotlabel(text, point, anchor, color) - add a label with a dot
		l.PushGoFunction(func(l *lua.State) int {
			text := lua.CheckString(l, 2)
			pos := checkPoint(l, 3)
			anchorStr := lua.OptString(l, 4, "center")
			anchor := parseAnchor(anchorStr)
			// Color is optional, default to black
			var color mp.Color
			if l.Top() >= 5 && !l.IsNil(5) {
				color = checkColor(l, 5)
			} else {
				color = mp.ColorCSS("black")
			}
			pic.DotLabel(text, pos, anchor, color)
			l.PushValue(1) // return self for chaining
			return 1
		})
		return 1

	case "labels":
		// pic:labels() - get all labels as a table
		l.PushGoFunction(func(l *lua.State) int {
			labels := pic.Labels()
			l.CreateTable(len(labels), 0)
			for i, lbl := range labels {
				pushLabel(l, lbl)
				l.RawSetInt(-2, i+1)
			}
			return 1
		})
		return 1
	}

	return 0
}

// pushLabel pushes a Label as userdata
func pushLabel(l *lua.State, label *mp.Label) {
	l.PushUserData(label)
	lua.SetMetaTableNamed(l, "hobby.label")
}

// checkLabel checks if value at index is a Label
func checkLabel(l *lua.State, index int) *mp.Label {
	ud := l.ToUserData(index)
	if lbl, ok := ud.(*mp.Label); ok {
		return lbl
	}
	lua.Errorf(l, "expected label at argument %d", index)
	return nil
}

// registerLabelMeta registers the metatable for Labels
func registerLabelMeta(l *lua.State) {
	lua.NewMetaTable(l, "hobby.label")
	l.PushGoFunction(labelIndex)
	l.SetField(-2, "__index")
	l.Pop(1)
}

func labelIndex(l *lua.State) int {
	label := checkLabel(l, 1)
	key := lua.CheckString(l, 2)

	switch key {
	case "text":
		l.PushString(label.Text)
		return 1

	case "position":
		l.PushGoFunction(func(l *lua.State) int {
			pushPoint(l, label.Position)
			return 1
		})
		return 1

	case "fontsize":
		l.PushNumber(label.FontSize)
		return 1

	case "setfontsize":
		l.PushGoFunction(func(l *lua.State) int {
			size := lua.CheckNumber(l, 2)
			label.FontSize = size
			l.PushValue(1)
			return 1
		})
		return 1

	case "color":
		l.PushGoFunction(func(l *lua.State) int {
			pushColor(l, label.Color)
			return 1
		})
		return 1

	case "setcolor":
		l.PushGoFunction(func(l *lua.State) int {
			c := checkColor(l, 2)
			label.Color = c
			l.PushValue(1)
			return 1
		})
		return 1
	}

	return 0
}
