package ui

import "log"

type Window struct {
	handle     WinHandle
	Active     bool
	Rect       Rectangle
	activeRect Rectangle
	Style      Style
	Background Background
	widgets    WidgetList

	HasBorders  bool
	BorderWidth float64
	BorderColor Color

	HasHeader        bool
	HeaderHeight     float64
	HeaderBackground Background
	headerRect       Rectangle
	HasHeaderTitle   bool
	HeaderTitle      string
	HeaderFont       Font
	HeaderFontSize   float64
	HeaderFontClr    Color
	headerTitlePos   Point

	HasCloseBtn  bool
	CloseBtn     Background
	closeBtnRect Rectangle
}

func (win *Window) initWindow() {
	if (win.HasCloseBtn || win.HasHeaderTitle) && !win.HasHeader {
		// What is the best behavior here? Should the UI force a header on the window?
		// Or should it disable the Close button?
		log.SetPrefix("[UI Error]: ")
		log.Fatalln("Can not add a Close button on a headerless Window")
		win.HasCloseBtn = false
	}
	if win.HasHeader {
		if win.HasCloseBtn {
			width := win.HeaderHeight - (win.Style.Margin[1] * 2)
			win.closeBtnRect = Rectangle{
				X:      win.Rect.X + win.Rect.Width - width - win.Style.Margin[0],
				Y:      win.Rect.Y + win.Style.Margin[1],
				Width:  width,
				Height: width,
			}
		}
		win.headerRect = Rectangle{
			X: win.Rect.X, Y: win.Rect.Y,
			Width: win.Rect.Width, Height: win.HeaderHeight,
		}
		win.activeRect = Rectangle{
			X: win.Rect.X, Y: win.Rect.Y + win.HeaderHeight,
			Width: win.Rect.Width, Height: win.Rect.Height - win.HeaderHeight,
		}
		if win.HasHeaderTitle {
			titleSize := win.HeaderFont.MeasureText(win.HeaderTitle, win.HeaderFontSize)
			win.headerTitlePos = Point{
				win.headerRect.X + (win.headerRect.Width/2 - titleSize[0]/2),
				win.headerRect.Y + (win.headerRect.Height/2 - titleSize[1]/2),
			}
		}
	} else {
		win.activeRect = win.Rect
	}
	win.widgets.initList(win.Style)
}

func (win *Window) update() {
	win.widgets.updateWidgets()
}

func (win *Window) draw(buf *renderBuffer) {
	bgEntry := win.Background.entry(win.Rect)
	buf.addEntry(bgEntry)

	if win.HasHeader {
		hdrEntry := win.HeaderBackground.entry(win.headerRect)
		buf.addEntry(hdrEntry)
		if win.HasHeaderTitle {
			buf.addEntry(RenderEntry{
				Kind: RenderText,
				Rect: Rectangle{
					X:      win.headerTitlePos[0],
					Y:      win.headerTitlePos[1],
					Height: win.HeaderFontSize,
				},
				Clr:  win.HeaderFontClr,
				Font: win.HeaderFont,
				Text: win.HeaderTitle,
			})
		}
		if win.HasCloseBtn {
			btnEntry := win.CloseBtn.entry(win.closeBtnRect)
			buf.addEntry(btnEntry)
		}
	}

	win.widgets.drawWidgets(buf)

	if win.HasBorders {
		buf.addEntry(RenderEntry{
			Kind: RenderRectangle,
			Rect: Rectangle{
				X: win.Rect.X, Y: win.Rect.Y,
				Width: win.BorderWidth, Height: win.Rect.Height,
			},
			Clr: win.BorderColor,
		})
		buf.addEntry(RenderEntry{
			Kind: RenderRectangle,
			Rect: Rectangle{
				X: win.Rect.X, Y: win.Rect.Y,
				Width: win.Rect.Width, Height: win.BorderWidth,
			},
			Clr: win.BorderColor,
		})
		buf.addEntry(RenderEntry{
			Kind: RenderRectangle,
			Rect: Rectangle{
				X: win.Rect.X + win.Rect.Width - win.BorderWidth, Y: win.Rect.Y,
				Width: win.BorderWidth, Height: win.Rect.Height,
			},
			Clr: win.BorderColor,
		})
		buf.addEntry(RenderEntry{
			Kind: RenderRectangle,
			Rect: Rectangle{
				X: win.Rect.X, Y: win.Rect.Y + win.Rect.Height - win.BorderWidth,
				Width: win.Rect.Width, Height: win.BorderWidth,
			},
			Clr: win.BorderColor,
		})
	}
}

func (w *Window) AddWidget(wgt Widget, length int) {
	w.widgets.addWidget(wgt, w.activeRect, length)
}

func (w *Window) RemainingLength() int {
	return w.widgets.getRemainingLen(w.activeRect)
}

type WinHandle struct {
	id  int
	gen uint
}

func (h WinHandle) AddWidget(wgt Widget, length int) {
	getWindow(h).AddWidget(wgt, length)
}

func (h WinHandle) RemainingLength() int {
	return getWindow(h).RemainingLength()
}
