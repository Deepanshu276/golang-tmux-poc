package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/rivo/tview"
)

// Move to the previous slide
func previousSlide() {
	currentActivePage = (currentActivePage - 1 + len(Tabs)) % len(Tabs)
	info.Highlight(strconv.Itoa(regionIds[currentActivePage])).
		ScrollToHighlight()
}

// Move to the next slide
func nextSlide() {
	currentActivePage = (currentActivePage + 1) % len(Tabs)
	info.Highlight(strconv.Itoa(regionIds[currentActivePage])).
		ScrollToHighlight()
}

// Remove the slide with the given index
func removeSlide(s int) {
	Tabs = append(Tabs[:s], Tabs[s+1:]...)
	regionIds = append(regionIds[:s], regionIds[s+1:]...)
	totalPageCount = len(Tabs)
	info.Clear()
	for index, tabSlide := range Tabs {
		oldIndex := tabSlide.index
		tabSlide.index = index
		fmt.Fprintf(info, `["%d"]%s[white][""]  `, oldIndex, fmt.Sprintf("%d %s", tabSlide.index+1, tabSlide.title))
	}
	pages.RemovePage(strconv.Itoa(s))
	previousSlide()
}

// Adds a slide to the end of currently present slides
func addSlide() {
	tabSlide := *newTab("bash", os.Getenv("SHELL"))
	Tabs = append(Tabs, tabSlide)
	pages.AddPage(strconv.Itoa(tabSlide.index), tabSlide.primitive, true, tabSlide.index == 0)
	fmt.Fprintf(info, `["%d"]%s[white][""]  `, tabSlide.index, fmt.Sprintf("%d %s", tabSlide.index+1, tabSlide.title))
	currentActivePage = tabSlide.index
	info.Highlight(strconv.Itoa(currentActivePage)).
		ScrollToHighlight()
}

func initTerminalMux() *tview.Flex {
	// Initial Slides
	Tabs = append(Tabs, *newTab("kite", os.Getenv("SHELL")))
	Tabs = append(Tabs, *newTab("bash", os.Getenv("SHELL")))

	// Set the bottom navigation bar
	info.
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			pages.SwitchToPage(added[0])
		})

	for _, slide := range Tabs {
		pages.AddPage(strconv.Itoa(slide.index), slide.primitive, true, slide.index == 0)
		fmt.Fprintf(info, `["%d"]%s[white][""]  `, slide.index, fmt.Sprintf("%d %s", slide.index+1, slide.title))
	}
	info.Highlight("0")

	// Returns the main view & layout for the app
	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(pages, 0, 1, true).
		AddItem(info, 1, 1, false)
}
