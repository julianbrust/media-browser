package cli

import (
	"github.com/gdamore/tcell/v2"
	"log"
	"strconv"
)

func SetupScreen() (tcell.Screen, tcell.Style) {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	s.Clear()

	return s, defStyle
}

type CLI struct {
	Screen tcell.Screen
	Style  tcell.Style
}

type Page struct {
	Current int
	Total   int
	Results int
	Content []Content
}

type Content struct {
	Display string
	ID      int
}

type ScreenDimensions struct {
	X1, X2, Y1, Y2 int
}

func GetDimensions(width int, height int) ScreenDimensions {
	return ScreenDimensions{
		X1: 0,
		X2: width - 1,
		Y1: 0,
		Y2: height - 1,
	}
}

func BuildScreen(page Page, index int, header []string, content []Content, list bool) []string {
	newText := make([]string, len(header)+page.Results+1)

	for i, line := range header {
		newText[i] = line
	}

	listIndex := len(header)
	browsableList := BuildBrowsableList(page, index, content, list)

	for _, entry := range browsableList {
		newText[listIndex] = entry
		listIndex++
	}

	return newText
}

func BuildBrowsableList(page Page, index int, content []Content, list bool) []string {
	entriesAmount := len(page.Content)

	entries := make([]string, page.Results+1)

	tabString := "  "
	entriesIndex := 0

	for _, entry := range content {
		if entriesIndex < entriesAmount && entriesIndex < page.Results {
			if entriesIndex == index {
				entries[entriesIndex] = tabString + "> " + entry.Display
			} else {
				entries[entriesIndex] = tabString + "  " + entry.Display
			}
			entriesIndex++
		}
	}

	if list {
		entries[page.Results] = "page " + strconv.FormatInt(int64(page.Current), 10) +
			"/" + strconv.FormatInt(int64(page.Total), 10)
	}

	return entries
}

func DrawScreen(screen tcell.Screen, style tcell.Style, dim ScreenDimensions, text []string) {
	row := dim.Y1
	col := dim.X1
	for _, line := range text {
		for _, r := range []rune(line) {
			screen.SetContent(col, row, r, nil, style)
			col++
			if col >= dim.X2 {
				row++
				col = dim.X1
			}
			if row > dim.Y2 {
				break
			}
		}
		col = dim.X1
		row++
	}
}
