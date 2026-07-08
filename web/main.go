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
	"sort"
	"strings"
	"syscall/js"
)

//go:embed  books/*
var Books embed.FS

const (
	SlopAvg    = 0.002318874154124976
	SlopStddev = 5.8209022598764776e-05
	NotAvg     = 0.002729698680992552
	NotStddev  = 3.5226320411986187e-07
)

// Text is a text type
type Text uint8

const (
	TextSlop Text = iota
	TextNot
)

// Rank is a text rank
type Rank struct {
	Rank float64
	Type Text
}

var Ranks = []Rank{
	{0.0027188210448006, TextNot},
	{0.0024715820022209, TextNot},
	{0.0026193556997382, TextNot},
	{0.0027671081299420, TextNot},
	{0.0026994349599398, TextNot},
	{0.0025840830287612, TextNot},
	{0.0032537747798464, TextNot},
	{0.0030260218133745, TextNot},
	{0.0030208543655926, TextNot},
	{0.0030199160180444, TextNot},
	{0.0029726773896805, TextNot},
	{0.0026546765932896, TextNot},
	{0.0042592065200255, TextNot},
	{0.0024739898736685, TextNot},
	{0.0027594223439506, TextNot},
	{0.0026971659298036, TextNot},
	{0.0025918806265325, TextNot},
	{0.0026751595676999, TextNot},
	{0.0025289134565670, TextNot},
	{0.0025590231275894, TextNot},
	{0.0025814640050012, TextNot},
	{0.0023368627091383, TextNot},
	{0.0024870842928849, TextNot},
	{0.0029075499031191, TextNot},
	{0.0025487047367912, TextNot},
	{0.0023944620415836, TextNot},
	{0.0028243499031931, TextNot},
	{0.0027792506686523, TextNot},
	{0.0022983617071924, TextNot},
	{0.0033354077537305, TextNot},
	{0.0026267206446545, TextNot},
	{0.0022798844561389, TextNot},
	{0.0021906492701910, TextNot},
	{0.0024939065349125, TextNot},
	{0.0022579607558760, TextNot},
	{0.0029258282744318, TextNot},
	{0.0027802379613480, TextNot},
	{0.0030600140400914, TextNot},
	{0.0025397396367300, TextNot},
	{0.0025992219533246, TextNot},
	{0.0028060770270868, TextNot},
	{0.0022546838432058, TextNot},
	{0.0024367762075795, TextNot},
	{0.0026651291875348, TextNot},
	{0.0028795887644575, TextNot},
	{0.0030133758611835, TextNot},
	{0.0023960381190513, TextNot},
	{0.0030995077704995, TextNot},
	{0.0025965666477240, TextNot},
	{0.0026769202985166, TextNot},
	{0.0029379048771785, TextNot},
	{0.0030094657475025, TextNot},
	{0.0025344136713206, TextNot},
	{0.0028319165180598, TextNot},
	{0.0025604733845274, TextNot},
	{0.0028865071627312, TextNot},
	{0.0026622015240734, TextNot},
	{0.0035105711405213, TextNot},
	{0.0027806458982686, TextNot},
	{0.0027587114059970, TextNot},
	{0.0026219374546689, TextNot},
	{0.0027395088538209, TextNot},
	{0.0027085489113351, TextNot},
	{0.0027325167866255, TextNot},
	{0.0024445933956482, TextSlop},
	{0.0015424269288494, TextSlop},
	{0.0025149471576455, TextSlop},
	{0.0024220115114023, TextSlop},
	{0.0023339019148041, TextSlop},
	{0.0017118832908543, TextSlop},
	{0.0026346499019590, TextSlop},
	{0.0023131888268327, TextSlop},
	{0.0019824144576281, TextSlop},
	{0.0020004694725476, TextSlop},
	{0.0024191015400646, TextSlop},
	{0.0024318700824397, TextSlop},
	{0.0025487131930542, TextSlop},
	{0.0024489039635493, TextSlop},
	{0.0026200062738841, TextSlop},
	{0.0021308244593731, TextSlop},
	{0.0025408850862659, TextSlop},
	{0.0023466259895508, TextSlop},
	{0.0027722206828839, TextSlop},
	{0.0026262351013152, TextSlop},
	{0.0020511025331169, TextSlop},
	{0.0025411085842163, TextSlop},
	{0.0025550910644093, TextSlop},
	{0.0023265712402190, TextSlop},
	{0.0022505472938499, TextSlop},
	{0.0015494579987679, TextSlop},
	{0.0024414392046835, TextSlop},
	{0.0023367264173967, TextSlop},
	{0.0024874132527845, TextSlop},
	{0.0025623734728518, TextSlop},
	{0.0024054167664316, TextSlop},
	{0.0022563102793485, TextSlop},
	{0.0025699314925276, TextSlop},
	{0.0022083485785579, TextSlop},
	{0.0024047174992061, TextSlop},
	{0.0021898323460208, TextSlop},
	{0.0024831399177127, TextSlop},
	{0.0022166256740873, TextSlop},
	{0.0024844674078830, TextSlop},
	{0.0021142262108663, TextSlop},
	{0.0019836141218865, TextSlop},
	{0.0024520837641847, TextSlop},
	{0.0023170101015378, TextSlop},
	{0.0024531369435419, TextSlop},
	{0.0024855228483818, TextSlop},
	{0.0021475658591521, TextSlop},
	{0.0026414136916751, TextSlop},
	{0.0021671147561427, TextSlop},
	{0.0022909557569519, TextSlop},
	{0.0020329115799536, TextSlop},
	{0.0029593567459492, TextSlop},
	{0.0026884490230808, TextSlop},
	{0.0024110046547876, TextSlop},
	{0.0024822767049441, TextSlop},
	{0.0022355475692285, TextSlop},
	{0.0025201338997538, TextSlop},
	{0.0022753110504954, TextSlop},
	{0.0021692623494004, TextSlop},
	{0.0018238113597638, TextSlop},
	{0.0021154275181101, TextSlop},
	{0.0017224302827641, TextSlop},
	{0.0024141895208603, TextSlop},
	{0.0018532019733349, TextSlop},
	{0.0025454933226278, TextSlop},
}

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
		{
			Real:  false,
			Name:  "gemma4.txt.bz2",
			Index: 18,
		},
		{
			Real:  false,
			Name:  "gpt-oss.txt.bz2",
			Index: 19,
		},
		{
			Real:  false,
			Name:  "llama3.1.txt.bz2",
			Index: 20,
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
func TestMode(sample string) (float64, float64, float64) {
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
		return result,
			(1 + math.Erf((result-SlopAvg)/(SlopStddev*math.Sqrt(2)))) / 2,
			(1 + math.Erf((result-NotAvg)/(NotStddev*math.Sqrt(2)))) / 2
	}
}

// Class is a model for a class
type Class struct {
	Graph
	Total float64
	List  []string
}

// Classes is a set of classes
type Classes []Class

// Score is the score function
func (c Classes) Score(a int, data []string) float64 {
	sum := 0.0
	for i := range c {
		sum += c[i].Total
	}
	p := math.Log(float64(c[a].Total+1) / (sum + float64(len(c))))
	length := float64(len(data))
	for _, symbol := range data {
		p += math.Log(float64(c[a].Ranks[symbol]+1) / (float64(c[a].Total) + float64(len(c[a].Ranks))))

	}
	return p / length
}

// TestMode2 test
func TestMode2(sample string) bool {
	books := LoadBooks()
	rng := rand.New(rand.NewSource(1))
	classes := make(Classes, len(books))
	for i, book := range books {
		text := string(book.Text)
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
			classes[i].Graph = g
			classes[i].Total = count
			classes[i].List = list
		}
	}

	max, index := -math.MaxFloat64, 0
	for i := range classes {
		score := classes.Score(i, classes[i].List)
		fmt.Println("score=", score)
		if score > max {
			max, index = score, i
		}
	}

	return index == 0
}

func processText(this js.Value, args []js.Value) any {
	if len(args) > 0 {
		index := args[0].Int()
		input := args[1].String()
		bytes := []byte(input)
		if len(bytes) < 1024 {
			return "number of bytes is less than 1024"
		}
		if 1024*index+1024 > len(bytes) {
			return ""
		}
		bytes = bytes[index*1024 : index*1024+1024]
		result, _, _ := TestMode(string(bytes))
		type Result struct {
			Rank
			Diff float64
		}
		results := make([]Result, 0, len(Ranks))
		for _, rank := range Ranks {
			diff := math.Abs(rank.Rank - result)
			results = append(results, Result{
				Rank: rank,
				Diff: diff,
			})
		}
		sort.Slice(results, func(i, j int) bool {
			return results[i].Diff < results[j].Diff
		})
		var histogram [2]int
		for i := range results[:64] {
			histogram[results[i].Type]++
		}
		return fmt.Sprintf("%d %f probability slop and %f probability not\n", index, float64(histogram[0])/64.0, float64(histogram[1])/64.0)
		//return fmt.Sprintf("%d %f probability slop and %f probability not\n", index, slop, not)
	}
	return ""
}

func processText2(this js.Value, args []js.Value) any {
	if len(args) > 0 {
		input := args[0].String()
		not := TestMode2(input)
		if not {
			return "not"
		}
	}
	return "slop"
}

func main() {
	js.Global().Set("goProcessText", js.FuncOf(processText2))
	select {}
}
