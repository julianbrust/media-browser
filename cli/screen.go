package cli

import (
	"github.com/gdamore/tcell/v2"
	"log"
	"strconv"
)

// SetupScreen prepares a new tcell CLI screen.
// It creates and returns tcell.Screen object.
// It sets a default style config and returns it as tcell.Style.
func SetupScreen() (tcell.Screen, tcell.Style) {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err = s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	s.Clear()

	return s, defStyle
}

// CLI represents a tcell object containing the necessary Screen and Style objects.
type CLI struct {
	Screen tcell.Screen // tcell.Screen object
	Style  tcell.Style  // style configuration for a tcell screen
}

// Page represents a Page of a list of browsable content for CLI views.
type Page struct {
	Current int       // current page
	Total   int       // total number of pages
	Results int       // number of results on the current page
	Content []Content // Content for the current page
}

// Content represents an Item to be displayed on a browsable page.
type Content struct {
	Display string // value to display
	ID      int    // ID of the content element
}

// ScreenDimensions represents values to define the dimensions of a CLI screen.
type ScreenDimensions struct {
	X1, X2, Y1, Y2 int // start and end values on the X and Y Axis
}

// GetDimensions sets the screen dimensions starting from the top left.
// It uses width and height to define the bottom left params of the screen.
func GetDimensions(width int, height int) ScreenDimensions {
	return ScreenDimensions{
		X1: 0,
		X2: width - 1,
		Y1: 0,
		Y2: height - 1,
	}
}

// BuildScreen creates the content for a tcell screen.
// It creates an array of lines to display.
// header lines will be added to the top.
// page, index and content can be used to add a browsable list.
func BuildScreen(page Page, index int, header []string, content []Content, list bool) []string {
	newText := make([]string, len(header)+page.Results+1)

	for i, line := range header {
		newText[i] = line
	}

	if index > len(content)-1 {
		index = len(content) - 1
	}

	listIndex := len(header)

	if list {
		browsableList := BuildBrowsableList(page, index, content)

		for _, entry := range browsableList {
			newText[listIndex] = entry
			listIndex++
		}
	}

	return newText
}

// BuildBrowsableList creates a list of lines for a browsable page.
// It uses the current page and content on that page.
// It will prepend the item at the current index with an arrow.
func BuildBrowsableList(page Page, index int, content []Content) []string {
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

	entries[page.Results] = "page " + strconv.FormatInt(int64(page.Current), 10) +
		"/" + strconv.FormatInt(int64(page.Total), 10)

	return entries
}

// BuildErrorScreen creates a list of lines consisting of the header and the error message.
func BuildErrorScreen(header []string, err string) []string {
	var newText []string

	for _, line := range header {
		newText = append(newText, line)
	}

	newText = append(newText, "  "+err)

	return newText
}

// DrawScreen draws a tcell screen line by line.
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
