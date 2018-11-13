package NetworkGA

import (
	"../FindwordGA"
	"fmt"
	"sync"
	"time"
)

type NetworkGA struct {
	sendChannel     []chan [][]byte //size 1
	receiveChanel   []chan [][]byte //size 1
	numberOfThreads int
	Duration        time.Duration
}

func (s *NetworkGA) NetworkGA(solution []byte, alphabet []byte, populationSize int, mutationRate float64, logPop bool, seed int64,
	exchangePercent float32, sendToArray [][]int, exchangeInterval int) {

	var doneChannelMutex = sync.Mutex{}
	var sendToChannel = make([][]chan [][]byte, len(sendToArray))
	var rcvFromChannel = make([][]chan [][]byte, len(sendToArray))
	var gaConnectionArray = make([]FindwordGA.FindwordGA, len(sendToArray))

	var wg sync.WaitGroup
	wg.Add(len(sendToArray))

	doneChannel := make(chan bool, 1) //always size one
	doneChannel <- false

	for i := 0; i < len(gaConnectionArray); i++ {
		gaConnectionArray[i] = FindwordGA.FindwordGA{}
		gaConnectionArray[i].SetGaID(i)
		gaConnectionArray[i].InitPopulation(solution, alphabet, populationSize, mutationRate, logPop, seed+int64(i*64), 1)
	}

	//init network
	//init rcv and send channelArrays in second dimension size of zero
	//itterate over maparray and APPEND channels
	for i := 0; i < len(sendToArray); i++ {
		sendToChannel[i] = make([]chan [][]byte, 0)
		rcvFromChannel[i] = make([]chan [][]byte, 0)
	}

	//build send an rcv channel for each monkey dimension1 sendsTo dimension2
	for i := 0; i < len(sendToArray); i++ {
		for j := 0; j < len(sendToArray[i]); j++ {
			if sendToArray[i][j] == 1 {
				var newChan = make(chan [][]byte, 1)
				sendToChannel[i] = append(sendToChannel[i], newChan)
				rcvFromChannel[j] = append(rcvFromChannel[j], newChan)
			}
		}
	}

	for i := 0; i < len(sendToArray); i++ {
		gaConnectionArray[i].InitNetwork(exchangePercent, sendToChannel[i], rcvFromChannel[i], &wg, doneChannel, &doneChannelMutex, exchangeInterval)
	}
	fmt.Printf("==========================================GA NETWORK STARTS ======================================\n")

	fmt.Printf("Starting word find genetic algorithm network with %d threads population size of %d each thread alphabet of %d chars and a solution size of %d\n", len(sendToArray), populationSize, len(alphabet), len(solution))
	startTime := time.Now()
	//run all findwordGA in individual thread
	for i := 0; i < len(gaConnectionArray); i++ {
		go gaConnectionArray[i].Run()
	}
	wg.Wait()
	fmt.Printf("==========================================GA NETWORK ENDS ========================================\n")

	s.Duration = time.Now().Sub(startTime)
}
