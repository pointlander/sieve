// Copyright 2026 The Sieve Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build js && wasm

package main

import (
	"compress/bzip2"
	"embed"
	"fmt"
	"io"
	"math"
	"math/rand"
	"strings"
	"syscall/js"
)

//go:embed 10.txt.utf-8.bz2
var Books embed.FS

const (
	Avg    = 0.002657625207559362
	Stddev = 4.0680864464795154e-05
)

// Book is a book
type Book struct {
	Name  string
	Text  []byte
	Real  bool
	Index int
}

// LoadBooks loads books
func LoadBooks() []Book {
	books := []Book{
		{
			Real:  true,
			Name:  "10.txt.utf-8.bz2",
			Index: 0,
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
		books[i].Text = load(fmt.Sprintf("%s", books[i].Name))
	}
	return books
}

// Node is a node in a graph
type Node struct {
	Links map[string]uint64
	Keys  []string
}

// Graph is a graph
type Graph struct {
	Keys  []string
	Graph map[string]Node
	Ranks map[string]uint64
}

// NewGraph makes a new graph
func NewGraph() Graph {
	return Graph{
		Graph: make(map[string]Node),
		Ranks: make(map[string]uint64),
	}
}

// LearnFast adds context to a model
func (g *Graph) LearnFast(delta float64, iterations int, rng *rand.Rand, words, list []string, size int) float64 {
	for i, word := range words[:len(words)-1] {
		{
			node := g.Graph[word]
			if node.Links == nil {
				g.Keys = append(g.Keys, word)
				node.Links = make(map[string]uint64)
				node.Keys = make([]string, 0, 8)
			}
			count, ok := node.Links[words[i+1]]
			if !ok {
				node.Keys = append(node.Keys, words[i+1])
			}
			count++
			node.Links[words[i+1]] = count
			g.Graph[word] = node
		}
	}
	word := words[0]
	node := g.Graph[word]
	previous := math.MaxFloat64
	for i := range iterations {
		g.Ranks[word]++
		if rng.Float64() > .9 {
			index := rng.Intn(len(words))
			word = words[index]
			node = g.Graph[word]
		}
		for len(node.Keys) == 0 {
			index := rng.Intn(len(words))
			word = words[index]
			node = g.Graph[word]
		}
		sum := uint64(0)
		for _, value := range node.Keys {
			sum += node.Links[value]
		}
		total, selected := uint64(0), uint64(rng.Intn(int(sum)))
		for _, value := range node.Keys {
			total += node.Links[value]
			if selected < total {
				word = value
				node = g.Graph[word]
				break
			}
		}
		if (i+1)%len(g.Graph) == 0 {
			current, count := 0.0, float64(i)
			for _, word := range list {
				current += float64(g.Ranks[word]) / count
			}
			current /= float64(size)
			if math.Abs(current-previous) < delta {
				return count
			}
			previous = current
		}
	}
	return -1
}

// TestMode test
func TestMode(sample string) float64 {
	books := LoadBooks()
	rng := rand.New(rand.NewSource(1))
	text := string(books[0].Text)
	words := strings.Fields(text)
	{
		suffix := strings.Fields(sample)
		cp := make([]string, len(words))
		copy(cp, words)
		has, list := make(map[string]bool), make([]string, 0, 8)
		for _, word := range suffix {
			if !has[word] {
				has[word] = true
				list = append(list, word)
			}
		}
		words := append(cp, suffix...)
		g := NewGraph()
		count := g.LearnFast(1e-5, 8*1024*1024, rng, words, list, len(list))
		sum := 0.0
		for _, value := range list {
			sum += float64(g.Ranks[value]) / float64(count)
		}
		result := float64(sum) / float64(len(list))
		return (1 + math.Erf((result-Avg)/(Stddev*math.Sqrt(2)))) / 2
	}
}

func processText(this js.Value, args []js.Value) any {
	if len(args) > 0 {
		input := args[0].String()
		bytes := []byte(input)
		if len(bytes) < 1024 {
			return "number of bytes is less than 1024"
		}
		bytes = bytes[:1024]
		pReal := TestMode(string(bytes))
		return fmt.Sprintf("%f probability real and %f probability AI", pReal, 1-pReal)
	}
	return ""
}

func main() {
	js.Global().Set("goProcessText", js.FuncOf(processText))
	select {}
}
