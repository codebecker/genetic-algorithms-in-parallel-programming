package main

import (
	"./FindwordGA"
	"./NetworkGA"
	"fmt"
)

func main() {

	var seed int64 = 222383033
	var solution = []byte("dragons are awesome")
	var alphabet = []byte("abcdefghijklmnopqrstuvwxyz ")
	var mutationRate float32 = 0.2
	var exchangeInPercent float32 = 0.1
	var populationSize = 100000

	fmt.Println("=====================TEST WITH A SINGLE GA IN SINGLETHREAD===============================")

	findWord1 := FindwordGA.FindwordGA{}
	findWord1.InitPopulation(solution, alphabet, populationSize, mutationRate, false, seed, 1)
	findWord1.Run()
	fmt.Printf("Duration for singlethread in GA : %s \n", findWord1.Duration.String())

	fmt.Println("=====================TEST WITH A SINGLE GA IN MULTITHREADING===============================")

	findWord2 := FindwordGA.FindwordGA{}
	findWord2.InitPopulation(solution, alphabet, populationSize, mutationRate, false, seed, 3)
	findWord2.Run()
	fmt.Printf("Duration for multithreading in GA : %s \n", findWord2.Duration.String())

	fmt.Println("=====================TEST WITH GA NETWORKS NO MULTITHREADING IN GA  ===============================")

	//four GA ring connection
	ring4GA := [][]int{
		{0, 1, 0, 1},
		{1, 0, 1, 0},
		{0, 1, 0, 1},
		{1, 0, 1, 0},
	}
	ring2GA := [][]int{
		{0, 1},
		{1, 0},
	}

	spiderSingle := new(NetworkGA.NetworkGA)
	spiderSingle.NetworkGA(solution, alphabet, populationSize, mutationRate, false, seed, exchangeInPercent, [][]int{{0}}, 0)

	spiderDual := new(NetworkGA.NetworkGA)
	spiderDual.NetworkGA(solution, alphabet, populationSize/2, mutationRate, false, seed, exchangeInPercent, ring2GA, 1)

	spiderQuadro := new(NetworkGA.NetworkGA)
	spiderQuadro.NetworkGA(solution, alphabet, populationSize/4, mutationRate, false, seed, exchangeInPercent, ring4GA, 1)

	fmt.Printf("==========================================COMPARISON BETWEEN NEWTORKS=============================\n")

	fmt.Printf("GA network single: %s \n", spiderSingle.Duration.String())
	fmt.Printf("GA network dual: %s \n", spiderDual.Duration.String())
	fmt.Printf("GA network quadro: %s \n", spiderQuadro.Duration.String())

}
