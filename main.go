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

// BookCount the number of books to use
const BookCount = 4

// BookSet the set of books to use
var BookSet = [BookCount]int{0, 18, 19, 20}

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
	Rank [BookCount]float64
	Type Text
}

var Ranks = []Rank{
	/*{[2]float64{0.0027188210448006, 0.0020126947761568}, TextNot},
	{[2]float64{0.0026193556997382, 0.0020685603897593}, TextNot},
	{[2]float64{0.0024715820022209, 0.0019123002569783}, TextNot},
	{[2]float64{0.0027671081299420, 0.0021477355714208}, TextNot},
	{[2]float64{0.0025840830287612, 0.0021372555593250}, TextNot},
	{[2]float64{0.0032537747798464, 0.0020615617711806}, TextNot},
	{[2]float64{0.0026994349599398, 0.0019261037728857}, TextNot},
	{[2]float64{0.0030199160180444, 0.0022990827118585}, TextNot},
	{[2]float64{0.0030208543655926, 0.0023278603327049}, TextNot},
	{[2]float64{0.0030260218133745, 0.0024490198119427}, TextNot},
	{[2]float64{0.0024739898736685, 0.0020164290349175}, TextNot},
	{[2]float64{0.0026546765932896, 0.0020332011394520}, TextNot},
	{[2]float64{0.0042592065200255, 0.0027719952123388}, TextNot},
	{[2]float64{0.0029726773896805, 0.0023629766121976}, TextNot},
	{[2]float64{0.0027594223439506, 0.0022106832202279}, TextNot},
	{[2]float64{0.0025918806265325, 0.0020510988128279}, TextNot},
	{[2]float64{0.0026971659298036, 0.0020137907831607}, TextNot},
	{[2]float64{0.0026751595676999, 0.0021136714257591}, TextNot},
	{[2]float64{0.0025590231275894, 0.0021666657580400}, TextNot},
	{[2]float64{0.0025289134565670, 0.0021540421433440}, TextNot},
	{[2]float64{0.0023368627091383, 0.0018457259376887}, TextNot},
	{[2]float64{0.0024870842928849, 0.0020619538432607}, TextNot},
	{[2]float64{0.0025814640050012, 0.0020900943751812}, TextNot},
	{[2]float64{0.0029075499031191, 0.0022693704089688}, TextNot},
	{[2]float64{0.0027792506686523, 0.0020802717712846}, TextNot},
	{[2]float64{0.0023944620415836, 0.0018747283271312}, TextNot},
	{[2]float64{0.0025487047367912, 0.0020191031972727}, TextNot},
	{[2]float64{0.0028243499031931, 0.0022662518618300}, TextNot},
	{[2]float64{0.0033354077537305, 0.0022703399064168}, TextNot},
	{[2]float64{0.0022983617071924, 0.0019108721678261}, TextNot},
	{[2]float64{0.0026267206446545, 0.0021581731817740}, TextNot},
	{[2]float64{0.0022798844561389, 0.0018281582651901}, TextNot},
	{[2]float64{0.0024939065349125, 0.0020192390258509}, TextNot},
	{[2]float64{0.0021906492701910, 0.0019705967366372}, TextNot},
	{[2]float64{0.0022579607558760, 0.0021707763478034}, TextNot},
	{[2]float64{0.0029258282744318, 0.0023802622024535}, TextNot},
	{[2]float64{0.0030600140400914, 0.0026314479934123}, TextNot},
	{[2]float64{0.0025992219533246, 0.0020924425449598}, TextNot},
	{[2]float64{0.0028060770270868, 0.0021003592437275}, TextNot},
	{[2]float64{0.0025397396367300, 0.0020157828804711}, TextNot},
	{[2]float64{0.0027802379613480, 0.0020728976838723}, TextNot},
	{[2]float64{0.0022546838432058, 0.0018341888037837}, TextNot},
	{[2]float64{0.0030133758611835, 0.0021951242552140}, TextNot},
	{[2]float64{0.0026651291875348, 0.0022537261217838}, TextNot},
	{[2]float64{0.0024367762075795, 0.0021327558987518}, TextNot},
	{[2]float64{0.0023960381190513, 0.0020321180795489}, TextNot},
	{[2]float64{0.0028795887644575, 0.0023119428501649}, TextNot},
	{[2]float64{0.0030995077704995, 0.0021011386103258}, TextNot},
	{[2]float64{0.0025965666477240, 0.0021410107812784}, TextNot},
	{[2]float64{0.0026769202985166, 0.0021887106251796}, TextNot},
	{[2]float64{0.0030094657475025, 0.0023897481368266}, TextNot},
	{[2]float64{0.0025344136713206, 0.0020005720388606}, TextNot},
	{[2]float64{0.0029379048771785, 0.0021993391130875}, TextNot},
	{[2]float64{0.0028319165180598, 0.0022608703948642}, TextNot},
	{[2]float64{0.0028865071627312, 0.0022612093554764}, TextNot},
	{[2]float64{0.0035105711405213, 0.0021281073199723}, TextNot},
	{[2]float64{0.0027806458982686, 0.0022995061647644}, TextNot},
	{[2]float64{0.0025604733845274, 0.0020330196217686}, TextNot},
	{[2]float64{0.0026622015240734, 0.0021399736915127}, TextNot},
	{[2]float64{0.0027587114059970, 0.0022435004332408}, TextNot},
	{[2]float64{0.0026219374546689, 0.0020189780622877}, TextNot},
	{[2]float64{0.0027395088538209, 0.0021888219204224}, TextNot},
	{[2]float64{0.0027085489113351, 0.0021866499197581}, TextNot},
	{[2]float64{0.0027325167866255, 0.0022972738357680}, TextNot},
	{[2]float64{0.0024445933956482, 0.0019862412736359}, TextSlop},
	{[2]float64{0.0025149471576455, 0.0019738924928700}, TextSlop},
	{[2]float64{0.0015424269288494, 0.0014938341177138}, TextSlop},
	{[2]float64{0.0024220115114023, 0.0020384929864067}, TextSlop},
	{[2]float64{0.0017118832908543, 0.0016371546721888}, TextSlop},
	{[2]float64{0.0026346499019590, 0.0022087475959419}, TextSlop},
	{[2]float64{0.0023339019148041, 0.0019073261603557}, TextSlop},
	{[2]float64{0.0023131888268327, 0.0019837243343036}, TextSlop},
	{[2]float64{0.0020004694725476, 0.0017521155106957}, TextSlop},
	{[2]float64{0.0019824144576281, 0.0015404338298265}, TextSlop},
	{[2]float64{0.0024191015400646, 0.0019010771130178}, TextSlop},
	{[2]float64{0.0024318700824397, 0.0019666544725462}, TextSlop},
	{[2]float64{0.0024489039635493, 0.0019906066490950}, TextSlop},
	{[2]float64{0.0025487131930542, 0.0020721469524445}, TextSlop},
	{[2]float64{0.0026200062738841, 0.0022413252379082}, TextSlop},
	{[2]float64{0.0021308244593731, 0.0019219290054469}, TextSlop},
	{[2]float64{0.0025408850862659, 0.0021009093606919}, TextSlop},
	{[2]float64{0.0027722206828839, 0.0023002548415115}, TextSlop},
	{[2]float64{0.0023466259895508, 0.0019546898084273}, TextSlop},
	{[2]float64{0.0026262351013152, 0.0022462855330569}, TextSlop},
	{[2]float64{0.0023265712402190, 0.0020203660903791}, TextSlop},
	{[2]float64{0.0025411085842163, 0.0021235251059460}, TextSlop},
	{[2]float64{0.0020511025331169, 0.0019632193392935}, TextSlop},
	{[2]float64{0.0025550910644093, 0.0020637651495933}, TextSlop},
	{[2]float64{0.0023367264173967, 0.0020243454128757}, TextSlop},
	{[2]float64{0.0024414392046835, 0.0019926771600527}, TextSlop},
	{[2]float64{0.0022505472938499, 0.0019734835325917}, TextSlop},
	{[2]float64{0.0025623734728518, 0.0021066516230483}, TextSlop},
	{[2]float64{0.0015494579987679, 0.0014992944010676}, TextSlop},
	{[2]float64{0.0024054167664316, 0.0020184889383590}, TextSlop},
	{[2]float64{0.0024874132527845, 0.0020640742268208}, TextSlop},
	{[2]float64{0.0025699314925276, 0.0021604241771556}, TextSlop},
	{[2]float64{0.0024047174992061, 0.0020722258672211}, TextSlop},
	{[2]float64{0.0022563102793485, 0.0019363341719843}, TextSlop},
	{[2]float64{0.0022083485785579, 0.0018674595302438}, TextSlop},
	{[2]float64{0.0021898323460208, 0.0019794324246494}, TextSlop},
	{[2]float64{0.0024844674078830, 0.0019770331421940}, TextSlop},
	{[2]float64{0.0024831399177127, 0.0020159619873255}, TextSlop},
	{[2]float64{0.0021142262108663, 0.0018861659418638}, TextSlop},
	{[2]float64{0.0019836141218865, 0.0018662628702630}, TextSlop},
	{[2]float64{0.0024520837641847, 0.0019947281153504}, TextSlop},
	{[2]float64{0.0022166256740873, 0.0018012521997888}, TextSlop},
	{[2]float64{0.0023170101015378, 0.0019439993565243}, TextSlop},
	{[2]float64{0.0026414136916751, 0.0022357455493961}, TextSlop},
	{[2]float64{0.0024531369435419, 0.0020623727948055}, TextSlop},
	{[2]float64{0.0022909557569519, 0.0018446530155318}, TextSlop},
	{[2]float64{0.0024855228483818, 0.0020537383030561}, TextSlop},
	{[2]float64{0.0021475658591521, 0.0017240581004856}, TextSlop},
	{[2]float64{0.0021671147561427, 0.0019353283914463}, TextSlop},
	{[2]float64{0.0020329115799536, 0.0015091663193545}, TextSlop},
	{[2]float64{0.0025201338997538, 0.0020970939652849}, TextSlop},
	{[2]float64{0.0026884490230808, 0.0022613757015215}, TextSlop},
	{[2]float64{0.0024110046547876, 0.0020069706136141}, TextSlop},
	{[2]float64{0.0024822767049441, 0.0020699164699811}, TextSlop},
	{[2]float64{0.0029593567459492, 0.0024116829952942}, TextSlop},
	{[2]float64{0.0022355475692285, 0.0019418266018627}, TextSlop},
	{[2]float64{0.0022753110504954, 0.0019787856896811}, TextSlop},
	{[2]float64{0.0021692623494004, 0.0017935157745070}, TextSlop},
	{[2]float64{0.0018238113597638, 0.0017921369989173}, TextSlop},
	{[2]float64{0.0021154275181101, 0.0017568265734081}, TextSlop},
	{[2]float64{0.0017224302827641, 0.0015215778422430}, TextSlop},
	{[2]float64{0.0024141895208603, 0.0020420342200405}, TextSlop},
	{[2]float64{0.0018532019733349, 0.0018377499719710}, TextSlop},
	{[2]float64{0.0025454933226278, 0.0020479447646390}, TextSlop},*/
	{[4]float64{0.0027671081299420, 0.0018650707586764, 0.0013793827281971, 0.0023159896423369}, TextNot},
	{[4]float64{0.0024715820022209, 0.0016703600086517, 0.0012103102055569, 0.0020669995199859}, TextNot},
	{[4]float64{0.0026193556997382, 0.0016321386256763, 0.0012813045173885, 0.0021114349639895}, TextNot},
	{[4]float64{0.0027188210448006, 0.0016112675112089, 0.0011980335560656, 0.0019507123641568}, TextNot},
	{[4]float64{0.0025840830287612, 0.0015360843762267, 0.0012141888928991, 0.0020338154475646}, TextNot},
	{[4]float64{0.0032537747798464, 0.0017799224455074, 0.0013297163389536, 0.0022458432297109}, TextNot},
	{[4]float64{0.0030208543655926, 0.0018681511748029, 0.0014125752332171, 0.0023784648090018}, TextNot},
	{[4]float64{0.0030199160180444, 0.0019334067281301, 0.0013979796324919, 0.0021907293955013}, TextNot},
	{[4]float64{0.0026994349599398, 0.0016057750943194, 0.0012341642690447, 0.0019187465888405}, TextNot},
	{[4]float64{0.0030260218133745, 0.0018656591442867, 0.0014291025680280, 0.0022521789657176}, TextNot},
	{[4]float64{0.0042592065200255, 0.0021477047372315, 0.0016568505432283, 0.0029556804780091}, TextNot},
	{[4]float64{0.0026546765932896, 0.0015590497598144, 0.0011987994050645, 0.0020798673656785}, TextNot},
	{[4]float64{0.0027594223439506, 0.0014226481847196, 0.0011452897846793, 0.0019554585536931}, TextNot},
	{[4]float64{0.0025918806265325, 0.0017900044308236, 0.0013223727622605, 0.0020638744373162}, TextNot},
	{[4]float64{0.0024739898736685, 0.0017262705492864, 0.0012911789957285, 0.0021394450159473}, TextNot},
	{[4]float64{0.0029726773896805, 0.0020243356887072, 0.0015458828520490, 0.0027373154756086}, TextNot},
	{[4]float64{0.0026751595676999, 0.0014551548868728, 0.0010794368903752, 0.0018356166904609}, TextNot},
	{[4]float64{0.0026971659298036, 0.0017842104279445, 0.0012172172118282, 0.0021493645607874}, TextNot},
	{[4]float64{0.0025590231275894, 0.0017896190481375, 0.0013235052936567, 0.0021476936862720}, TextNot},
	{[4]float64{0.0025289134565670, 0.0017751251214074, 0.0013457025763817, 0.0021735279176970}, TextNot},
	{[4]float64{0.0024870842928849, 0.0018971211979799, 0.0013961886841250, 0.0022894699219625}, TextNot},
	{[4]float64{0.0028243499031931, 0.0018792346258460, 0.0014061418941597, 0.0023959017614545}, TextNot},
	{[4]float64{0.0029075499031191, 0.0017067022510393, 0.0013012494919279, 0.0020190872751754}, TextNot},
	{[4]float64{0.0025487047367912, 0.0015490486062910, 0.0012245135153085, 0.0020291899928694}, TextNot},
	{[4]float64{0.0023368627091383, 0.0013681380005472, 0.0010704903983396, 0.0017967574795253}, TextNot},
	{[4]float64{0.0033354077537305, 0.0019324535550147, 0.0014835319305378, 0.0024872608510278}, TextNot},
	{[4]float64{0.0025814640050012, 0.0015829768809638, 0.0012210945934466, 0.0019210489437091}, TextNot},
	{[4]float64{0.0027792506686523, 0.0016179358026104, 0.0012829792104358, 0.0021180168265349}, TextNot},
	{[4]float64{0.0022579607558760, 0.0021040864578644, 0.0015771881603995, 0.0024520965741404}, TextNot},
	{[4]float64{0.0022983617071924, 0.0014880350293102, 0.0011571709359766, 0.0018849867897466}, TextNot},
	{[4]float64{0.0021906492701910, 0.0014422009163109, 0.0010887245876769, 0.0017136452460130}, TextNot},
	{[4]float64{0.0023944620415836, 0.0014866468371461, 0.0010962547938595, 0.0019682193226720}, TextNot},
	{[4]float64{0.0024939065349125, 0.0016193450246939, 0.0012739571381013, 0.0021525411413082}, TextNot},
	{[4]float64{0.0026267206446545, 0.0018836503483245, 0.0013808541635061, 0.0022707490487142}, TextNot},
	{[4]float64{0.0029258282744318, 0.0018067623085900, 0.0013946897747693, 0.0022963222732492}, TextNot},
	{[4]float64{0.0025992219533246, 0.0015668949177283, 0.0011930663686259, 0.0020319664462862}, TextNot},
	{[4]float64{0.0028060770270868, 0.0017713149926344, 0.0013377761264926, 0.0022601035811464}, TextNot},
	{[4]float64{0.0022798844561389, 0.0013143588443162, 0.0010262563241891, 0.0017040073424912}, TextNot},
	{[4]float64{0.0022546838432058, 0.0017282817641441, 0.0012615211785127, 0.0021323583071830}, TextNot},
	{[4]float64{0.0025397396367300, 0.0015419160524272, 0.0012548900597187, 0.0021441577359923}, TextNot},
	{[4]float64{0.0030600140400914, 0.0025834606762385, 0.0019633810144307, 0.0032363109608477}, TextNot},
	{[4]float64{0.0027802379613480, 0.0017359995776188, 0.0013927247017437, 0.0024912842619776}, TextNot},
	{[4]float64{0.0030133758611835, 0.0019094061653227, 0.0014183357470283, 0.0023317105018946}, TextNot},
	{[4]float64{0.0024367762075795, 0.0017017106928510, 0.0012804827829991, 0.0020814788292197}, TextNot},
	{[4]float64{0.0026651291875348, 0.0017895782732075, 0.0013630992421401, 0.0022754321009824}, TextNot},
	{[4]float64{0.0030995077704995, 0.0017468124060504, 0.0013127856015900, 0.0022724064349816}, TextNot},
	{[4]float64{0.0026769202985166, 0.0018177669763381, 0.0013507330081268, 0.0022426397596886}, TextNot},
	{[4]float64{0.0023960381190513, 0.0014123210093946, 0.0010927826747631, 0.0017263080229883}, TextNot},
	{[4]float64{0.0028795887644575, 0.0018299247966788, 0.0014365647364997, 0.0022708453458075}, TextNot},
	{[4]float64{0.0029379048771785, 0.0016077957824390, 0.0012224373910801, 0.0020095851494316}, TextNot},
	{[4]float64{0.0025965666477240, 0.0016350595476538, 0.0013039337988181, 0.0021441235270373}, TextNot},
	{[4]float64{0.0030094657475025, 0.0019781815687612, 0.0015488110941487, 0.0025413048899338}, TextNot},
	{[4]float64{0.0028319165180598, 0.0019825050600167, 0.0015045604880051, 0.0025270078361879}, TextNot},
	{[4]float64{0.0027806458982686, 0.0018493496593454, 0.0014126805942297, 0.0022454873098372}, TextNot},
	{[4]float64{0.0035105711405213, 0.0017318933571370, 0.0013697835930100, 0.0024752261302374}, TextNot},
	{[4]float64{0.0026219374546689, 0.0015552739270680, 0.0011849255952772, 0.0019332611602403}, TextNot},
	{[4]float64{0.0028865071627312, 0.0019310947888454, 0.0013723204699741, 0.0022574286830043}, TextNot},
	{[4]float64{0.0025604733845274, 0.0015936288495069, 0.0012596690926269, 0.0021659987181540}, TextNot},
	{[4]float64{0.0025344136713206, 0.0016149897136105, 0.0012300947941434, 0.0019920846959129}, TextNot},
	{[4]float64{0.0026622015240734, 0.0014946753130285, 0.0011967298712465, 0.0019012997466023}, TextNot},
	{[4]float64{0.0027587114059970, 0.0016756235263486, 0.0013132322658963, 0.0021825822585113}, TextNot},
	{[4]float64{0.0027395088538209, 0.0014938964631065, 0.0011926046211550, 0.0019579938987141}, TextNot},
	{[4]float64{0.0027325167866255, 0.0018667952383545, 0.0014027109165932, 0.0021826164356068}, TextNot},
	{[4]float64{0.0027085489113351, 0.0016084067647778, 0.0012446771367358, 0.0019713691694717}, TextNot},
	{[4]float64{0.0026346499019590, 0.0023795203050075, 0.0017955075845974, 0.0032754571472987}, TextSlop},
	{[4]float64{0.0024445933956482, 0.0024092609374042, 0.0017141814942684, 0.0028886185566527}, TextSlop},
	{[4]float64{0.0024220115114023, 0.0023974000028139, 0.0016467332379604, 0.0027843534200460}, TextSlop},
	{[4]float64{0.0023339019148041, 0.0021628689339562, 0.0015738865415321, 0.0032362903964112}, TextSlop},
	{[4]float64{0.0017118832908543, 0.0018170029724951, 0.0021392519908878, 0.0018682503862603}, TextSlop},
	{[4]float64{0.0025149471576455, 0.0024485643391622, 0.0016796535998588, 0.0029205778421187}, TextSlop},
	{[4]float64{0.0024191015400646, 0.0019107666044944, 0.0014700996502386, 0.0031496299726917}, TextSlop},
	{[4]float64{0.0019824144576281, 0.0019312987918119, 0.0022716121951284, 0.0021784324977752}, TextSlop},
	{[4]float64{0.0015424269288494, 0.0017566465319890, 0.0020610959923539, 0.0019019256013487}, TextSlop},
	{[4]float64{0.0020004694725476, 0.0022330231652716, 0.0015627532578322, 0.0024925241788301}, TextSlop},
	{[4]float64{0.0026200062738841, 0.0024486334325742, 0.0017455360515521, 0.0032573507283052}, TextSlop},
	{[4]float64{0.0023131888268327, 0.0024583644049485, 0.0017092356392660, 0.0027858316225631}, TextSlop},
	{[4]float64{0.0024318700824397, 0.0023190636795571, 0.0017551250339100, 0.0031680268623566}, TextSlop},
	{[4]float64{0.0025487131930542, 0.0021932326471764, 0.0017267526709273, 0.0032052455347723}, TextSlop},
	{[4]float64{0.0024489039635493, 0.0022528549930097, 0.0023873360316658, 0.0024977878777055}, TextSlop},
	{[4]float64{0.0021308244593731, 0.0022015443525626, 0.0015919022445399, 0.0024546518863252}, TextSlop},
	{[4]float64{0.0025408850862659, 0.0023003682794513, 0.0017960063618873, 0.0031919003595331}, TextSlop},
	{[4]float64{0.0025411085842163, 0.0024953221883164, 0.0017578404437131, 0.0033047267310241}, TextSlop},
	{[4]float64{0.0023466259895508, 0.0022371868050534, 0.0016203116263492, 0.0027252849633706}, TextSlop},
	{[4]float64{0.0023265712402190, 0.0023283279250706, 0.0016951306656207, 0.0032374448701393}, TextSlop},
	{[4]float64{0.0026262351013152, 0.0023819765239135, 0.0018178791854970, 0.0035523100170029}, TextSlop},
	{[4]float64{0.0020511025331169, 0.0021992819108060, 0.0023532187901836, 0.0024602432547103}, TextSlop},
	{[4]float64{0.0027722206828839, 0.0024637633421720, 0.0018530187095403, 0.0036098153017385}, TextSlop},
	{[4]float64{0.0025550910644093, 0.0023640809796906, 0.0017372239718328, 0.0032861048181083}, TextSlop},
	{[4]float64{0.0022505472938499, 0.0023181125871234, 0.0019505166297320, 0.0028416916324516}, TextSlop},
	{[4]float64{0.0024414392046835, 0.0023140589383058, 0.0017103087621624, 0.0031396087748239}, TextSlop},
	{[4]float64{0.0023367264173967, 0.0027105985986072, 0.0018108764119995, 0.0029906823540580}, TextSlop},
	{[4]float64{0.0024054167664316, 0.0024034405981757, 0.0016037725279223, 0.0026770483764420}, TextSlop},
	{[4]float64{0.0022083485785579, 0.0019514925340856, 0.0022959806575223, 0.0022024093022974}, TextSlop},
	{[4]float64{0.0015494579987679, 0.0018132412510090, 0.0021451877310601, 0.0019622975640310}, TextSlop},
	{[4]float64{0.0025623734728518, 0.0022833753559109, 0.0016825561237019, 0.0031404440600489}, TextSlop},
	{[4]float64{0.0024874132527845, 0.0021797487809452, 0.0015569091373338, 0.0025904994599736}, TextSlop},
	{[4]float64{0.0025699314925276, 0.0023985005524013, 0.0017497498493262, 0.0027218776374848}, TextSlop},
	{[4]float64{0.0024047174992061, 0.0025697045390787, 0.0022295541325069, 0.0025791968609730}, TextSlop},
	{[4]float64{0.0021898323460208, 0.0020826736643211, 0.0015764918923868, 0.0026617146588443}, TextSlop},
	{[4]float64{0.0022563102793485, 0.0022596402820084, 0.0016437662664370, 0.0026909452662668}, TextSlop},
	{[4]float64{0.0024844674078830, 0.0020984952516059, 0.0015611331254355, 0.0031239709527028}, TextSlop},
	{[4]float64{0.0024831399177127, 0.0022343474652055, 0.0016905007805257, 0.0033228287668854}, TextSlop},
	{[4]float64{0.0022166256740873, 0.0020725813989429, 0.0017442168564542, 0.0025373148954389}, TextSlop},
	{[4]float64{0.0021142262108663, 0.0023130981358424, 0.0016437570568294, 0.0024523279240260}, TextSlop},
	{[4]float64{0.0023170101015378, 0.0023202720866677, 0.0020664684116205, 0.0024963387087092}, TextSlop},
	{[4]float64{0.0024520837641847, 0.0023388613184226, 0.0023167183239559, 0.0026209111789547}, TextSlop},
	{[4]float64{0.0026414136916751, 0.0023926665178070, 0.0018219088644347, 0.0034094353438639}, TextSlop},
	{[4]float64{0.0021671147561427, 0.0024313214849180, 0.0021717390115399, 0.0026210297140863}, TextSlop},
	{[4]float64{0.0019836141218865, 0.0021233867582165, 0.0023462720281855, 0.0021398558715555}, TextSlop},
	{[4]float64{0.0024531369435419, 0.0021116419212823, 0.0022520171660487, 0.0024107586827703}, TextSlop},
	{[4]float64{0.0022909557569519, 0.0023558807081804, 0.0020516241896899, 0.0024266154973533}, TextSlop},
	{[4]float64{0.0020329115799536, 0.0017406398892910, 0.0015172138627792, 0.0023045412831901}, TextSlop},
	{[4]float64{0.0021475658591521, 0.0019862651934050, 0.0016385523818418, 0.0025795809144743}, TextSlop},
	{[4]float64{0.0025201338997538, 0.0022817872961519, 0.0018530788553058, 0.0028254761715467}, TextSlop},
	{[4]float64{0.0029593567459492, 0.0024774495127289, 0.0018610242789919, 0.0035196998338891}, TextSlop},
	{[4]float64{0.0024855228483818, 0.0024656745634715, 0.0016695970492317, 0.0028240817394148}, TextSlop},
	{[4]float64{0.0024110046547876, 0.0025668941005722, 0.0022755690049728, 0.0028910763978388}, TextSlop},
	{[4]float64{0.0026884490230808, 0.0025318076249983, 0.0018107055026872, 0.0033817823100969}, TextSlop},
	{[4]float64{0.0022753110504954, 0.0025513932673083, 0.0017587809071002, 0.0028039251670335}, TextSlop},
	{[4]float64{0.0022355475692285, 0.0023475875850619, 0.0016002292408293, 0.0025819426636744}, TextSlop},
	{[4]float64{0.0021154275181101, 0.0022111192850091, 0.0015108472486788, 0.0024916963319585}, TextSlop},
	{[4]float64{0.0024822767049441, 0.0024337161774482, 0.0018998172350586, 0.0032970619373794}, TextSlop},
	{[4]float64{0.0018238113597638, 0.0018779118981321, 0.0022770294771607, 0.0020808902742223}, TextSlop},
	{[4]float64{0.0024141895208603, 0.0022195975195853, 0.0015050583242763, 0.0029928010483749}, TextSlop},
	{[4]float64{0.0021692623494004, 0.0022507270742492, 0.0016063860131893, 0.0025530925109486}, TextSlop},
	{[4]float64{0.0017224302827641, 0.0022577847542304, 0.0019588983410035, 0.0021738749407778}, TextSlop},
	{[4]float64{0.0018532019733349, 0.0021615454486076, 0.0015400123348836, 0.0023237067582403}, TextSlop},
	{[4]float64{0.0025454933226278, 0.0022253883055910, 0.0016909115234523, 0.0032573932284007}, TextSlop},
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

// Copy copies a node
func (n *Node) Copy() Node {
	cp := Node{
		Links: make(map[string]uint64),
		Keys:  make([]string, len(n.Keys)),
	}
	copy(cp.Keys, n.Keys)
	for key, value := range n.Links {
		cp.Links[key] = value
	}
	return cp
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

// Copy produces a copy of a graph
func (g *Graph) Copy() Graph {
	cp := Graph{
		Keys:  make([]string, len(g.Keys)),
		Graph: make(map[string]Node),
		Ranks: make(map[string]uint64),
	}
	copy(cp.Keys, g.Keys)
	for key, value := range g.Graph {
		cp.Graph[key] = value.Copy()
	}
	return cp
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
func (g *Graph) LearnFast(delta float64, iterations int, rng *rand.Rand, words, list []string) float64 {
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
			current /= float64(len(list))
			if math.Abs(current-previous) < delta {
				return count
			}
			previous = current
		}
	}
	return float64(iterations)
}

// LearnFastList adds context to a model
func (g *Graph) LearnFastList(delta float64, iterations int, rng *rand.Rand, words, list []string) float64 {
	for i, word := range list[:len(list)-1] {
		{
			node := g.Graph[word]
			if node.Links == nil {
				g.Keys = append(g.Keys, word)
				node.Links = make(map[string]uint64)
				node.Keys = make([]string, 0, 8)
			}
			count, ok := node.Links[list[i+1]]
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
			current /= float64(len(list))
			if math.Abs(current-previous) < delta {
				return count
			}
			previous = current
		}
	}
	return float64(iterations)
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
		Words []string
		Value uint64
		Cost  float64
	}
	set := make([]Trace, 0, 8)
	for range 8 {
		word := words[len(words)-1]
		node := g.Graph[word]
		trace := Trace{}
		for _, word := range words[:len(words)-1] {
			trace.Trace = trace.Trace + word + " "
			trace.Words = append(trace.Words, word)
			trace.Value += g.Ranks[word]
		}
		for range 33 {
			trace.Trace = trace.Trace + word + " "
			trace.Words = append(trace.Words, word)
			trace.Value += g.Ranks[word]
			sum := uint64(0)
			for _, w := range node.Keys {
				sum += g.Ranks[w]
			}
			for sum == 0 {
				node = g.Graph[words[rng.Intn(len(words))]]
				for _, w := range node.Keys {
					sum += g.Ranks[w]
				}
			}
			total, selected := uint64(0), uint64(rng.Intn(int(sum)))
			for _, w := range node.Keys {
				total += g.Ranks[w]
				if selected < total {
					word = w
					node = g.Graph[word]
					break
				}
			}
		}
		set = append(set, trace)
	}
	mark := make(map[string]int)
	var search func(depth int, word string)
	search = func(depth int, word string) {
		if depth > 4 {
			return
		}
		node := g.Graph[word]
		for key := range node.Links {
			if value, found := mark[key]; !found {
				mark[key] = depth
				search(depth+1, key)
			} else if value > depth {
				mark[key] = depth
				search(depth+1, key)
			}
		}
	}
	search(0, words[0])
	fmt.Println(mark)
	for i, trace := range set {
		cp := make([]string, len(words))
		copy(cp, words)
		cp = append(cp, trace.Words...)
		has, list := make(map[string]bool), make([]string, 0, 8)
		for _, word := range trace.Words {
			if !has[word] {
				list = append(list, word)
			}
		}
		gcp := g.Copy()
		count := gcp.LearnFastList(1e-6, 8*1024*1024, rng, cp, list)
		for _, word := range words {
			set[i].Cost += float64(gcp.Ranks[word]) / float64(count)
		}
	}

	sort.Slice(set, func(i, j int) bool {
		return set[i].Cost < set[j].Cost
	})
	for _, trace := range set {
		fmt.Println(trace.Cost, trace.Trace)
	}
}

// PreMode pre-generate model
func PreMode(text string) {
	rng := rand.New(rand.NewSource(1))
	words := strings.Fields(text)
	g := NewGraph()
	g.LearnFast(1e-6, 8*1024*1024, rng, words, g.Keys)
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

	done := make(chan [BookCount]float64, 8)
	cal := func(t string, seed int64, alt string) {
		var results [BookCount]float64
		for i, book := range BookSet {
			text := string(books[book].Text)
			words := strings.Fields(text)
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
			{
				words := append(cp, suffix...)
				g := NewGraph()
				count := g.LearnFast(1e-6, 8*1024*1024, rng, words, list)
				sum := 0.0
				for _, value := range list {
					sum += float64(g.Ranks[value]) / float64(count)
				}
				results[i] = float64(sum) / float64(len(list))
			}
		}
		fmt.Printf("{[%d]float64{", BookCount)
		for _, result := range results {
			fmt.Printf("%.16f,", result)
		}
		fmt.Printf("}, %s},\n", t)
		done <- results
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
	results := make([][BookCount]float64, 0, Samples)
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
		results := make([][BookCount]float64, 0, Samples)
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

// Class is a model for a class
type Class struct {
	Graph
	Total float64
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

// TestMode test
func TestMode() {
	rng := rand.New(rand.NewSource(1))

	context := gradient.Context[float64]{}
	set := context.NewSet()
	set.Add("w0", BookCount, 8)
	set.AddBias("b0", 8)
	set.Add("w1", 16, 2)
	set.AddBias("b1", 2)
	set.AddData("inputs", BookCount, len(Ranks))
	set.AddData("in", BookCount)
	set.AddData("outputs", 2, len(Ranks))
	set.InitAdam(rng)

	inputs := set.ByName["inputs"]
	outputs := set.ByName["outputs"]
	for i, rank := range Ranks {
		for ii := range rank.Rank {
			inputs.X[BookCount*i+ii] = rank.Rank[ii]
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
	var list [BookCount][]string
	var count [BookCount]float64
	var g [BookCount]Graph
	var result [BookCount]float64
	for i, book := range BookSet {
		text := string(books[book].Text)
		words := strings.Fields(text)
		suffix := strings.Fields(samples[*FlagTest] /*[:1024]*/)
		cp := make([]string, len(words))
		copy(cp, words)
		list[i] = make([]string, 0, 8)
		has := make(map[string]bool)
		for _, word := range suffix {
			if !has[word] {
				has[word] = true
				list[i] = append(list[i], word)
			}
		}
		{
			words := append(cp, suffix...)
			g[i] = NewGraph()
			count[i] = g[i].LearnFast(1e-5, 8*1024*1024, rng, words, list[i])
		}
		sum := 0.0
		for _, value := range list[i] {
			sum += float64(g[i].Ranks[value]) / float64(count[i])
		}
		result[i] = float64(sum) / float64(len(list[i]))
	}

	classes := make(Classes, BookCount)
	for i := range classes {
		classes[i].Graph = g[i]
		classes[i].Total = count[i]
	}

	max, index := -math.MaxFloat64, 0
	for i, result := range result {
		fmt.Printf("result%d=%.16f\n", i, result)
		type Result struct {
			Rank
			Diff float64
		}
		results := make([]Result, 0, len(Ranks))
		for _, rank := range Ranks {
			diff := math.Abs(rank.Rank[i] - result)
			results = append(results, Result{
				Rank: rank,
				Diff: diff,
			})
		}
		sort.Slice(results, func(i, j int) bool {
			return results[i].Diff < results[j].Diff
		})
		var histogram [2]int
		for i := range results[:16] {
			histogram[results[i].Type]++
		}
		fmt.Println(histogram)

		score := classes.Score(i, list[i])
		fmt.Println("score=", score)
		if score > max {
			max, index = score, i
		}
	}
	if index == 0 {
		fmt.Println("not")
	} else {
		fmt.Println("slop")
	}

	{
		inputs := set.ByName["in"]
		for i, result := range result {
			inputs.X[i] = result
		}

		l0 := Everett(Add(Mul(set.Get("w0"), set.Get("in")), set.Get("b0")))
		l1 := Add(Mul(set.Get("w1"), l0), set.Get("b1"))
		l1(func(a *gradient.V[float64]) bool {
			fmt.Println("l1", a.X[0], a.X[1])
			return true
		})
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
