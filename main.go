// Copyright 2026 The Sieve Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"strings"
)

//go:embed books/*
var Books embed.FS

//go:embed archive.zip
var Archive embed.FS

// Book is a book
type Book struct {
	Name string
	Text []byte
	Real bool
}

// LoadBooks loads books
func LoadBooks() []Book {
	books := []Book{
		{
			Real: true,
			Name: "10.txt.utf-8.bz2",
		},
		{
			Real: true,
			Name: "11.txt.utf-8.bz2",
		},
		{
			Real: true,
			Name: "43.txt.utf-8.bz2",
		},
		{
			Real: true,
			Name: "pg74.txt.bz2",
		},
		{
			Real: true,
			Name: "76.txt.utf-8.bz2",
		},
		{
			Real: true,
			Name: "84.txt.utf-8.bz2",
		},
		{
			Real: true,
			Name: "100.txt.utf-8.bz2",
		},
		{
			Real: true,
			Name: "145.txt.utf-8.bz2",
		},
		{
			Real: true,
			Name: "768.txt.utf-8.bz2",
		},
		{
			Real: true,
			Name: "1260.txt.utf-8.bz2",
		},
		{
			Real: true,
			Name: "1342.txt.utf-8.bz2",
		},
		{
			Real: true,
			Name: "1837.txt.utf-8.bz2",
		},
		{
			Real: true,
			Name: "2641.txt.utf-8.bz2",
		},
		{
			Real: true,
			Name: "2701.txt.utf-8.bz2",
		},
		{
			Real: true,
			Name: "3176.txt.utf-8.bz2",
		},
		{
			Real: true,
			Name: "37106.txt.utf-8.bz2",
		},
		{
			Real: true,
			Name: "64317.txt.utf-8.bz2",
		},
		{
			Real: true,
			Name: "67979.txt.utf-8.bz2",
		},
		{
			Real: false,
			Name: "gemma4.txt.bz2",
		},
		{
			Real: false,
			Name: "gpt-oss.txt.bz2",
		},
		{
			Real: false,
			Name: "llama3.1.txt.bz2",
		},
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

// Prompt is a llm prompt
type Prompt struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// Query submits a query to the llm
func Query(query string) string {
	prompt := Prompt{
		Model:  *FlagModel,
		Prompt: query,
	}
	data, err := json.Marshal(prompt)
	if err != nil {
		panic(err)
	}
	buffer := bytes.NewBuffer(data)
	response, err := http.Post("http://localhost:11434/api/generate", "application/json", buffer)
	if err != nil {
		panic(err)
	}
	reader, answer := bufio.NewReader(response.Body), ""
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		data := map[string]interface{}{}
		err = json.Unmarshal([]byte(line), &data)
		text := data["response"].(string)
		answer += text
	}
	return answer
}

// Symbols are some symbols
type Symbols [2]byte

// Iterate iterate the symbols
func (s *Symbols) Iterate(n byte) {
	for i := range s[:len(s)-1] {
		s[i] = s[i+1]
	}
	s[len(s)-1] = n
}

// Target is a model to target
type Target struct {
	Count map[Symbols]uint64
	Total uint64
}

// Targets is a set of targets
type Targets []Target

// Score is the score function
func (t Targets) Score(a int, data []byte) float64 {
	var symbols Symbols
	sum := 0.0
	for i := range t {
		sum += float64(t[i].Total)
	}
	p := math.Log(float64(t[a].Total+1) / (sum + float64(len(t))))
	set := make(map[Symbols]bool)
	for _, symbol := range data {
		if symbol == '\r' || symbol == '\n' {
			continue
		}
		symbols.Iterate(symbol)
		set[symbols] = true
	}
	for symbols := range set {
		p += math.Log(float64(t[a].Count[symbols]+1) / (float64(t[a].Total) + float64(len(t[a].Count))))

	}
	return p
}

var (
	// FlagNN nearest neighbor mode
	FlagNN = flag.Bool("nn", false, "nearest neighbor mode")
	// FlagQuery submit a query to the llm
	FlagQuery = flag.String("query", "", "query the llm")
	// FlagModel the model to use
	FlagModel = flag.String("model", "gemma4", "the model to use")
	// FlagGenerate generates content
	FlagGenerate = flag.Bool("generate", false, "generate content")
)

// NNMode is the nearest neighbor mode
func NNMode() {
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
	cs := func(a, b []float64) float64 {
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
	fmt.Println(cs(histograms[0][:], histograms[1][:]))
	fmt.Println(cs(histograms[0][:], histograms[2][:]))
	fmt.Println(cs(histograms[1][:], histograms[2][:]))
}

func main() {
	flag.Parse()

	if *FlagNN {
		NNMode()
		return
	}

	if *FlagQuery != "" {
		fmt.Println(Query(*FlagQuery))
		return
	}

	if *FlagGenerate {
		rng := rand.New(rand.NewSource(1))
		results := Query("What is the meaning of life? Be verbose in your answer.")
		fmt.Println(results)
		for {
			words := strings.Fields(results)
			next := words[rng.Intn(len(words))]
			next = strings.ToLower(strings.Trim(next, ".!?,"))
			results = Query(fmt.Sprintf("What is the meaning of %s? Be verbose in your answer.", next))
			fmt.Println(results)
		}
	}

	books := LoadBooks()
	data := [][]byte{
		books[4].Text[9*1024 : 10*1024],
		books[5].Text[8*1024 : 9*1024],
		[]byte(fake[:1024]),
	}
	var classes [][]byte
	{
		count := 0
		for _, b := range books {
			if !b.Real {
				count++
			}
		}
		index := 0
		classes = make([][]byte, count)
		for i := range books {
			if !books[i].Real {
				classes[index] = books[i].Text
				index++
			}
		}
	}
	targets := make(Targets, len(classes))
	for i := range targets {
		targets[i].Count = make(map[Symbols]uint64)
	}
	for i, d := range classes {
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

	count := 0
	for i := range books {
		if books[i].Real {
			count++
		}
	}
	reals := make(Targets, count)
	for i := range reals {
		reals[i].Count = make(map[Symbols]uint64)
	}
	for i, d := range books {
		if !d.Real {
			continue
		}
		var symbols Symbols
		for _, symbol := range d.Text {
			if symbol == '\r' || symbol == '\n' {
				continue
			}
			symbols.Iterate(symbol)
			reals[i].Count[symbols]++
			reals[i].Total++
		}
	}

	test := func(i int) {
		data := data[i]
		scores := make([]float64, len(targets))
		for i := range scores {
			scores[i] = targets.Score(i, data)
		}
		histogram := make([]int, len(scores))
		for r := range reals {
			score := reals.Score(r, data)
			for i := range scores {
				if scores[i] < score {
					histogram[i]++
				}
			}
		}
		score := 0
		for i := range histogram {
			if histogram[i] < count/2 {
				score++
			}
		}
		if score > 0 {
			fmt.Println(scores, histogram, "fake")
		} else {
			fmt.Println(scores, histogram, "real")
		}
	}
	for i := range data {
		test(i)
	}
}
