package NetworkGA

import (
	"../FindwordGA"
	"sync"
	"time"
)

type NetworkGA struct {
	sendChannel     []chan [][]byte //size 1
	receiveChanel   []chan [][]byte //size 1
	numberOfThreads int
	Duration        time.Duration
}

func (s *NetworkGA) NetworkGA(solution []byte, alphabet []byte, populationSize int, mutationRate float32, seed int64,
	exchangePercent float32, connections [][]int, exchangeInterval int) {

	if mutationRate < 0 {
		panic("mutation rate must be greater than 0")
	} else if len(alphabet) < 1 {
		panic("alphabet must include chars")
	} else if populationSize < 1 {
		panic("population size must be greater than 0")
	} else if len(solution) < 1 {
		panic("solution must include at least one char")
	} else if len(solution) < 1 {
		panic("solution must include at least one char")
	} else if exchangePercent < 0 {
		panic("exchange rate must be at least 0 percenet")
	} else if len(connections) < 1 {
		panic("connections must contain network connection information")
	} else if exchangeInterval < 0 {
		panic("exchange interval must be greater than 0")
	}

	for i := 0; i < len(connections); i++ {
		if len(connections) != len(connections[i]) {
			panic("connections Array must have same width and height in both dimensions")
		}
	}

	var doneChannelMutex = sync.Mutex{}
	var sendToChannel = make([][]chan [][]byte, len(connections))
	var rcvFromChannel = make([][]chan [][]byte, len(connections))
	var gaConnectionArray = make([]FindwordGA.FindwordGA, len(connections))

	var wg sync.WaitGroup
	wg.Add(len(connections))

	doneChannel := make(chan bool, 1) //always size of one so threads run async
	doneChannel <- false

	for i := 0; i < len(gaConnectionArray); i++ {
		gaConnectionArray[i] = FindwordGA.FindwordGA{}
		gaConnectionArray[i].SetGaID(i)
		gaConnectionArray[i].InitPopulation(solution, alphabet, populationSize, mutationRate, seed+int64(i*64), 1)
	}

	//init network
	//init rcv and send channelArrays in second dimension size of zero
	//itterate over maparray and APPEND channels
	for i := 0; i < len(connections); i++ {
		sendToChannel[i] = make([]chan [][]byte, 0)
		rcvFromChannel[i] = make([]chan [][]byte, 0)
	}

	//build send an rcv channel for each monkey dimension1 sendsTo dimension2
	for i := 0; i < len(connections); i++ {
		for j := 0; j < len(connections[i]); j++ {
			if connections[i][j] == 1 {
				var newChan = make(chan [][]byte, 1)
				sendToChannel[i] = append(sendToChannel[i], newChan)
				rcvFromChannel[j] = append(rcvFromChannel[j], newChan)
			}
		}
	}

	for i := 0; i < len(connections); i++ {
		gaConnectionArray[i].InitNetwork(exchangePercent, sendToChannel[i], rcvFromChannel[i], &wg, doneChannel, &doneChannelMutex, exchangeInterval)
	}

	startTime := time.Now()
	//run all findwordGA in individual threads
	for i := 0; i < len(gaConnectionArray); i++ {
		go gaConnectionArray[i].Run()
	}
	wg.Wait()
	//fmt.Printf("==========================================GA NETWORK ENDS ========================================\n")

	s.Duration = time.Now().Sub(startTime)
}
