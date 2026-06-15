// Copyright 2026 The Sieve Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"compress/bzip2"
	"embed"
	"fmt"
	"io"
	"math"
)

//go:embed books/*
var Books embed.FS

// Book is a book
type Book struct {
	Name string
	Text []byte
}

// LoadBooks loads books
func LoadBooks() []Book {
	books := []Book{
		{Name: "10.txt.utf-8.bz2"},
		{Name: "11.txt.utf-8.bz2"},
		{Name: "43.txt.utf-8.bz2"},
		{Name: "pg74.txt.bz2"},
		{Name: "76.txt.utf-8.bz2"},
		{Name: "84.txt.utf-8.bz2"},
		{Name: "100.txt.utf-8.bz2"},
		{Name: "145.txt.utf-8.bz2"},
		{Name: "768.txt.utf-8.bz2"},
		{Name: "1260.txt.utf-8.bz2"},
		{Name: "1342.txt.utf-8.bz2"},
		{Name: "1837.txt.utf-8.bz2"},
		{Name: "2641.txt.utf-8.bz2"},
		{Name: "2701.txt.utf-8.bz2"},
		{Name: "3176.txt.utf-8.bz2"},
		{Name: "37106.txt.utf-8.bz2"},
		{Name: "64317.txt.utf-8.bz2"},
		{Name: "67979.txt.utf-8.bz2"},
	}
	load := func(book string) []byte {
		file, err := Books.Open(book)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		breader := bzip2.NewReader(file)
		data, err := io.ReadAll(breader)
		if err != nil {
			panic(err)
		}
		return data
	}
	for i := range books {
		books[i].Text = load(fmt.Sprintf("books/%s", books[i].Name))
	}
	return books
}

const fake = `A shimmering twilight permanently blankets the Cerulean Hinterlands, a vast, primordial basin where nature defies conventional biology. The ground beneath is a spongy, resilient carpet of bioluminescent moss that pulses in tandem with a slow, rhythmic planetary heartbeat. Towering above this neon floor are the Goliath Redwoods, colossal flora whose bark resembles liquid obsidian, reflecting the twin moons that hang suspended in a violet sky. Instead of leaves, these silent giants sprout iridescent, translucent crystalline fronds that chime like distant glass wind chimes whenever the atmospheric thermal currents sweep through the valley.

Meandering through the heart of the hinterlands is the Whisperwind River, a stream of liquid quicksilver that flows uphill, defying gravity by climbing the tiered obsidian terraces. The water glows with a soft, internal amber warmth, casting dancing shadows on the surrounding stone formations. Flocks of featherless, moth-winged avians dance above the water's surface, leaving trails of stardust in their wake.

The air is thick with the sweet, crisp scent of crushed ozone and wild, oversized vanilla orchids that bloom only in the shadows. There are no harsh winds here, only a perpetual, comforting breeze that carries the distant, melodic echoes of the valley’s deep caverns. It is a sanctuary of surreal stillness, where the boundary between organic life and mineral magic blurs entirely, creating an untamed wilderness that feels both anciently grounded and beautifully alien.`

// Symbols are some symbols
type Symbols [2]byte

type Target struct {
	Count map[Symbols]uint64
	Total uint64
}

func main() {
	books := LoadBooks()
	a, b := books[4].Text[9*1024:10*1024], books[5].Text[8*1024:9*1024]
	fake := []byte(fake[:1024])
	fmt.Println(len(a), len(b), len(fake))
	data := [][]byte{
		a,
		b,
		fake,
	}
	var histograms [3][256]float64
	for i, d := range data {
		for _, symbol := range d {
			histograms[i][symbol]++
		}
	}
	dot := func(a, b []float64) float64 {
		sum := 0.0
		for i, value := range a {
			sum += value * b[i]
		}
		return sum
	}
	coss := func(a, b []float64) float64 {
		aa := dot(a, a)
		bb := dot(b, b)
		if aa == 0 {
			return 0
		}
		if bb == 0 {
			return 0
		}
		return dot(a, b) / (math.Sqrt(aa) * math.Sqrt(bb))
	}
	fmt.Println(coss(histograms[0][:], histograms[1][:]))
	fmt.Println(coss(histograms[0][:], histograms[2][:]))
	fmt.Println(coss(histograms[1][:], histograms[2][:]))

	targets := make([]Target, len(data))
	for i := range targets {
		targets[i].Count = make(map[Symbols]uint64)
	}
	for i, d := range data {
		var symbols Symbols
		for _, symbol := range d {
			if symbol == '\r' || symbol == '\n' {
				continue
			}
			symbols[0], symbols[1] = symbols[1], symbol
			targets[i].Count[symbols]++
			targets[i].Total++
		}
	}
	prob := func(a, b int) float64 {
		var symbols Symbols
		sum := 0.0
		for i := range targets {
			sum += float64(targets[i].Total)
		}
		p := math.Log(float64(targets[a].Total+1) / (sum + float64(len(targets))))
		set := make(map[Symbols]bool)
		for _, symbol := range data[b] {
			if symbol == '\r' || symbol == '\n' {
				continue
			}
			symbols[0], symbols[1] = symbols[1], symbol
			set[symbols] = true
		}
		for symbols := range set {
			p += math.Log(float64(targets[a].Count[symbols]+1) / (float64(targets[a].Total) + float64(len(targets[a].Count))))

		}
		return p
	}
	fmt.Println()
	fmt.Println(prob(0, 1))
	fmt.Println(prob(0, 2))
	fmt.Println(prob(1, 2))
}
