package main

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func RandomCloseRangeCryptoLevel(lowerLimit int, upperLimit int) int {

	rangeSize := upperLimit - lowerLimit + 1 // +1 as generator upperLimit is open
	nBig, _ := crand.Int(crand.Reader, big.NewInt(int64(rangeSize)))
	randomInteger := int(nBig.Int64())
	return lowerLimit + randomInteger
}

// Only to be initialice once, if we initialize it every time we call function, randomnes is not well distributed
var source rand.Source = rand.NewSource(time.Now().UnixNano())

func RandomCloseRange(lowerLimit int, upperLimit int) int {

	rangeSize := upperLimit - lowerLimit + 1 // +1 as generator upperLimit is open

	randomizer := rand.New(source)

	return lowerLimit + randomizer.Intn(rangeSize)
}

func Dice(diceSides int) int {
	return RandomCloseRange(1, 20)
}

func Advantage(dicesSides int) int {
	firstDice := Dice(dicesSides)
	secondDice := Dice(dicesSides)
	return int(math.Max(float64(firstDice), float64(secondDice)))
}

func Disadvantage(dicesSides int) int {
	firstDice := Dice(dicesSides)
	secondDice := Dice(dicesSides)
	return int(math.Min(float64(firstDice), float64(secondDice)))
}

func generateSamplesForBarchart(rollGenerator func(int) int, diceSides int) ([]int, []opts.BarData) {
	rollSamples := make([]opts.BarData, 0, 40)
	rollSamplesCount := make(map[int]int)
	const numberOfSamples int = 1000000
	for i := 0; i < numberOfSamples; i++ {
		rollSample := rollGenerator(diceSides)
		count, exists := rollSamplesCount[rollSample]
		if exists {
			rollSamplesCount[rollSample] = count + 1
		} else {
			rollSamplesCount[rollSample] = 1
		}
	}

	rolls := make([]int, 0, len(rollSamplesCount))
	for roll := range rollSamplesCount {
		rolls = append(rolls, roll) // strconv.Itoa(roll)
	}

	// sort the slice by keys
	sort.Ints(rolls)

	// iterate by sorted keys
	for _, roll := range rolls {
		// roll, _ := strconv.Atoi(roll)
		rollSamples = append(rollSamples, opts.BarData{Value: rollSamplesCount[roll]})
	}

	return rolls, rollSamples
}
func GenerateBarchart(categories1 []int, series1 []opts.BarData, categories2 []int, series2 []opts.BarData, categories3 []int, series3 []opts.BarData) {

	// create a new bar instance
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title:    "Roll comparison",
		Subtitle: "Disadvantage, Single Dice, Advantage",
	}))

	// Put data into instance
	bar.SetXAxis(categories1).
		AddSeries("Disadvantage", series1).
		SetXAxis(categories2).
		AddSeries("Single Dice", series2).
		SetXAxis(categories3).
		AddSeries("Advantage", series3)

	// Where the magic happens
	f, _ := os.Create("bar.html")
	bar.Render(f)
}

func main() {
	categories1, series1 := generateSamplesForBarchart(Disadvantage, 20)
	categories2, series2 := generateSamplesForBarchart(Dice, 20)
	categories3, series3 := generateSamplesForBarchart(Advantage, 20)
	GenerateBarchart(categories1, series1, categories2, series2, categories3, series3)
}
