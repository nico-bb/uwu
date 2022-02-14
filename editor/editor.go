package editor

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/nico-ec/uwu/ui"
)

type Editor struct {
	ctx *ui.Context

	// Editor's resources
	font   Font
	header Image
	layout Image

	window   ui.Handle
	treeView *ui.List

	statusbar statusBar
}

func (ed *Editor) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return fmt.Errorf("closing editor")
	}
	var runes []rune
	runes = ebiten.AppendInputChars(runes[:0])
	for _, r := range runes {
		ed.ctx.AppendCharPressed(r)
		// key = rl.GetCharPressed()
	}
	mx, my := ebiten.CursorPosition()
	mleft := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	ed.ctx.UpdateUI(ui.Input{
		MPos:  ui.Point{float64(mx), float64(my)},
		MLeft: mleft,
		Enter: ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeyKPEnter),
		Del:   ebiten.IsKeyPressed(ebiten.KeyBackspace),
		Ctrl:  ebiten.IsKeyPressed(ebiten.KeyControlLeft) || ebiten.IsKeyPressed(ebiten.KeyControlRight),
		Tab:   ebiten.IsKeyPressed(ebiten.KeyTab),
		Left:  ebiten.IsKeyPressed(ebiten.KeyLeft),
		Right: ebiten.IsKeyPressed(ebiten.KeyRight),
		Up:    ebiten.IsKeyPressed(ebiten.KeyUp),
		Down:  ebiten.IsKeyPressed(ebiten.KeyDown),
	})

	ed.statusbar.updateStatus()
	return nil
}

func (ed *Editor) Draw(screen *ebiten.Image) {
	uiBuf := ed.ctx.DrawUI()
	for _, e := range uiBuf {
		switch e.Kind {
		case ui.RenderText:
			font := e.Font.(*Font)
			// FIXME: This is a bit hacky.
			// It isn't wrong, but it centers the text at the bottom of the "line"
			// Should provide more options such as center to middle
			// Instead of offsetting with the font size, offset by the text height
			// if user wants to center to middle
			text.Draw(
				screen,
				e.Text,
				font.faces[int(e.Rect.Height)],
				int(e.Rect.X),
				int(e.Rect.Y+e.Rect.Height),
				e.Clr,
			)
		case ui.RenderRectangle:
			ebitenutil.DrawRect(
				screen,
				e.Rect.X, e.Rect.Y,
				e.Rect.Width, e.Rect.Height,
				e.Clr,
			)

		case ui.RenderImageSlice:
			dstRects := [9]ui.Rectangle{}
			srcRects := [9]ui.Rectangle{}

			l := e.Constr.Left
			r := e.Constr.Right
			u := e.Constr.Up
			d := e.Constr.Down

			imgW := e.Img.GetWidth()
			imgH := e.Img.GetHeight()

			srcX0 := float64(0)
			srcX1 := l
			srcX2 := imgW - r

			srcY0 := float64(0)
			srcY1 := u
			srcY2 := imgH - d

			dstL := l
			dstR := r
			dstU := u
			dstD := d

			// if scale > 0 {
			// 	dstL *= scale
			// 	dstR *= scale
			// 	dstU *= scale
			// 	dstD *= scale
			// }

			dstX0 := e.Rect.X
			dstX1 := e.Rect.X + dstL
			dstX2 := e.Rect.X + e.Rect.Width - dstR

			dstY0 := e.Rect.Y
			dstY1 := e.Rect.Y + dstU
			dstY2 := e.Rect.Y + e.Rect.Height - dstD

			// TOP
			dstRects[0] = ui.Rectangle{X: dstX0, Y: dstY0, Width: dstL, Height: dstU}
			srcRects[0] = ui.Rectangle{X: srcX0, Y: srcY0, Width: l, Height: u}
			//
			dstRects[1] = ui.Rectangle{X: dstX1, Y: dstY0, Width: e.Rect.Width - (dstL + dstR), Height: dstU}
			srcRects[1] = ui.Rectangle{X: srcX1, Y: srcY0, Width: imgW - (l + r), Height: u}
			//
			dstRects[2] = ui.Rectangle{X: dstX2, Y: dstY0, Width: dstR, Height: dstU}
			srcRects[2] = ui.Rectangle{X: srcX2, Y: srcY0, Width: r, Height: u}
			//
			// MIDDLE
			dstRects[3] = ui.Rectangle{X: dstX0, Y: dstY1, Width: dstL, Height: e.Rect.Height - (dstU + dstD)}
			srcRects[3] = ui.Rectangle{X: srcX0, Y: srcY1, Width: l, Height: imgH - (u + d)}
			//
			dstRects[4] = ui.Rectangle{X: dstX1, Y: dstY1, Width: e.Rect.Width - (dstL + dstR), Height: e.Rect.Height - (dstU + dstD)}
			srcRects[4] = ui.Rectangle{X: srcX1, Y: srcY1, Width: imgW - (l + r), Height: imgH - (u + d)}
			//
			dstRects[5] = ui.Rectangle{X: dstX2, Y: dstY1, Width: dstR, Height: e.Rect.Height - (dstU + dstD)}
			srcRects[5] = ui.Rectangle{X: srcX2, Y: srcY1, Width: r, Height: imgH - (u + d)}
			//
			// BOTTOM
			dstRects[6] = ui.Rectangle{X: dstX0, Y: dstY2, Width: dstL, Height: dstD}
			srcRects[6] = ui.Rectangle{X: srcX0, Y: srcY2, Width: l, Height: d}
			//
			dstRects[7] = ui.Rectangle{X: dstX1, Y: dstY2, Width: e.Rect.Width - (dstL + dstR), Height: dstD}
			srcRects[7] = ui.Rectangle{X: srcX1, Y: srcY2, Width: imgW - (l + r), Height: d}
			//
			dstRects[8] = ui.Rectangle{X: dstX2, Y: dstY2, Width: dstR, Height: dstD}
			srcRects[8] = ui.Rectangle{X: srcX2, Y: srcY2, Width: r, Height: d}

			img := e.Img.(*Image)
			for i := 0; i < 9; i += 1 {
				opt := &ebiten.DrawImageOptions{}
				opt.GeoM.Scale(
					dstRects[i].Width/srcRects[i].Width,
					dstRects[i].Height/srcRects[i].Height,
				)
				opt.GeoM.Translate(dstRects[i].X, dstRects[i].Y)
				if e.Clr[3] != 0 {
					r, g, b, a := e.Clr.RGBA()
					opt.ColorM.Scale(
						float64(r)/float64(a),
						float64(g)/float64(a),
						float64(b)/float64(a),
						float64(a)/0xffff,
					)
				}
				r := image.Rect(
					int(srcRects[i].X), int(srcRects[i].Y),
					int(srcRects[i].X+srcRects[i].Width), int(srcRects[i].Y+srcRects[i].Height),
				)
				screen.DrawImage(img.data.SubImage(r).(*ebiten.Image), opt)
			}
		}
	}
}

func (e *Editor) Layout(w, h int) (int, int) {
	return 1600, 900
}

func NewEditor() *Editor {
	ed := new(Editor)
	ed.ctx = ui.NewContext()
	ed.ctx.SetCursorShapeCallback(changeEditorCursorShape)
	ui.MakeContextCurrent(ed.ctx)
	ed.font = NewFont("assets/CozetteVector.ttf", 72, []int{12})

	// TODO: refactor this into a function
	i, _, err := ebitenutil.NewImageFromFile("assets/uiHeader.png")
	if err != nil {
		panic(err)
	}
	ed.header = Image{
		data: i,
	}

	i, _, err = ebitenutil.NewImageFromFile("assets/uiLayout.png")
	if err != nil {
		panic(err)
	}
	ed.layout = Image{
		data: i,
	}

	ed.window = ui.AddWindow(
		ui.Window{
			Active: true,
			Rect:   ui.Rectangle{0, 0, 1600, 900},
			Style: ui.Style{
				Ordering: ui.StyleOrderRow,
				Padding:  0,
				Margin:   ui.Point{0, 0},
			},
			Background: ui.Background{
				Visible: true,
				Kind:    ui.BackgroundSolidColor,
				Clr:     uwuBackgroundClr,
			},
			HasHeader:    true,
			HeaderHeight: 25,
			HeaderBackground: ui.Background{
				Visible: true,
				Kind:    ui.BackgroundImageSlice,
				Clr:     ui.Color{232, 152, 168, 255},
				Img:     &ed.header,
				Constr:  ui.Constraint{2, 2, 2, 2},
			},
			HasHeaderTitle: true,
			HeaderTitle:    "UwU",
			HeaderFont:     &ed.font,
			HeaderFontSize: 12,
			HeaderFontClr:  uwuTextClr,
			// HasCloseBtn:    true,
			// CloseBtn: ui.Background{
			// 	Visible: true,
			// 	Kind:    ui.BackgroundImageSlice,
			// 	Clr:     ui.Color{255, 255, 255, 255},
			// 	Img:     &uiBtn,
			// 	Constr:  ui.Constraint{2, 2, 2, 2},
			// },
		},
	)

	rem := ui.ContainerRemainingLength(ed.window)
	lyt := ui.AddWidget(ed.window, &ui.Layout{
		Background: ui.Background{
			Visible: false,
		},
		Style: ui.Style{
			Ordering: ui.StyleOrderColumn,
			Padding:  0,
			Margin:   ui.Point{0, 0},
		},
	}, rem-20)
	ed.treeView = &ui.List{
		Background: ui.Background{
			Visible: true,
			Kind:    ui.BackgroundImageSlice,
			Clr:     ui.Color{232, 152, 168, 255},
			Img:     &ed.layout,
			Constr:  ui.Constraint{2, 2, 2, 2},
		},
		Style: ui.Style{
			Padding: 3,
			Margin:  ui.Point{5, 0},
		},

		Name:       "Root",
		Font:       &ed.font,
		TextSize:   12,
		TextClr:    uwuTextClr,
		IndentSize: 10,
	}
	ui.AddWidget(lyt, ed.treeView, 140)
	// subFolder := ui.NewSubList("subFolder")
	{
		p := openProject("./")
		fmt.Printf("%#v\n", p)
		// files, err := os.ReadDir("./")
		// if err != nil {
		// 	panic(err)
		// }
		// for _, f := range files {
		// 	fmt.Println(f.Name())
		// 	switch f.IsDir() {
		// 	case true:
		// 		ed.treeView.AddItem(&ui.SubList{Name: f.Name()})
		// 	case false:
		// 		ed.treeView.AddItem(&ui.ListItem{Name: f.Name()})
		// 	}
		// }
	}
	// ed.treeView.AddItem(&subFolder)
	// subFolder.AddItem(&ui.ListItem{Name: "file1"})
	// subFolder.AddItem(&ui.ListItem{Name: "file2"})
	// subFolder.AddItem(&ui.ListItem{Name: "file3"})
	// ed.treeView.AddItem(&ui.ListItem{Name: "file1"})
	// ed.treeView.AddItem(&ui.ListItem{Name: "file2"})
	// ed.treeView.AddItem(&ui.ListItem{Name: "file3"})
	// ed.treeView.AddItem(&ui.ListItem{Name: "file4"})
	// ed.treeView.AddItem(&ui.ListItem{Name: "file5"})

	// Text editor
	editor := &ui.TextBox{
		Background: ui.Background{
			Visible: false,
		},
		Cap:                500,
		Margin:             10,
		Font:               &ed.font,
		TextSize:           12,
		TabSize:            2,
		AutoIndent:         true,
		HasRuler:           true,
		HasSyntaxHighlight: true,
		ShowCurrentLine:    true,
	}
	editor.SetLexKeywords([]string{
		"func",
	})
	editor.SetSyntaxColors(ui.ColorStyle{
		Normal:  uwuTextClr,
		Keyword: uwuKeywordClr,
		Digit:   uwuDigitClr,
	})
	txtEdit := ui.AddWidget(lyt, editor, ui.FitContainer)

	// Status bar
	ed.statusbar = newStatusBar(ed.window, txtEdit, &ed.font)

	return ed
}

func changeEditorCursorShape(s ui.CursorShape) {
	var ebitenCursorShape ebiten.CursorShapeType
	switch s {
	case ui.CursorShapeDefault:
		ebitenCursorShape = ebiten.CursorShapeDefault
	case ui.CursorShapeText:
		ebitenCursorShape = ebiten.CursorShapeText
	}
	ebiten.SetCursorShape(ebitenCursorShape)
}
