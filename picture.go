package hobby

import (
	"math"

	"github.com/boxesandglue/mpgo/draw"
	"github.com/boxesandglue/mpgo/mp"
	"github.com/boxesandglue/mpgo/svg"
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
		// pic:label(text, point, anchor[, {color=, fontsize=}])
		l.PushGoFunction(func(l *lua.State) int {
			text := lua.CheckString(l, 2)
			pos := checkPoint(l, 3)
			anchorStr := lua.OptString(l, 4, "center")
			anchor := parseAnchor(anchorStr)
			label := pic.LabelWithStyle(text, pos, anchor)
			if l.Top() >= 5 && l.IsTable(5) {
				applyLabelOptions(l, 5, label)
			}
			l.PushValue(1) // return self for chaining
			return 1
		})
		return 1

	case "dotlabel":
		// pic:dotlabel(text, point, anchor[, color | {color=, fontsize=}])
		l.PushGoFunction(func(l *lua.State) int {
			text := lua.CheckString(l, 2)
			pos := checkPoint(l, 3)
			anchorStr := lua.OptString(l, 4, "center")
			anchor := parseAnchor(anchorStr)
			label := pic.LabelWithStyle(text, pos, anchor)
			color := mp.ColorCSS("black")
			if l.Top() >= 5 && !l.IsNil(5) {
				if l.IsTable(5) {
					applyLabelOptions(l, 5, label)
					color = label.Color
				} else {
					color = checkColor(l, 5)
					label.Color = color
				}
			}
			// Add the dot (filled circle at position)
			dot := mp.FullCircle()
			dot = mp.Scaled(mp.DefaultDotLabelDiam).ApplyToPath(dot)
			dot = mp.Shifted(pos.X, pos.Y).ApplyToPath(dot)
			dot.Style.Fill = color
			dot.Style.Stroke = mp.ColorCSS("none")
			pic.AddPath(dot)
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

	case "converttopaths":
		// pic:converttopaths(face) - convert all labels to glyph outline paths
		l.PushGoFunction(func(l *lua.State) int {
			face := checkFace(l, 2)
			if err := pic.ConvertLabelsToPathsWithFont(face); err != nil {
				lua.Errorf(l, "converttopaths: %s", err)
				return 0
			}
			l.PushValue(1) // return self for chaining
			return 1
		})
		return 1

	case "llcorner":
		minX, minY, _, _ := pictureBBox(pic)
		pushPoint(l, mp.P(minX, minY))
		return 1

	case "lrcorner":
		_, minY, maxX, _ := pictureBBox(pic)
		pushPoint(l, mp.P(maxX, minY))
		return 1

	case "ulcorner":
		minX, _, _, maxY := pictureBBox(pic)
		pushPoint(l, mp.P(minX, maxY))
		return 1

	case "urcorner":
		_, _, maxX, maxY := pictureBBox(pic)
		pushPoint(l, mp.P(maxX, maxY))
		return 1

	case "center":
		minX, minY, maxX, maxY := pictureBBox(pic)
		pushPoint(l, mp.P((minX+maxX)/2, (minY+maxY)/2))
		return 1

	case "bbox":
		l.PushGoFunction(func(l *lua.State) int {
			minX, minY, maxX, maxY := pictureBBox(pic)
			pushPath(l, bboxPath(minX, minY, maxX, maxY))
			return 1
		})
		return 1
	}

	return 0
}

// pictureBBox computes the bounding box of a picture.
// If the picture has a clip path, the clip path bounds are used (like MetaPost).
// Otherwise, bounds of all paths are aggregated with stroke width padding.
func pictureBBox(pic *draw.Picture) (minX, minY, maxX, maxY float64) {
	if clip := pic.ClipPath(); clip != nil {
		return svg.PathBBox(clip)
	}
	minX, minY = math.Inf(1), math.Inf(1)
	maxX, maxY = math.Inf(-1), math.Inf(-1)
	for _, p := range pic.Paths() {
		pMinX, pMinY, pMaxX, pMaxY := svg.PathBBox(p)
		halfStroke := p.Style.StrokeWidth / 2
		pMinX -= halfStroke
		pMinY -= halfStroke
		pMaxX += halfStroke
		pMaxY += halfStroke
		if pMinX < minX {
			minX = pMinX
		}
		if pMinY < minY {
			minY = pMinY
		}
		if pMaxX > maxX {
			maxX = pMaxX
		}
		if pMaxY > maxY {
			maxY = pMaxY
		}
	}
	if math.IsInf(minX, 1) {
		return 0, 0, 0, 0
	}
	return
}

// applyLabelOptions reads label options from a Lua table at the given stack index.
// Supported keys: color (Color), fontsize (number).
func applyLabelOptions(l *lua.State, index int, label *mp.Label) {
	l.Field(index, "color")
	if !l.IsNil(-1) {
		label.Color = checkColor(l, l.Top())
	}
	l.Pop(1)

	l.Field(index, "fontsize")
	if !l.IsNil(-1) {
		label.FontSize = lua.CheckNumber(l, l.Top())
	}
	l.Pop(1)
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
