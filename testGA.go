package main

import (
	"math/rand"
)

type testGA struct {
	seedForGA         int64
	solution          []byte
	alphabet          []byte
	mutationRate      float64
	exchangeInPercent float32
	populationSize    int

	//only for generated test cases
	initSeed         int
	randomSourceSeed int64
	random           *rand.Rand
}

/*Generates a whole single GA setup providing one seed only
*The same seed always results in the same GA setup -> deterministic
 */
func (t *testGA) generateTestBySeedSingleGA(initSeed int64) {
}

func (t *testGA) runTestSetupTimes(initSeed int64) {
	t.randomSourceSeed = initSeed
	var seedSource = rand.NewSource(t.randomSourceSeed)
	t.random = rand.New(seedSource)

	t.seedForGA = 22238380
	t.solution = []byte("dragons are awesome")
	t.alphabet = []byte("abcdefghijklmnopqrstuvwxyz ")
	t.mutationRate = 0.2
	t.exchangeInPercent = 0.1
	t.populationSize = 10000
}

func (t *testGA) logTest() {

}
