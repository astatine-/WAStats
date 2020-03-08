package main

import (
	"bufio"
	"image/color"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/psykhi/wordclouds"
)

type maskConf struct {
	File  string     `json:"file"`
	Color color.RGBA `json:"color"`
}

type wcConf struct {
	FontMaxSize     int          `json:"font_max_size"`
	FontMinSize     int          `json:"font_min_size"`
	RandomPlacement bool         `json:"random_placement"`
	FontFile        string       `json:"font_file"`
	Colors          []color.RGBA `json:"colors"`
	Width           int          `json:"width"`
	Height          int          `json:"height"`
	Mask            maskConf     `json:"mask"`
}

var wcColors = []color.RGBA{
	{0x1b, 0x1b, 0x1b, 0xff},
	{0x48, 0x48, 0x4B, 0xff},
	{0x59, 0x3a, 0xee, 0xff},
	{0x65, 0xCD, 0xFA, 0xff},
	{0x70, 0xD6, 0xBF, 0xff},
}

var defaultConf = wcConf{
	FontMaxSize:     700,
	FontMinSize:     10,
	RandomPlacement: false,
	FontFile:        "Roboto-Regular.ttf",
	Colors:          wcColors,
	Width:           4096,
	Height:          4096,
	Mask: maskConf{"", color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 0,
	}},
}

func generateWordcloud(text string) {
	stoplist := make(map[string]int, 0)
	s, err := os.Open("stoplist.txt")
	if err != nil {
		return
	}
	defer s.Close()

	slscanner := bufio.NewScanner(s)
	for slscanner.Scan() {
		line := slscanner.Text()
		stoplist[line] = 1
	}

	textlist := strings.Split(strings.ToLower(text), " ")

	filteredList := make(map[string]int)
	for _, w := range textlist {
		if stoplist[w] == 1 || len(w) <= 2 { // word in text is present in stoplist
			continue
		} else {
			filteredList[w]++
		}
	}

	for w1 := range filteredList {
		if filteredList[w1] <= 1 {
			delete(filteredList, w1)
		}
	}
	println("Generating word cloud with with:  " + strconv.Itoa(len(filteredList)) + " words (of 4 or more repetitions)")

	colors := make([]color.Color, 0)
	for _, c := range defaultConf.Colors {
		colors = append(colors, c)
	}
	var boxes []*wordclouds.Box

	filteredList["(c) 2020, Arun Tanksali"] = 6
	wc := wordclouds.NewWordcloud(filteredList,
		wordclouds.FontFile(defaultConf.FontFile),
		wordclouds.FontMaxSize(defaultConf.FontMaxSize),
		wordclouds.FontMinSize(defaultConf.FontMinSize),
		wordclouds.Colors(colors),
		wordclouds.MaskBoxes(boxes),
		wordclouds.Height(defaultConf.Height),
		wordclouds.Width(defaultConf.Width),
		wordclouds.RandomPlacement(defaultConf.RandomPlacement))

	img := wc.Draw()

	outputF, err := os.Create("wc.png")
	png.Encode(outputF, img)
	outputF.Close()
	genPDFAddWordCloud()
}
