// Copyright 2026 The Sieve Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"bytes"
	"compress/bzip2"
	"embed"
	"encoding/gob"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strings"

	"github.com/pointlander/gradient"
)

//go:embed books/*
var Books embed.FS

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
			Real:  true,
			Name:  "11.txt.utf-8.bz2",
			Index: 1,
		},
		{
			Real:  true,
			Name:  "43.txt.utf-8.bz2",
			Index: 2,
		},
		{
			Real:  true,
			Name:  "pg74.txt.bz2",
			Index: 3,
		},
		{
			Real:  true,
			Name:  "76.txt.utf-8.bz2",
			Index: 4,
		},
		{
			Real:  true,
			Name:  "84.txt.utf-8.bz2",
			Index: 5,
		},
		{
			Real:  true,
			Name:  "100.txt.utf-8.bz2",
			Index: 6,
		},
		{
			Real:  true,
			Name:  "145.txt.utf-8.bz2",
			Index: 7,
		},
		{
			Real:  true,
			Name:  "768.txt.utf-8.bz2",
			Index: 8,
		},
		{
			Real:  true,
			Name:  "1260.txt.utf-8.bz2",
			Index: 9,
		},
		{
			Real:  true,
			Name:  "1342.txt.utf-8.bz2",
			Index: 10,
		},
		{
			Real:  true,
			Name:  "1837.txt.utf-8.bz2",
			Index: 11,
		},
		{
			Real:  true,
			Name:  "2641.txt.utf-8.bz2",
			Index: 12,
		},
		{
			Real:  true,
			Name:  "2701.txt.utf-8.bz2",
			Index: 13,
		},
		{
			Real:  true,
			Name:  "3176.txt.utf-8.bz2",
			Index: 14,
		},
		{
			Real:  true,
			Name:  "37106.txt.utf-8.bz2",
			Index: 15,
		},
		{
			Real:  true,
			Name:  "64317.txt.utf-8.bz2",
			Index: 16,
		},
		{
			Real:  true,
			Name:  "67979.txt.utf-8.bz2",
			Index: 17,
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
		{
			Real:  true,
			Name:  "lm.txt.bz2",
			Index: 21,
		},
		{
			Real:  true,
			Name:  "accelerando.txt.bz2",
			Index: 22,
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

const gemini = `A shimmering twilight permanently blankets the Cerulean Hinterlands, a vast, primordial basin where nature defies conventional biology. The ground beneath is a spongy, resilient carpet of bioluminescent moss that pulses in tandem with a slow, rhythmic planetary heartbeat. Towering above this neon floor are the Goliath Redwoods, colossal flora whose bark resembles liquid obsidian, reflecting the twin moons that hang suspended in a violet sky. Instead of leaves, these silent giants sprout iridescent, translucent crystalline fronds that chime like distant glass wind chimes whenever the atmospheric thermal currents sweep through the valley.

Meandering through the heart of the hinterlands is the Whisperwind River, a stream of liquid quicksilver that flows uphill, defying gravity by climbing the tiered obsidian terraces. The water glows with a soft, internal amber warmth, casting dancing shadows on the surrounding stone formations. Flocks of featherless, moth-winged avians dance above the water's surface, leaving trails of stardust in their wake.

The air is thick with the sweet, crisp scent of crushed ozone and wild, oversized vanilla orchids that bloom only in the shadows. There are no harsh winds here, only a perpetual, comforting breeze that carries the distant, melodic echoes of the valley’s deep caverns. It is a sanctuary of surreal stillness, where the boundary between organic life and mineral magic blurs entirely, creating an untamed wilderness that feels both anciently grounded and beautifully alien.`

// Text is a text type
type Text uint8

const (
	TextSlop Text = iota
	TextNot
)

// Rank is a text rank
type Rank struct {
	Rank [2]float64
	Type Text
}

var Ranks = []Rank{
	/*{[2]float64{0.0027188210448006, 0.0021707177135268}, TextNot},
	{[2]float64{0.0024715820022209, 0.0022843656070799}, TextNot},
	{[2]float64{0.0027671081299420, 0.0023964612455301}, TextNot},
	{[2]float64{0.0026193556997382, 0.0021191478071765}, TextNot},
	{[2]float64{0.0032537747798464, 0.0022569680064573}, TextNot},
	{[2]float64{0.0025840830287612, 0.0022436937998044}, TextNot},
	{[2]float64{0.0030199160180444, 0.0025708445147745}, TextNot},
	{[2]float64{0.0030208543655926, 0.0026152698951034}, TextNot},
	{[2]float64{0.0026994349599398, 0.0022076699527420}, TextNot},
	{[2]float64{0.0030260218133745, 0.0025658719290011}, TextNot},
	{[2]float64{0.0027594223439506, 0.0021298467539631}, TextNot},
	{[2]float64{0.0024739898736685, 0.0023793495891937}, TextNot},
	{[2]float64{0.0029726773896805, 0.0026090682614000}, TextNot},
	{[2]float64{0.0042592065200255, 0.0030658209062800}, TextNot},
	{[2]float64{0.0026546765932896, 0.0021784202039178}, TextNot},
	{[2]float64{0.0025918806265325, 0.0023566801635863}, TextNot},
	{[2]float64{0.0026971659298036, 0.0022283595796568}, TextNot},
	{[2]float64{0.0025590231275894, 0.0024211534646999}, TextNot},
	{[2]float64{0.0026751595676999, 0.0020610967936296}, TextNot},
	{[2]float64{0.0025289134565670, 0.0024787705252416}, TextNot},
	{[2]float64{0.0023368627091383, 0.0018932748412559}, TextNot},
	{[2]float64{0.0024870842928849, 0.0024176143717155}, TextNot},
	{[2]float64{0.0025814640050012, 0.0022536222561623}, TextNot},
	{[2]float64{0.0028243499031931, 0.0024085641545411}, TextNot},
	{[2]float64{0.0029075499031191, 0.0023319154135278}, TextNot},
	{[2]float64{0.0025487047367912, 0.0020746365999183}, TextNot},
	{[2]float64{0.0022983617071924, 0.0020691707063021}, TextNot},
	{[2]float64{0.0027792506686523, 0.0022741979528189}, TextNot},
	{[2]float64{0.0023944620415836, 0.0020980970592896}, TextNot},
	{[2]float64{0.0022579607558760, 0.0026524501110297}, TextNot},
	{[2]float64{0.0033354077537305, 0.0024943948065777}, TextNot},
	{[2]float64{0.0026267206446545, 0.0022779221093241}, TextNot},
	{[2]float64{0.0021906492701910, 0.0018664165974540}, TextNot},
	{[2]float64{0.0024939065349125, 0.0021234736726677}, TextNot},
	{[2]float64{0.0022798844561389, 0.0019426968391971}, TextNot},
	{[2]float64{0.0030600140400914, 0.0032287975459727}, TextNot},
	{[2]float64{0.0029258282744318, 0.0024894356274089}, TextNot},
	{[2]float64{0.0028060770270868, 0.0023128213538140}, TextNot},
	{[2]float64{0.0025397396367300, 0.0022325897122978}, TextNot},
	{[2]float64{0.0025992219533246, 0.0022279109831160}, TextNot},
	{[2]float64{0.0027802379613480, 0.0023351546670216}, TextNot},
	{[2]float64{0.0022546838432058, 0.0023215665772443}, TextNot},
	{[2]float64{0.0030133758611835, 0.0026073213972252}, TextNot},
	{[2]float64{0.0024367762075795, 0.0021641497698719}, TextNot},
	{[2]float64{0.0026651291875348, 0.0024166674588935}, TextNot},
	{[2]float64{0.0030995077704995, 0.0024482251131621}, TextNot},
	{[2]float64{0.0023960381190513, 0.0019835105408201}, TextNot},
	{[2]float64{0.0025965666477240, 0.0021773752708994}, TextNot},
	{[2]float64{0.0026769202985166, 0.0025215607345782}, TextNot},
	{[2]float64{0.0029379048771785, 0.0023425463694939}, TextNot},
	{[2]float64{0.0028795887644575, 0.0025370055632383}, TextNot},
	{[2]float64{0.0030094657475025, 0.0027909382211062}, TextNot},
	{[2]float64{0.0025344136713206, 0.0021558403024649}, TextNot},
	{[2]float64{0.0026219374546689, 0.0021517113313310}, TextNot},
	{[2]float64{0.0028865071627312, 0.0025492552184272}, TextNot},
	{[2]float64{0.0025604733845274, 0.0022736778704723}, TextNot},
	{[2]float64{0.0028319165180598, 0.0025514659201326}, TextNot},
	{[2]float64{0.0026622015240734, 0.0021657275843487}, TextNot},
	{[2]float64{0.0035105711405213, 0.0023411923361343}, TextNot},
	{[2]float64{0.0027806458982686, 0.0025126437639720}, TextNot},
	{[2]float64{0.0027587114059970, 0.0024436857058114}, TextNot},
	{[2]float64{0.0027085489113351, 0.0022838944749878}, TextNot},
	{[2]float64{0.0027395088538209, 0.0022410249216797}, TextNot},
	{[2]float64{0.0027325167866255, 0.0024290516533453}, TextNot},
	{[2]float64{0.0025149471576455, 0.0024180855017632}, TextSlop},
	{[2]float64{0.0024445933956482, 0.0024937129839479}, TextSlop},
	{[2]float64{0.0023339019148041, 0.0024091981326603}, TextSlop},
	{[2]float64{0.0024220115114023, 0.0024734012048123}, TextSlop},
	{[2]float64{0.0026346499019590, 0.0027190505596740}, TextSlop},
	{[2]float64{0.0017118832908543, 0.0020221873790233}, TextSlop},
	{[2]float64{0.0019824144576281, 0.0019871751676600}, TextSlop},
	{[2]float64{0.0015424269288494, 0.0019297543163752}, TextSlop},
	{[2]float64{0.0020004694725476, 0.0022205196071603}, TextSlop},
	{[2]float64{0.0023131888268327, 0.0024205198865430}, TextSlop},
	{[2]float64{0.0024191015400646, 0.0023294245100570}, TextSlop},
	{[2]float64{0.0024489039635493, 0.0024009054204522}, TextSlop},
	{[2]float64{0.0024318700824397, 0.0025913612141239}, TextSlop},
	{[2]float64{0.0021308244593731, 0.0023762324464568}, TextSlop},
	{[2]float64{0.0026200062738841, 0.0028136367008170}, TextSlop},
	{[2]float64{0.0025487131930542, 0.0026058196577575}, TextSlop},
	{[2]float64{0.0025411085842163, 0.0026570683085240}, TextSlop},
	{[2]float64{0.0027722206828839, 0.0029362074069110}, TextSlop},
	{[2]float64{0.0025408850862659, 0.0025782601573532}, TextSlop},
	{[2]float64{0.0026262351013152, 0.0028439701041294}, TextSlop},
	{[2]float64{0.0023466259895508, 0.0024142207744834}, TextSlop},
	{[2]float64{0.0023265712402190, 0.0025738195853859}, TextSlop},
	{[2]float64{0.0020511025331169, 0.0024751187759248}, TextSlop},
	{[2]float64{0.0022505472938499, 0.0024848456272955}, TextSlop},
	{[2]float64{0.0023367264173967, 0.0025625272925423}, TextSlop},
	{[2]float64{0.0024414392046835, 0.0025781801901391}, TextSlop},
	{[2]float64{0.0025550910644093, 0.0027009314865210}, TextSlop},
	{[2]float64{0.0024054167664316, 0.0024654532541205}, TextSlop},
	{[2]float64{0.0015494579987679, 0.0020116032119858}, TextSlop},
	{[2]float64{0.0022083485785579, 0.0021456996779022}, TextSlop},
	{[2]float64{0.0025623734728518, 0.0026475525531631}, TextSlop},
	{[2]float64{0.0022563102793485, 0.0023468475209205}, TextSlop},
	{[2]float64{0.0024874132527845, 0.0023535562688416}, TextSlop},
	{[2]float64{0.0024047174992061, 0.0025192893013401}, TextSlop},
	{[2]float64{0.0025699314925276, 0.0026921091495425}, TextSlop},
	{[2]float64{0.0024844674078830, 0.0024600745234944}, TextSlop},
	{[2]float64{0.0024831399177127, 0.0025620861030160}, TextSlop},
	{[2]float64{0.0021898323460208, 0.0024379108415762}, TextSlop},
	{[2]float64{0.0021142262108663, 0.0022580771184123}, TextSlop},
	{[2]float64{0.0019836141218865, 0.0022558896961813}, TextSlop},
	{[2]float64{0.0023170101015378, 0.0022806461830852}, TextSlop},
	{[2]float64{0.0022166256740873, 0.0023195168969092}, TextSlop},
	{[2]float64{0.0026414136916751, 0.0027393168390889}, TextSlop},
	{[2]float64{0.0024520837641847, 0.0024813832485300}, TextSlop},
	{[2]float64{0.0024531369435419, 0.0026076107217002}, TextSlop},
	{[2]float64{0.0022909557569519, 0.0022879198848630}, TextSlop},
	{[2]float64{0.0024855228483818, 0.0025300244613007}, TextSlop},
	{[2]float64{0.0021671147561427, 0.0024029837028294}, TextSlop},
	{[2]float64{0.0020329115799536, 0.0019898823108718}, TextSlop},
	{[2]float64{0.0021475658591521, 0.0021566988284742}, TextSlop},
	{[2]float64{0.0026884490230808, 0.0028583502464631}, TextSlop},
	{[2]float64{0.0025201338997538, 0.0028413143359485}, TextSlop},
	{[2]float64{0.0024110046547876, 0.0025791073744682}, TextSlop},
	{[2]float64{0.0022355475692285, 0.0024845382145101}, TextSlop},
	{[2]float64{0.0029593567459492, 0.0029612084025578}, TextSlop},
	{[2]float64{0.0024822767049441, 0.0026270851279060}, TextSlop},
	{[2]float64{0.0018238113597638, 0.0020914255523993}, TextSlop},
	{[2]float64{0.0021154275181101, 0.0022610911266189}, TextSlop},
	{[2]float64{0.0022753110504954, 0.0025180764112618}, TextSlop},
	{[2]float64{0.0021692623494004, 0.0022429177526345}, TextSlop},
	{[2]float64{0.0024141895208603, 0.0025486852722917}, TextSlop},
	{[2]float64{0.0025454933226278, 0.0026603229520664}, TextSlop},
	{[2]float64{0.0017224302827641, 0.0021212288055982}, TextSlop},
	{[2]float64{0.0018532019733349, 0.0022807511625862}, TextSlop},*/
	/*{[2]float64{0.0027188210448006, 0.0021282287660669}, TextNot},
	{[2]float64{0.0026193556997382, 0.0024044173312497}, TextNot},
	{[2]float64{0.0024715820022209, 0.0021212090382835}, TextNot},
	{[2]float64{0.0027671081299420, 0.0025811905600564}, TextNot},
	{[2]float64{0.0025840830287612, 0.0025292737536941}, TextNot},
	{[2]float64{0.0026994349599398, 0.0021494315019226}, TextNot},
	{[2]float64{0.0032537747798464, 0.0022762730406897}, TextNot},
	{[2]float64{0.0030260218133745, 0.0029627821842587}, TextNot},
	{[2]float64{0.0030208543655926, 0.0025083438711416}, TextNot},
	{[2]float64{0.0030199160180444, 0.0026504767698154}, TextNot},
	{[2]float64{0.0029726773896805, 0.0027246100401880}, TextNot},
	{[2]float64{0.0026546765932896, 0.0023973491864544}, TextNot},
	{[2]float64{0.0042592065200255, 0.0030226055237778}, TextNot},
	{[2]float64{0.0024739898736685, 0.0021394438167413}, TextNot},
	{[2]float64{0.0027594223439506, 0.0026887096324170}, TextNot},
	{[2]float64{0.0025918806265325, 0.0024274177640227}, TextNot},
	{[2]float64{0.0026971659298036, 0.0025417582180633}, TextNot},
	{[2]float64{0.0026751595676999, 0.0025348115007655}, TextNot},
	{[2]float64{0.0025590231275894, 0.0025481693569581}, TextNot},
	{[2]float64{0.0025289134565670, 0.0026012711794089}, TextNot},
	{[2]float64{0.0025814640050012, 0.0022847602763865}, TextNot},
	{[2]float64{0.0023368627091383, 0.0023527487129493}, TextNot},
	{[2]float64{0.0024870842928849, 0.0027899137330198}, TextNot},
	{[2]float64{0.0029075499031191, 0.0025635695974503}, TextNot},
	{[2]float64{0.0025487047367912, 0.0026126659208511}, TextNot},
	{[2]float64{0.0023944620415836, 0.0023723021284361}, TextNot},
	{[2]float64{0.0028243499031931, 0.0026870580056950}, TextNot},
	{[2]float64{0.0027792506686523, 0.0023035596148951}, TextNot},
	{[2]float64{0.0033354077537305, 0.0025572683625419}, TextNot},
	{[2]float64{0.0022983617071924, 0.0022596898368182}, TextNot},
	{[2]float64{0.0026267206446545, 0.0023300995755301}, TextNot},
	{[2]float64{0.0024939065349125, 0.0024329179034035}, TextNot},
	{[2]float64{0.0021906492701910, 0.0023855225557142}, TextNot},
	{[2]float64{0.0022798844561389, 0.0025045502518692}, TextNot},
	{[2]float64{0.0022579607558760, 0.0023367439539443}, TextNot},
	{[2]float64{0.0029258282744318, 0.0026649250898223}, TextNot},
	{[2]float64{0.0027802379613480, 0.0022328361821122}, TextNot},
	{[2]float64{0.0030600140400914, 0.0028403719854826}, TextNot},
	{[2]float64{0.0025397396367300, 0.0023561740984793}, TextNot},
	{[2]float64{0.0028060770270868, 0.0025879352569381}, TextNot},
	{[2]float64{0.0025992219533246, 0.0024416036732517}, TextNot},
	{[2]float64{0.0022546838432058, 0.0019351564479413}, TextNot},
	{[2]float64{0.0024367762075795, 0.0026077335314441}, TextNot},
	{[2]float64{0.0026651291875348, 0.0029192280224200}, TextNot},
	{[2]float64{0.0028795887644575, 0.0027077536442543}, TextNot},
	{[2]float64{0.0023960381190513, 0.0023759071727425}, TextNot},
	{[2]float64{0.0030133758611835, 0.0021616373371680}, TextNot},
	{[2]float64{0.0030995077704995, 0.0022919324414205}, TextNot},
	{[2]float64{0.0025965666477240, 0.0025240358563622}, TextNot},
	{[2]float64{0.0026769202985166, 0.0026807408358443}, TextNot},
	{[2]float64{0.0029379048771785, 0.0025491713250228}, TextNot},
	{[2]float64{0.0030094657475025, 0.0026039956586450}, TextNot},
	{[2]float64{0.0025344136713206, 0.0023030433062537}, TextNot},
	{[2]float64{0.0028319165180598, 0.0024527623448128}, TextNot},
	{[2]float64{0.0025604733845274, 0.0024219066518318}, TextNot},
	{[2]float64{0.0028865071627312, 0.0024051401117624}, TextNot},
	{[2]float64{0.0035105711405213, 0.0024538978203499}, TextNot},
	{[2]float64{0.0027806458982686, 0.0024491288329963}, TextNot},
	{[2]float64{0.0026622015240734, 0.0025886279732304}, TextNot},
	{[2]float64{0.0027587114059970, 0.0026264821485636}, TextNot},
	{[2]float64{0.0026219374546689, 0.0024089574805724}, TextNot},
	{[2]float64{0.0027395088538209, 0.0027303335480118}, TextNot},
	{[2]float64{0.0027085489113351, 0.0026419589000266}, TextNot},
	{[2]float64{0.0027325167866255, 0.0027131712298715}, TextNot},
	{[2]float64{0.0024445933956482, 0.0020767255520293}, TextSlop},
	{[2]float64{0.0024220115114023, 0.0021838305930206}, TextSlop},
	{[2]float64{0.0025149471576455, 0.0020855339490726}, TextSlop},
	{[2]float64{0.0015424269288494, 0.0014159845922331}, TextSlop},
	{[2]float64{0.0017118832908543, 0.0016807806667700}, TextSlop},
	{[2]float64{0.0026346499019590, 0.0022840275405665}, TextSlop},
	{[2]float64{0.0023339019148041, 0.0019350160749709}, TextSlop},
	{[2]float64{0.0023131888268327, 0.0020574931691047}, TextSlop},
	{[2]float64{0.0020004694725476, 0.0018321972983021}, TextSlop},
	{[2]float64{0.0019824144576281, 0.0015764245885660}, TextSlop},
	{[2]float64{0.0024191015400646, 0.0020446624859874}, TextSlop},
	{[2]float64{0.0024318700824397, 0.0020158633102074}, TextSlop},
	{[2]float64{0.0024489039635493, 0.0020840413388966}, TextSlop},
	{[2]float64{0.0025487131930542, 0.0022129795512482}, TextSlop},
	{[2]float64{0.0026200062738841, 0.0022712573138998}, TextSlop},
	{[2]float64{0.0021308244593731, 0.0021417630710741}, TextSlop},
	{[2]float64{0.0025408850862659, 0.0022710599956221}, TextSlop},
	{[2]float64{0.0023466259895508, 0.0020419253531920}, TextSlop},
	{[2]float64{0.0027722206828839, 0.0024426519234148}, TextSlop},
	{[2]float64{0.0026262351013152, 0.0023240574898432}, TextSlop},
	{[2]float64{0.0020511025331169, 0.0019756426006426}, TextSlop},
	{[2]float64{0.0025411085842163, 0.0022603015313247}, TextSlop},
	{[2]float64{0.0023265712402190, 0.0020181353919352}, TextSlop},
	{[2]float64{0.0025550910644093, 0.0020713539731139}, TextSlop},
	{[2]float64{0.0015494579987679, 0.0014954210224641}, TextSlop},
	{[2]float64{0.0022505472938499, 0.0019561162943465}, TextSlop},
	{[2]float64{0.0024414392046835, 0.0021041998322754}, TextSlop},
	{[2]float64{0.0023367264173967, 0.0021307405875052}, TextSlop},
	{[2]float64{0.0025623734728518, 0.0021704819185588}, TextSlop},
	{[2]float64{0.0024054167664316, 0.0021442135757522}, TextSlop},
	{[2]float64{0.0024874132527845, 0.0020536908364958}, TextSlop},
	{[2]float64{0.0022563102793485, 0.0020065473414054}, TextSlop},
	{[2]float64{0.0025699314925276, 0.0022323439902567}, TextSlop},
	{[2]float64{0.0022083485785579, 0.0021980362193465}, TextSlop},
	{[2]float64{0.0024047174992061, 0.0021270170471341}, TextSlop},
	{[2]float64{0.0021898323460208, 0.0020673127765840}, TextSlop},
	{[2]float64{0.0024831399177127, 0.0021638525002218}, TextSlop},
	{[2]float64{0.0022166256740873, 0.0019468676539144}, TextSlop},
	{[2]float64{0.0024844674078830, 0.0019955432066698}, TextSlop},
	{[2]float64{0.0021142262108663, 0.0019658859244526}, TextSlop},
	{[2]float64{0.0019836141218865, 0.0019098311668259}, TextSlop},
	{[2]float64{0.0024520837641847, 0.0020489432750839}, TextSlop},
	{[2]float64{0.0023170101015378, 0.0019910025462555}, TextSlop},
	{[2]float64{0.0024531369435419, 0.0021487895670600}, TextSlop},
	{[2]float64{0.0026414136916751, 0.0023313062883395}, TextSlop},
	{[2]float64{0.0021475658591521, 0.0018878953288244}, TextSlop},
	{[2]float64{0.0024855228483818, 0.0022273031844062}, TextSlop},
	{[2]float64{0.0021671147561427, 0.0020200631311045}, TextSlop},
	{[2]float64{0.0022909557569519, 0.0018680341935259}, TextSlop},
	{[2]float64{0.0020329115799536, 0.0015220188046325}, TextSlop},
	{[2]float64{0.0029593567459492, 0.0023852213685887}, TextSlop},
	{[2]float64{0.0026884490230808, 0.0022668207997967}, TextSlop},
	{[2]float64{0.0024110046547876, 0.0021241398852009}, TextSlop},
	{[2]float64{0.0024822767049441, 0.0020731029732809}, TextSlop},
	{[2]float64{0.0022355475692285, 0.0020400013889371}, TextSlop},
	{[2]float64{0.0025201338997538, 0.0021392070238943}, TextSlop},
	{[2]float64{0.0022753110504954, 0.0019848225050747}, TextSlop},
	{[2]float64{0.0021692623494004, 0.0018108496672041}, TextSlop},
	{[2]float64{0.0018238113597638, 0.0019296174663262}, TextSlop},
	{[2]float64{0.0021154275181101, 0.0017493220061098}, TextSlop},
	{[2]float64{0.0017224302827641, 0.0016407467308337}, TextSlop},
	{[2]float64{0.0024141895208603, 0.0020361376497181}, TextSlop},
	{[2]float64{0.0018532019733349, 0.0019107640659964}, TextSlop},
	{[2]float64{0.0025454933226278, 0.0021929932017456}, TextSlop},*/
	{[2]float64{0.0027188210448006, 0.0020153591888537}, TextNot},
	{[2]float64{0.0024715820022209, 0.0019242743405651}, TextNot},
	{[2]float64{0.0026193556997382, 0.0020576325182144}, TextNot},
	{[2]float64{0.0027671081299420, 0.0021436104376134}, TextNot},
	{[2]float64{0.0025840830287612, 0.0021179512454144}, TextNot},
	{[2]float64{0.0026994349599398, 0.0019191603815500}, TextNot},
	{[2]float64{0.0032537747798464, 0.0020504564297818}, TextNot},
	{[2]float64{0.0030260218133745, 0.0024476575172545}, TextNot},
	{[2]float64{0.0030199160180444, 0.0022939100072987}, TextNot},
	{[2]float64{0.0030208543655926, 0.0023267928549992}, TextNot},
	{[2]float64{0.0026546765932896, 0.0020276422506498}, TextNot},
	{[2]float64{0.0029726773896805, 0.0023720884825085}, TextNot},
	{[2]float64{0.0024739898736685, 0.0020021534311304}, TextNot},
	{[2]float64{0.0042592065200255, 0.0027879466429192}, TextNot},
	{[2]float64{0.0027594223439506, 0.0022052960766603}, TextNot},
	{[2]float64{0.0025918806265325, 0.0020370916877594}, TextNot},
	{[2]float64{0.0026751595676999, 0.0021262862143748}, TextNot},
	{[2]float64{0.0026971659298036, 0.0020096204023475}, TextNot},
	{[2]float64{0.0025590231275894, 0.0021636465402645}, TextNot},
	{[2]float64{0.0025289134565670, 0.0021440158287557}, TextNot},
	{[2]float64{0.0025814640050012, 0.0020813871153066}, TextNot},
	{[2]float64{0.0023368627091383, 0.0018488544880637}, TextNot},
	{[2]float64{0.0024870842928849, 0.0020672730542037}, TextNot},
	{[2]float64{0.0029075499031191, 0.0022647088934315}, TextNot},
	{[2]float64{0.0025487047367912, 0.0020271922075845}, TextNot},
	{[2]float64{0.0023944620415836, 0.0018746246933936}, TextNot},
	{[2]float64{0.0028243499031931, 0.0022525734404524}, TextNot},
	{[2]float64{0.0022983617071924, 0.0019206327207547}, TextNot},
	{[2]float64{0.0033354077537305, 0.0022541349584679}, TextNot},
	{[2]float64{0.0027792506686523, 0.0020609486792167}, TextNot},
	{[2]float64{0.0024939065349125, 0.0020012515773152}, TextNot},
	{[2]float64{0.0026267206446545, 0.0021831201129929}, TextNot},
	{[2]float64{0.0022798844561389, 0.0018429681177720}, TextNot},
	{[2]float64{0.0021906492701910, 0.0019489394124910}, TextNot},
	{[2]float64{0.0022579607558760, 0.0021353920656477}, TextNot},
	{[2]float64{0.0029258282744318, 0.0023917057948012}, TextNot},
	{[2]float64{0.0027802379613480, 0.0020761603593752}, TextNot},
	{[2]float64{0.0025992219533246, 0.0020917123143762}, TextNot},
	{[2]float64{0.0030600140400914, 0.0026367860079752}, TextNot},
	{[2]float64{0.0028060770270868, 0.0021058685250767}, TextNot},
	{[2]float64{0.0022546838432058, 0.0018419211791633}, TextNot},
	{[2]float64{0.0025397396367300, 0.0019959608554870}, TextNot},
	{[2]float64{0.0024367762075795, 0.0021393366220599}, TextNot},
	{[2]float64{0.0030995077704995, 0.0021062503189450}, TextNot},
	{[2]float64{0.0026651291875348, 0.0022583702836640}, TextNot},
	{[2]float64{0.0030133758611835, 0.0021846678876456}, TextNot},
	{[2]float64{0.0028795887644575, 0.0023074893422157}, TextNot},
	{[2]float64{0.0023960381190513, 0.0020341781734990}, TextNot},
	{[2]float64{0.0025965666477240, 0.0021341478340150}, TextNot},
	{[2]float64{0.0026769202985166, 0.0021756915409638}, TextNot},
	{[2]float64{0.0030094657475025, 0.0023922431290126}, TextNot},
	{[2]float64{0.0029379048771785, 0.0021876268812398}, TextNot},
	{[2]float64{0.0025344136713206, 0.0020060859349115}, TextNot},
	{[2]float64{0.0028319165180598, 0.0022471089783557}, TextNot},
	{[2]float64{0.0028865071627312, 0.0022701492631787}, TextNot},
	{[2]float64{0.0035105711405213, 0.0021382688613224}, TextNot},
	{[2]float64{0.0027806458982686, 0.0023052606476845}, TextNot},
	{[2]float64{0.0025604733845274, 0.0020261871650150}, TextNot},
	{[2]float64{0.0026622015240734, 0.0021555421346520}, TextNot},
	{[2]float64{0.0026219374546689, 0.0020311560521464}, TextNot},
	{[2]float64{0.0027587114059970, 0.0022552463239985}, TextNot},
	{[2]float64{0.0027395088538209, 0.0021667070862886}, TextNot},
	{[2]float64{0.0027085489113351, 0.0021957886096007}, TextNot},
	{[2]float64{0.0027325167866255, 0.0023020257826888}, TextNot},
	{[2]float64{0.0024445933956482, 0.0020028273779083}, TextSlop},
	{[2]float64{0.0025149471576455, 0.0019630047290875}, TextSlop},
	{[2]float64{0.0015424269288494, 0.0014880280303932}, TextSlop},
	{[2]float64{0.0024220115114023, 0.0020432889636761}, TextSlop},
	{[2]float64{0.0017118832908543, 0.0016328263291246}, TextSlop},
	{[2]float64{0.0026346499019590, 0.0021974862385253}, TextSlop},
	{[2]float64{0.0023339019148041, 0.0019054348968605}, TextSlop},
	{[2]float64{0.0023131888268327, 0.0019912428105422}, TextSlop},
	{[2]float64{0.0019824144576281, 0.0015434500025006}, TextSlop},
	{[2]float64{0.0020004694725476, 0.0017620451311242}, TextSlop},
	{[2]float64{0.0024191015400646, 0.0018787524730116}, TextSlop},
	{[2]float64{0.0024318700824397, 0.0019668845688196}, TextSlop},
	{[2]float64{0.0024489039635493, 0.0019790738863886}, TextSlop},
	{[2]float64{0.0025487131930542, 0.0020806857700000}, TextSlop},
	{[2]float64{0.0026200062738841, 0.0022377455005621}, TextSlop},
	{[2]float64{0.0025408850862659, 0.0021007701972370}, TextSlop},
	{[2]float64{0.0021308244593731, 0.0019257031320394}, TextSlop},
	{[2]float64{0.0027722206828839, 0.0022823965212420}, TextSlop},
	{[2]float64{0.0025411085842163, 0.0021288430939625}, TextSlop},
	{[2]float64{0.0023466259895508, 0.0019494255513368}, TextSlop},
	{[2]float64{0.0020511025331169, 0.0019569937960823}, TextSlop},
	{[2]float64{0.0023265712402190, 0.0020034248337703}, TextSlop},
	{[2]float64{0.0026262351013152, 0.0022449337549669}, TextSlop},
	{[2]float64{0.0025550910644093, 0.0020621577582183}, TextSlop},
	{[2]float64{0.0024414392046835, 0.0019852673247552}, TextSlop},
	{[2]float64{0.0022505472938499, 0.0019716871821269}, TextSlop},
	{[2]float64{0.0015494579987679, 0.0015092462936354}, TextSlop},
	{[2]float64{0.0024054167664316, 0.0020312009549497}, TextSlop},
	{[2]float64{0.0023367264173967, 0.0020207406535453}, TextSlop},
	{[2]float64{0.0025623734728518, 0.0021231580659022}, TextSlop},
	{[2]float64{0.0024874132527845, 0.0020603244773397}, TextSlop},
	{[2]float64{0.0022083485785579, 0.0018698615741617}, TextSlop},
	{[2]float64{0.0024047174992061, 0.0020807617314118}, TextSlop},
	{[2]float64{0.0022563102793485, 0.0019359366913888}, TextSlop},
	{[2]float64{0.0025699314925276, 0.0021406499123334}, TextSlop},
	{[2]float64{0.0021898323460208, 0.0019957321259814}, TextSlop},
	{[2]float64{0.0024844674078830, 0.0019680385361137}, TextSlop},
	{[2]float64{0.0024831399177127, 0.0020298479462247}, TextSlop},
	{[2]float64{0.0022166256740873, 0.0018058583718967}, TextSlop},
	{[2]float64{0.0019836141218865, 0.0018591032597489}, TextSlop},
	{[2]float64{0.0021142262108663, 0.0018892215989844}, TextSlop},
	{[2]float64{0.0024520837641847, 0.0019906395112097}, TextSlop},
	{[2]float64{0.0023170101015378, 0.0019627772973229}, TextSlop},
	{[2]float64{0.0024531369435419, 0.0020616107154292}, TextSlop},
	{[2]float64{0.0021671147561427, 0.0019382646219550}, TextSlop},
	{[2]float64{0.0024855228483818, 0.0020428985262783}, TextSlop},
	{[2]float64{0.0020329115799536, 0.0015064495019467}, TextSlop},
	{[2]float64{0.0021475658591521, 0.0017487631211274}, TextSlop},
	{[2]float64{0.0026414136916751, 0.0022416830336247}, TextSlop},
	{[2]float64{0.0022909557569519, 0.0018652014194489}, TextSlop},
	{[2]float64{0.0029593567459492, 0.0024085086855719}, TextSlop},
	{[2]float64{0.0025201338997538, 0.0021205778850723}, TextSlop},
	{[2]float64{0.0024110046547876, 0.0020002481239142}, TextSlop},
	{[2]float64{0.0024822767049441, 0.0020771096741988}, TextSlop},
	{[2]float64{0.0022753110504954, 0.0019761822792511}, TextSlop},
	{[2]float64{0.0026884490230808, 0.0022703044772992}, TextSlop},
	{[2]float64{0.0022355475692285, 0.0019397128037616}, TextSlop},
	{[2]float64{0.0021692623494004, 0.0017782344157116}, TextSlop},
	{[2]float64{0.0018238113597638, 0.0017918176575694}, TextSlop},
	{[2]float64{0.0021154275181101, 0.0017458904084874}, TextSlop},
	{[2]float64{0.0018532019733349, 0.0018357947592309}, TextSlop},
	{[2]float64{0.0024141895208603, 0.0020280162245580}, TextSlop},
	{[2]float64{0.0017224302827641, 0.0015185570100682}, TextSlop},
	{[2]float64{0.0025454933226278, 0.0020334444076903}, TextSlop},
}

// Prompt is a llm prompt
type Prompt struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// Query submits a query to the llm
func Query(query, model string) string {
	prompt := Prompt{
		Model:  model,
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
	length := float64(len(set))
	for symbols := range set {
		p += math.Log(float64(t[a].Count[symbols]+1) / (float64(t[a].Total) + float64(len(t[a].Count))))

	}
	return p / length
}

// Markov is a markov state
type Markov [4]byte

// Iterate iterate the symbols
func (m *Markov) Iterate(n byte) {
	for i := range m[:len(m)-1] {
		m[i] = m[i+1]
	}
	m[len(m)-1] = n
}

// Model is a model
type Model struct {
	Book
	Model map[Markov][256]uint64
}

// Lookup looks up a histogram
func (m *Model) Lookup(markov Markov) [256]uint64 {
	for i := range markov {
		h, contains := m.Model[markov]
		if contains {
			return h
		}
		markov[len(markov)-1-i] = 0
	}
	return m.Model[markov]
}

// MarkovMode markov mode
func MarkovMode() {
	books := LoadBooks()
	models := make([]Model, 0, 8)
	for _, book := range books {
		if book.Real {
			model := Model{
				Book:  book,
				Model: make(map[Markov][256]uint64),
			}
			markov := Markov{}
			for _, symbol := range model.Text {
				current := markov
				for i := range current {
					histogram := model.Model[current]
					histogram[symbol]++
					model.Model[current] = histogram
					current[len(current)-1-i] = 0
				}
				histogram := model.Model[current]
				histogram[symbol]++
				model.Model[current] = histogram
				markov.Iterate(symbol)
			}
			models = append(models, model)
		}
	}

	dot := func(a, b [256]uint64) float64 {
		sum := 0.0
		for i, value := range a {
			sum += float64(value) * float64(b[i])
		}
		return sum
	}
	cs := func(a, b [256]uint64) float64 {
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
	entropy := func(a [256]uint64) float64 {
		sum := 0.0
		for _, value := range a {
			sum += float64(value)
		}
		entropy := 0.0
		for _, value := range a {
			p := float64(value) / sum
			entropy += p * math.Log2(p)
		}
		return -entropy
	}

	markov := Markov{}
	context := Model{
		Model: make(map[Markov][256]uint64),
	}
	ctxt := []byte("What is the meaning of life?")
	for _, symbol := range ctxt {
		current := markov
		for i := range current {
			histogram := context.Model[current]
			histogram[symbol]++
			context.Model[current] = histogram
			current[len(current)-1-i] = 0
		}
		histogram := context.Model[current]
		histogram[symbol]++
		context.Model[current] = histogram
		markov.Iterate(symbol)
	}

	rng := rand.New(rand.NewSource(1))
	for range 8 * 1024 {
		set := make([][256]uint64, 0, 8)
		for i := range models {
			set = append(set, models[i].Lookup(markov))
		}
		//target := context.Lookup(markov)
		min, index := math.MaxFloat64, 0
		_ = cs
		for i := range set {
			//s := cs(set[i], target)
			s := entropy(set[i])
			if s < min {
				min, index = s, i
			}
		}
		sum := uint64(0)
		for _, count := range set[index] {
			sum += count
		}
		symbol, selected, total := byte(0), uint64(rng.Intn(int(sum))), uint64(0)
		for i, value := range set[index] {
			total += value
			if selected < total {
				symbol = byte(i)
				break
			}
		}
		fmt.Printf("%c", symbol)

		current := markov
		for i := range current {
			histogram := context.Model[current]
			histogram[symbol]++
			context.Model[current] = histogram
			current[len(current)-1-i] = 0
		}
		histogram := context.Model[current]
		histogram[symbol]++
		context.Model[current] = histogram
		markov.Iterate(symbol)
	}
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
	// FlagSample generates some samples
	FlagSample = flag.Bool("sample", false, "generate samples")
	// FlagMarkov markov mode
	FlagMarkov = flag.Bool("markov", false, "markov mode")
	// FlagGraph graphical model
	FlagGraph = flag.Bool("graph", false, "graphical model")
	// FlagVerse generate text
	FlagVerse = flag.String("verse", "", "generate text")
	// FlagPre pre-generate model
	FlagPre = flag.Bool("pre", false, "pre-generate model")
	// FlagCal calibrate
	FlagCal = flag.Bool("cal", false, "calibrate")
	// FlagTest test
	FlagTest = flag.Int("test", -1, "test")
)

// NNMode is the nearest neighbor mode
func NNMode() {
	books := LoadBooks()
	a, b := books[4].Text[9*1024:10*1024], books[5].Text[8*1024:9*1024]
	fake := []byte(gemini[:1024])
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
	Diff  map[string]uint64
}

// NewGraph makes a new graph
func NewGraph() Graph {
	return Graph{
		Graph: make(map[string]Node),
		Ranks: make(map[string]uint64),
	}
}

// Learn learns a model
func (g *Graph) Learn(iterations int, rng *rand.Rand, words []string) {
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
	for range iterations {
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

// Add adds context to a model
func (g *Graph) Add(iterations int, rng *rand.Rand, words []string) {
	sum, count := uint64(0), uint64(0)
	for _, key := range g.Keys {
		node := g.Graph[key]
		for _, key := range node.Keys {
			sum += node.Links[key]
			count++
		}
	}
	avg := float64(sum) / float64(count)
	stddev := 0.0
	for _, key := range g.Keys {
		node := g.Graph[key]
		for _, key := range node.Keys {
			diff := float64(node.Links[key]) - avg
			stddev += diff * diff
		}
	}
	stddev = math.Sqrt(stddev / float64(count))
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
			count += uint64(3 * stddev)
			node.Links[words[i+1]] = count
			g.Graph[word] = node
		}
	}
	if g.Diff == nil {
		g.Diff = make(map[string]uint64)
	}
	word := words[0]
	node := g.Graph[word]
	for range iterations {
		g.Diff[word]++
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
	}
	for key, value := range g.Diff {
		diff := int(value) - int(g.Ranks[key])
		if diff < 0 {
			diff = -diff
		}
		g.Diff[key] = uint64(diff)
	}
}

// Result is a result
type Result struct {
	Graph
	List []string
}

// GraphResults graph results
type GraphResults struct {
	A     chan Graph
	B     chan Result
	C     chan Result
	Text  int
	Books []Book
}

// Process process results
func (g GraphResults) Process() {
	rng := rand.New(rand.NewSource(1))
	text := string(g.Books[g.Text].Text)
	words := strings.Fields(text)

	gA := <-g.A
	gB := <-g.B
	gC := <-g.C

	fmt.Println(g.Books[g.Text].Name)
	word := "God"
	node := gA.Graph[word]
	for range 33 {
		fmt.Printf(" %s", word)
		sum := uint64(0)
		for _, w := range node.Keys {
			sum += gA.Ranks[w]
		}
		for sum == 0 {
			node = gA.Graph[words[rng.Intn(len(words))]]
			for _, w := range node.Keys {
				sum += gA.Ranks[w]
			}
		}
		total, selected := uint64(0), uint64(rng.Intn(int(sum)))
		for _, w := range node.Keys {
			total += gA.Ranks[w]
			if selected < total {
				word = w
				node = gA.Graph[word]
				break
			}
		}
	}
	fmt.Println()
	sum := uint64(0)
	for _, value := range gA.Keys {
		sum += gA.Ranks[value]
	}
	entropy := 0.0
	for _, value := range gA.Keys {
		if gA.Ranks[value] == 0 {
			continue
		}
		p := float64(gA.Ranks[value]) / float64(sum)
		entropy += p * math.Log2(p)
	}
	fmt.Println(-entropy)
	process := func(r Result) (result float64) {
		sum := uint64(0)
		for _, value := range r.Keys {
			sum += r.Ranks[value]
		}
		entropy := 0.0
		for _, value := range r.Keys {
			if r.Ranks[value] == 0 {
				continue
			}
			p := float64(r.Ranks[value]) / float64(sum)
			entropy += p * math.Log2(p)
		}
		fmt.Println(-entropy)
		{
			max := 0.0
			for range r.List {
				p := 1 / float64(len(r.List))
				max += p * math.Log2(p)
			}
			fmt.Println("max", -max)
			sum := uint64(0)
			for _, value := range r.List {
				sum += r.Ranks[value]
			}
			entropy := 0.0
			for _, value := range r.List {
				if r.Ranks[value] == 0 {
					continue
				}
				p := float64(r.Ranks[value]) / float64(sum)
				entropy += p * math.Log2(p)
			}
			fmt.Println(-entropy)
		}
		{
			sum := uint64(0)
			for _, value := range r.List {
				sum += r.Ranks[value]
			}
			result = float64(sum) / float64(len(r.List))
			fmt.Println(result)
		}
		{
			sum := uint64(0)
			for _, value := range words {
				sum += gB.Ranks[value]
			}
			fmt.Println(float64(sum) / float64(len(words)))
		}
		return result
	}
	rB := process(gB)
	rC := process(gC)
	if g.Books[g.Text].Real {
		if rB < rC {
			fmt.Println("correct real", rB, rC)
		} else {
			fmt.Println("incorrect real", rB, rC)
		}
	} else {
		if rB > rC {
			fmt.Println("correct fake", rB, rC)
		} else {
			fmt.Println("incorrect fake", rB, rC)
		}
	}
}

// GraphMode is a graphical model
func GraphMode(books []Book, t int, alt string) GraphResults {
	text := string(books[t].Text)
	words := strings.Fields(text)
	doneA := make(chan Graph, 8)
	go func() {
		rng := rand.New(rand.NewSource(1))
		g := NewGraph()
		g.Learn(8*1024*1024, rng, words)
		doneA <- g
	}()

	doneB := make(chan Result, 8)
	go func() {
		rng := rand.New(rand.NewSource(1))
		suffix := strings.Fields(samples[0][:1024])
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
		g.Learn(8*1024*1024, rng, words)
		doneB <- Result{
			Graph: g,
			List:  list,
		}
	}()
	doneC := make(chan Result, 8)
	go func() {
		rng := rand.New(rand.NewSource(1))
		suffix := strings.Fields(alt)
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
		g.Learn(8*1024*1024, rng, words)
		doneC <- Result{
			Graph: g,
			List:  list,
		}
	}()
	return GraphResults{
		A:     doneA,
		B:     doneB,
		C:     doneC,
		Books: books,
		Text:  t,
	}
}

// VerseMode generate text
func VerseMode(text string) {
	rng := rand.New(rand.NewSource(1))
	words := strings.Fields(text)
	g := NewGraph()
	input, err := os.Open("pre.gob")
	if err != nil {
		panic(err)
	}
	defer input.Close()
	decoder := gob.NewDecoder(input)
	err = decoder.Decode(&g)
	if err != nil {
		panic(err)
	}
	type Trace struct {
		Trace string
		Value uint64
	}
	set := make([]Trace, 0, 8)
	g.Add(8*1024*1024, rng, words)
	for range 1024 {
		word := words[0]
		node := g.Graph[word]
		trace := Trace{}
		for range 33 {
			trace.Trace = trace.Trace + word + " "
			trace.Value += g.Diff[word]
			sum := uint64(0)
			for _, w := range node.Keys {
				sum += g.Diff[w]
			}
			for sum == 0 {
				node = g.Graph[words[rng.Intn(len(words))]]
				for _, w := range node.Keys {
					sum += g.Diff[w]
				}
			}
			total, selected := uint64(0), uint64(rng.Intn(int(sum)))
			for _, w := range node.Keys {
				total += g.Diff[w]
				if selected < total {
					word = w
					node = g.Graph[word]
					break
				}
			}
		}
		set = append(set, trace)
	}
	sort.Slice(set, func(i, j int) bool {
		return set[i].Value < set[j].Value
	})
	for _, trace := range set {
		fmt.Println(trace.Value, trace.Trace)
	}
}

// PreMode pre-generate model
func PreMode(text string) {
	rng := rand.New(rand.NewSource(1))
	words := strings.Fields(text)
	g := NewGraph()
	g.Learn(8*1024*1024, rng, words)
	output, err := os.Create("pre.gob")
	if err != nil {
		panic(err)
	}
	defer output.Close()
	encoder := gob.NewEncoder(output)
	err = encoder.Encode(g)
	if err != nil {
		panic(err)
	}
}

const (
	SlopAvg    = 0.002318874154124976
	SlopStddev = 5.8209022598764776e-05
	NotAvg     = 0.002729698680992552
	NotStddev  = 3.5226320411986187e-07
)

// CalMode calibrate
func CalMode() {
	const Samples = 64
	books := LoadBooks()
	rng := rand.New(rand.NewSource(1))
	text := string(books[0].Text)
	words := strings.Fields(text)
	words2 := strings.Fields(string(books[22].Text))

	done := make(chan [2]float64, 8)
	cal := func(t string, seed int64, alt string) {
		rng := rand.New(rand.NewSource(1))
		suffix := strings.Fields(alt)
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
		count := g.LearnFast(1e-6, 8*1024*1024, rng, words, list, len(list))
		result := 0.0
		{
			sum := 0.0
			for _, value := range list {
				sum += float64(g.Ranks[value]) / float64(count)
			}
			result = float64(sum) / float64(len(list))
		}
		result2 := 0.0
		{
			suffix := strings.Fields(alt)
			cp := make([]string, len(words2))
			copy(cp, words2)
			has, list := make(map[string]bool), make([]string, 0, 8)
			for _, word := range suffix {
				if !has[word] {
					has[word] = true
					list = append(list, word)
				}
			}
			words := append(cp, suffix...)
			g := NewGraph()
			count := g.LearnFast(1e-6, 8*1024*1024, rng, words, list, len(list))
			{
				sum := 0.0
				for _, value := range list {
					sum += float64(g.Ranks[value]) / float64(count)
				}
				result2 = float64(sum) / float64(len(list))
			}
		}
		fmt.Printf("{[2]float64{%.16f, %.16f}, %s},\n", result, result2, t)
		done <- [2]float64{result, result2}
	}
	c, flight, cpus := 0, 0, runtime.NumCPU()
	for c < Samples && flight < cpus {
		book := rng.Intn(18)
		length := len(books[book].Text)
		count := length/1024 - 1
		index := rng.Intn(count)
		go cal("TextNot", rng.Int63(), string(books[book].Text[index*1024:(index+1)*1024]))
		c++
		flight++
	}
	results := make([][2]float64, 0, Samples)
	for c < Samples {
		result := <-done
		results = append(results, result)
		flight--

		book := rng.Intn(18)
		length := len(books[book].Text)
		count := length/1024 - 1
		index := rng.Intn(count)
		go cal("TextNot", rng.Int63(), string(books[book].Text[index*1024:(index+1)*1024]))
		c++
		flight++
	}
	for range flight {
		result := <-done
		results = append(results, result)
	}

	{
		c, flight, cpus := 0, 0, runtime.NumCPU()
		for c < Samples && flight < cpus {
			book := rng.Intn(3) + 18
			length := len(books[book].Text)
			count := length/1024 - 1
			index := rng.Intn(count)
			go cal("TextSlop", rng.Int63(), string(books[book].Text[index*1024:(index+1)*1024]))
			c++
			flight++
		}
		results := make([][2]float64, 0, Samples)
		for c < Samples {
			result := <-done
			results = append(results, result)
			flight--

			book := rng.Intn(3) + 18
			length := len(books[book].Text)
			count := length/1024 - 1
			index := rng.Intn(count)
			go cal("TextSlop", rng.Int63(), string(books[book].Text[index*1024:(index+1)*1024]))
			c++
			flight++
		}
		for range flight {
			result := <-done
			results = append(results, result)
		}

		for i := range 2 {
			sum := 0.0
			for _, value := range results {
				sum += value[i]
			}
			avg := sum / float64(len(results))
			stddev := 0.0
			for _, value := range results {
				diff := value[i] - avg
				stddev = diff * diff
			}
			stddev = math.Sqrt(stddev / float64(len(results)))
			fmt.Println("slop", avg, stddev)
		}
	}

	for i := range 2 {
		sum := 0.0
		for _, value := range results {
			sum += value[i]
		}
		avg := sum / float64(len(results))
		stddev := 0.0
		for _, value := range results {
			diff := value[i] - avg
			stddev = diff * diff
		}
		stddev = math.Sqrt(stddev / float64(len(results)))
		fmt.Println("not", avg, stddev)
	}
}

// TestMode test
func TestMode() {
	rng := rand.New(rand.NewSource(1))

	context := gradient.Context[float64]{}
	set := context.NewSet()
	set.Add("w0", 2, 8)
	set.AddBias("b0", 8)
	set.Add("w1", 16, 2)
	set.AddBias("b1", 2)
	set.AddData("inputs", 2, len(Ranks))
	set.AddData("in", 2)
	set.AddData("outputs", 2, len(Ranks))
	set.InitAdam(rng)

	inputs := set.ByName["inputs"]
	outputs := set.ByName["outputs"]
	for i, rank := range Ranks {
		for ii := range rank.Rank {
			inputs.X[2*i+ii] = rank.Rank[ii]
		}
		if rank.Type == TextSlop {
			outputs.X[2*i+0] = 1.0
			outputs.X[2*i+1] = 0.0
		} else {

			outputs.X[2*i+0] = 0.0
			outputs.X[2*i+1] = 1.0
		}
	}

	Add := context.B(context.Add)
	Mul := context.B(context.Mul)
	Everett := context.U(context.Everett)
	//Sigmoid := context.U(context.Sigmoid)
	Quadratic := context.B(context.Quadratic)
	Avg := context.U(context.Avg)

	l0 := Everett(Add(Mul(set.Get("w0"), set.Get("inputs")), set.Get("b0")))
	l1 := Add(Mul(set.Get("w1"), l0), set.Get("b1"))
	loss := Avg(Quadratic(l1, set.Get("outputs")))

	for range 16 * 1024 {
		set.Zero()
		l := gradient.Gradient(loss)
		set.Adam(gradient.B1, gradient.B2, .01)
		fmt.Println(l.X[0])
	}

	books := LoadBooks()
	text := string(books[0].Text)
	words := strings.Fields(text)
	words2 := strings.Fields(string(books[22].Text))
	{
		suffix := strings.Fields(samples[*FlagTest][:1024])
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
		var count2 float64
		g2 := NewGraph()
		list2 := make([]string, 0, 8)
		{
			suffix := strings.Fields(samples[*FlagTest][:1024])
			cp := make([]string, len(words2))
			copy(cp, words2)
			has := make(map[string]bool)
			for _, word := range suffix {
				if !has[word] {
					has[word] = true
					list2 = append(list2, word)
				}
			}
			words := append(cp, suffix...)
			count2 = g2.LearnFast(1e-5, 8*1024*1024, rng, words, list2, len(list2))
		}
		{
			sum := 0.0
			for _, value := range list {
				sum += float64(g.Ranks[value]) / float64(count)
			}
			sum2 := 0.0
			for _, value := range list2 {
				sum2 += float64(g2.Ranks[value]) / float64(count2)
			}
			result := float64(sum) / float64(len(list))
			result2 := float64(sum2) / float64(len(list2))
			{
				inputs := set.ByName["in"]
				inputs.X[0] = result
				inputs.X[1] = result2

				Add := context.B(context.Add)
				Mul := context.B(context.Mul)
				Everett := context.U(context.Everett)
				//Sigmoid := context.U(context.Sigmoid)

				l0 := Everett(Add(Mul(set.Get("w0"), set.Get("in")), set.Get("b0")))
				l1 := Add(Mul(set.Get("w1"), l0), set.Get("b1"))
				l1(func(a *gradient.V[float64]) bool {
					fmt.Println("l1", a.X[0], a.X[1])
					return true
				})
			}
			type Result struct {
				Rank
				Diff float64
			}
			results := make([]Result, 0, len(Ranks))
			for _, rank := range Ranks {
				diff := math.Abs(rank.Rank[0] - result)
				//diff2 := rank.Rank[1] - result2
				results = append(results, Result{
					Rank: rank,
					Diff: diff,
				})
			}
			sort.Slice(results, func(i, j int) bool {
				return results[i].Diff < results[j].Diff
			})
			var histogram [2]int
			for i := range results[:32] {
				histogram[results[i].Type]++
			}
			fmt.Println(histogram)
			{
				results := make([]Result, 0, len(Ranks))
				for _, rank := range Ranks {
					diff := math.Abs(rank.Rank[1] - result2)
					results = append(results, Result{
						Rank: rank,
						Diff: diff,
					})
				}
				sort.Slice(results, func(i, j int) bool {
					return results[i].Diff < results[j].Diff
				})
				var histogram [2]int
				for i := range results[:32] {
					histogram[results[i].Type]++
				}
				fmt.Println(histogram)
			}

			fmt.Printf("result=%.16f result2=%.16f\n", result, result2)
			fmt.Println((1+math.Erf((result-SlopAvg)/(SlopStddev*math.Sqrt(2))))/2,
				(1+math.Erf((result-NotAvg)/(NotStddev*math.Sqrt(2))))/2)
			for i := 1; i < 4; i++ {
				fmt.Printf("%d %.16f\n", i, NotAvg-float64(i)*NotStddev)
			}
			if result < NotAvg {
				fmt.Println("fake")
			}
		}
	}
}

func main() {
	flag.Parse()

	if *FlagNN {
		NNMode()
		return
	}

	if *FlagQuery != "" {
		fmt.Println(Query(*FlagQuery, *FlagModel))
		return
	}

	if *FlagGenerate {
		rng := rand.New(rand.NewSource(1))
		results := Query("What is the meaning of life? Be verbose in your answer.", *FlagModel)
		fmt.Println(results)
		for {
			words := strings.Fields(results)
			next := words[rng.Intn(len(words))]
			next = strings.ToLower(strings.Trim(next, ".!?,"))
			results = Query(fmt.Sprintf("What is the meaning of %s? Be verbose in your answer.", next), *FlagModel)
			fmt.Println(results)
		}
	}

	if *FlagSample {
		output, err := os.Create("samples.go")
		if err != nil {
			panic(err)
		}
		defer output.Close()
		fmt.Fprintf(output, `// Copyright 2026 The Sieve Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
		
package main

var samples = []string{
`)

		models := []string{"gemma4", "gpt-oss", "llama3.1:8b"}
		for _, model := range models {
			results := Query("Describe a fictional scene in at least 1024 symbols.", model)
			fmt.Fprintf(output, "`%s`,\n", results)
		}
		fmt.Fprintf(output, "}")
		return
	}

	if *FlagMarkov {
		MarkovMode()
		return
	}

	if *FlagGraph {
		books := LoadBooks()
		g := GraphMode(books, 21, string(books[1].Text[8*1024:9*1024]))
		g.Process()
		g = GraphMode(books, 18, string(books[1].Text[8*1024:9*1024]))
		g.Process()
		g = GraphMode(books, 19, string(books[1].Text[8*1024:9*1024]))
		g.Process()
		g = GraphMode(books, 20, string(books[1].Text[8*1024:9*1024]))
		g.Process()
		return
	}

	if *FlagVerse != "" {
		VerseMode(*FlagVerse)
		return
	}

	if *FlagPre {
		books := LoadBooks()
		PreMode(string(books[0].Text))
		return
	}

	if *FlagCal {
		CalMode()
		return
	}

	if *FlagTest >= 0 {
		TestMode()
		return
	}

	books := LoadBooks()
	data := [][]byte{
		books[4].Text[9*1024 : 10*1024],
		books[5].Text[8*1024 : 9*1024],
		[]byte(samples[0])[:1024],
		[]byte(samples[1])[:1024],
		[]byte(samples[2])[:1024],
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
			symbols.Iterate(symbol)
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
