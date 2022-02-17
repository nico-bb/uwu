package ui

const (
	tabWidth            = 80
	tabViewerInitialCap = 10
)

type (
	TabViewer struct {
		widgetRoot

		HeaderBackground Background
		HeaderHeight     float64
		headerRect       Rectangle

		TabFont     Font
		TabTextSize float64
		TabClr      Color
		tabRect     Rectangle
		tabs        []tab
		tabGens     []uint
		tabCount    int
		currentTab  tab
	}

	tab struct {
		name   string
		rect   Rectangle
		widget Widget
	}
)

func (t *TabViewer) init() {
	t.headerRect = Rectangle{
		X: t.rect.X, Y: t.rect.Y,
		Width: t.rect.Width, Height: t.HeaderHeight,
	}

	t.tabRect = Rectangle{
		X: t.rect.X, Y: t.rect.Y + t.HeaderHeight,
		Width: t.rect.Width, Height: t.rect.Height - t.HeaderHeight,
	}
	t.tabs = make([]tab, tabViewerInitialCap)
	t.tabGens = make([]uint, tabViewerInitialCap)
}

func (t *TabViewer) update() {
	// check if a new tab has been selected

	t.currentTab.widget.update()
}

func (t *TabViewer) draw(buf *renderBuffer) {
	bgEntry := t.HeaderBackground.entry(t.headerRect)
	buf.addEntry(bgEntry)

	for i := 0; i < t.tabCount; i += 1 {
		tab := &t.tabs[i]
		buf.addEntry(RenderEntry{
			Kind: RenderRectangle,
			Rect: tab.rect,
			Clr:  t.TabClr,
		})
		textSize := t.TabFont.MeasureText(tab.name, t.TabTextSize)
		buf.addEntry(RenderEntry{
			Kind: RenderText,
			Rect: Rectangle{
				X:      tab.rect.X + (tab.rect.Width/2 - textSize[0]/2),
				Y:      tab.rect.Y + (tab.rect.Height/2 - textSize[1]/2),
				Height: t.TabTextSize,
			},
			Clr:  Color{255, 255, 255, 255},
			Font: t.TabFont,
			Text: tab.name,
		})

	}
	t.currentTab.widget.draw(buf)
}

// Should the newly added tab be set as the active one?
func (t *TabViewer) AddTab(name string, w Widget) {
	w.setRect(t.tabRect)

	t.tabGens[t.tabCount] += 1
	t.tabs[t.tabCount] = tab{
		name: name,
		rect: Rectangle{
			X: t.rect.X, Y: t.rect.Y,
			Width: tabWidth, Height: t.HeaderHeight,
		},
		widget: w,
	}
	t.tabCount += 1
	t.currentTab = t.tabs[t.tabCount-1]
	w.init()
}
