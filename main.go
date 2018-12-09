package main

import (
	"fmt"
	"time"
)
import "./FindwordGA"
import "./NetworkGA"

func main() {
	/*
		var seed int64 = 222383033
		var solution = []byte("abcdefghijklmn")
		var alphabet = []byte("abcdefghijklmnopqrstuvwxyz ")
		var mutationRate float32 = 0.2
		var populationSize = 500
		var numberOfThreads = 4
		//benchmarkGA(solution, alphabet, populationSize, mutationRate, seed, numberOfThreads)
		//runBechmarkNetwork()
	*/
	runDifferentGAs()

}

func runDifferentGAs() {

	var seed int64 = 222383033
	var solution = []byte("monkey")
	var alphabet = []byte("abcdefghijklmnopqrstuvwxyz ")
	var mutationRate float32 = 0.2
	var exchangeInPercent float32 = 0.05
	var populationSize = 1000
	fmt.Printf("===================================START OF GA PRESENTATION====================================\n")
	fmt.Printf("============================USING SAME OPTIONS FOR ALL GA's====================================\n")

	fmt.Println("solution: ", string(solution))
	fmt.Println("alphabet: ", string(alphabet))
	fmt.Println("mutatuion rate: ", mutationRate)
	fmt.Println("exchange rate: ", exchangeInPercent)
	fmt.Println("population size: ", populationSize)

	fmt.Println("=====================TEST WITH A SINGLE GA IN SINGLETHREAD==============================")

	findWord1 := FindwordGA.FindwordGA{}
	findWord1.InitPopulation(solution, alphabet, populationSize, mutationRate, seed, 1)
	findWord1.Run()
	fmt.Printf("GA single thread: %s \n", findWord1.Duration.String())
	fmt.Println("=====================TEST WITH A SINGLE GA IN MULTITHREADING 2 THREADS==================")

	findWord2 := FindwordGA.FindwordGA{}
	findWord2.InitPopulation(solution, alphabet, populationSize, mutationRate, seed, 2)
	findWord2.Run()
	fmt.Printf("GA multi thread: %s \n", findWord2.Duration.String())

	fmt.Println("=====================TEST WITH GA SINGLE NETWORK NO MULTITHREADING IN GA================")

	ring2GA := [][]int{
		{0, 1},
		{1, 0},
	}

	networkGADual := new(NetworkGA.NetworkGA)
	networkGADual.NetworkGA(solution, alphabet, populationSize, mutationRate, seed, exchangeInPercent, ring2GA, 100)
	fmt.Printf("GA network dual: %s \n", networkGADual.Duration.String())

	fmt.Println("=====================TEST WITH GA QUADRO NETWORK NO MULTITHREADING IN GA================")
	//four GA ring connection
	ring4GA := [][]int{
		{0, 1, 0, 1},
		{1, 0, 1, 0},
		{0, 1, 0, 1},
		{1, 0, 1, 0},
	}

	networkGAQuadro := new(NetworkGA.NetworkGA)
	networkGAQuadro.NetworkGA(solution, alphabet, populationSize, mutationRate, seed, exchangeInPercent, ring4GA, 100)

	fmt.Printf("GA network quadro: %s \n", networkGAQuadro.Duration.String())

	fmt.Printf("==========================================END OF GA PRESENTATION=============================\n")

}

/*runs a GA configuration 10 times increases the solution size each run*/
func runBechmark() {
	var seed int64 = 222383033
	var solution = []byte("abc")
	var alphabet = []byte("abcdefghijklmnopqrstuvwxyz ")
	var mutationRate float32 = 0.2
	var populationSize = 10000
	var numberOfThreads = 4

	var solutionSteps = 2

	//length of array is the number of repeating
	timeForThreads := make([]int64, 10)

	//teste den GA erhöhe die länge von solution um 2 zeichen
	for i := 0; i < len(timeForThreads); i++ {
		//build new longer solution
		solution = make([]byte, solutionSteps*(i+1))
		for j := 0; j < solutionSteps*(i+1); j++ {
			solution[j] = alphabet[j]
		}
		//test 10 times
		timeForThreads[i] = benchmarkGA(solution, alphabet, populationSize, mutationRate, seed, numberOfThreads)
	}

	for i := 0; i < len(timeForThreads); i++ {
		println(timeForThreads[i])
	}
}

/*runs a GA configuration 10 times calcs the avg*/
func benchmarkGA(solution []byte, alphabet []byte, populationSize int, mutationRate float32, seed int64, numberOfThreds int) int64 {

	ga := FindwordGA.FindwordGA{}

	var duration time.Duration

	for i := 0; i < 10; i++ {
		ga = FindwordGA.FindwordGA{}
		ga.InitPopulation(solution, alphabet, populationSize, mutationRate, seed, numberOfThreds)
		ga.Run()
		fmt.Println(ga.Duration.String())
		duration += ga.Duration
	}

	return duration.Nanoseconds() / 10
}

func runBechmarkNetwork() {
	var solutionSteps = 2
	var seed int64 = 222383033
	var solution = []byte("abc")
	var alphabet = []byte("abcdefghijklmnopqrstuvwxyz ")
	var mutationRate float32 = 0.2
	var populationSize = 10000
	var exhangeRateInPercent = 0.2

	//four GA ring connection

	ring4GA := [][]int{
		{0, 1, 0, 1},
		{1, 0, 1, 0},
		{0, 1, 0, 1},
		{1, 0, 1, 0},
	}
	/*
		ring2GA := [][]int{
			{0, 1},
			{1, 0},
		}*/

	//length of array is the number of repeating
	timeForThreads := make([]int64, 10)

	//teste den GA erhöhe die länge von solution um 2 zeichen
	for i := 0; i < len(timeForThreads); i++ {
		//build new longer solution
		solution = make([]byte, solutionSteps*(i+1))
		for j := 0; j < solutionSteps*(i+1); j++ {
			solution[j] = alphabet[j]
		}
		//test 10 times
		timeForThreads[i] = benchmarkGANetwork(solution, alphabet, populationSize, mutationRate, seed, float32(exhangeRateInPercent), ring4GA, 1)
	}

	for i := 0; i < len(timeForThreads); i++ {
		println(timeForThreads[i])
	}
}

/*runs a GA NETWORK configuration 10 times calcs the avg*/
func benchmarkGANetwork(solution []byte, alphabet []byte, populationSize int, mutationRate float32, seed int64,
	exchangeInPercent float32, topology [][]int, exchangeInterval int) int64 {

	ga := new(NetworkGA.NetworkGA)

	var duration time.Duration

	for i := 0; i < 10; i++ {
		ga = new(NetworkGA.NetworkGA)
		ga.NetworkGA(solution, alphabet, populationSize, mutationRate, seed, exchangeInPercent, topology, exchangeInterval)
		duration += ga.Duration
	}
	return duration.Nanoseconds() / 10
}
