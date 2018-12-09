package FindwordGA

import (
	"math/rand"
	"sync"
	"time"
)

/////////////////////////////////////////////////////////////////////
//Implementation of word finding genetic algorithm///////////////////
/////////////////////////////////////////////////////////////////////

//TODO Store all or some populations in a list/map/array ? watch memory!
//TODO Produce some kind of logfile
//TODO struct for all measured values

type FindwordGA struct {
	//parameters
	gaID             int
	mutationRate     float32
	solution         []byte
	alphabet         []byte
	populationSize   int
	randomSourceSeed int64
	numberOfThreads  int

	//init variables
	random           *rand.Rand
	population       [][]byte
	alphabetLength   int
	solutionLength   int
	doneChannel      chan bool
	doneChannelMutex *sync.Mutex
	randomForThreads []*rand.Rand
	popBlockStart    []int
	popBlockEnd      []int

	//runtime variables
	populationFitness  []int
	selectedParents    []int
	populationBuffer   [][]byte
	threadsDoneChannel chan int

	//statistics
	solutionIndex              int
	Duration                   time.Duration
	AverageFitnessInPercent    float32
	StandardDeviationInPercent float32
	evolutionCount             int

	//information for running GA's in a network
	exchangeGA exchangeGA

	//experimental variables
	selectedParentsDensityFunction []float32
}

//set ID of GA
func (m *FindwordGA) SetGaID(gaID int) {
	m.gaID = gaID
}

//setup the GA population with information about the population and solution
func (m *FindwordGA) InitPopulation(initSolution []byte, initAlphabet []byte, initPopulationSize int, initMutationRate float32,
	initSeed int64, initNumberOfThreads int) {

	if initMutationRate < 0 {
		panic("mutation rate must be greater than 0")
	} else if len(initAlphabet) < 1 {
		panic("alphabet must include chars")
	} else if initPopulationSize < 1 {
		panic("population size must be greater than 0")
	} else if initNumberOfThreads < 1 {
		panic("GA has to run with at least one thread")
	} else if len(initSolution) < 1 {
		panic("solution must include at least one char")
	}

	m.mutationRate = initMutationRate
	m.solution = initSolution
	m.alphabet = initAlphabet
	m.solutionLength = len(m.solution)
	m.alphabetLength = len(m.alphabet)
	m.populationSize = initPopulationSize
	m.population = make([][]byte, initPopulationSize)
	m.populationBuffer = make([][]byte, initPopulationSize)
	m.populationFitness = make([]int, initPopulationSize)
	m.selectedParents = make([]int, initPopulationSize)
	//m.selectedParentsDensityFunction = make([]float32, m.populationSize)
	m.numberOfThreads = initNumberOfThreads
	m.threadsDoneChannel = make(chan int, m.numberOfThreads)

	//if it is running alone
	m.doneChannel = make(chan bool, 1)
	m.doneChannelMutex = &sync.Mutex{}
	m.doneChannel <- false

	if initSeed != 0 {
		m.randomSourceSeed = initSeed
	} else {
		m.randomSourceSeed = time.Now().UnixNano()
	}
	m.random = rand.New(rand.NewSource(m.randomSourceSeed))

	m.randomForThreads = make([]*rand.Rand, m.numberOfThreads)
	for i := 0; i < m.numberOfThreads; i++ {
		m.randomForThreads[i] = rand.New(rand.NewSource(m.randomSourceSeed + int64(i)))
	}

	//filling the population with random strings
	for i := 0; i < m.populationSize; i++ {
		m.population[i] = m.randomString(m.solutionLength)
	}

	//calculate population blocks for threads
	m.popBlockStart = make([]int, m.numberOfThreads)
	m.popBlockEnd = make([]int, m.numberOfThreads)

	for i := 0; i < m.numberOfThreads; i++ {
		//when this is the last thread and there is a rest from dividing through numberOfThreads
		if i == m.numberOfThreads-1 {
			m.popBlockStart[i] = i * (m.populationSize / m.numberOfThreads)
			m.popBlockEnd[i] = (i+1)*int(m.populationSize/m.numberOfThreads) + (m.populationSize % m.numberOfThreads) - 1
		} else {
			m.popBlockStart[i] = i * (m.populationSize / m.numberOfThreads)
			m.popBlockEnd[i] = (i+1)*int(m.populationSize/m.numberOfThreads) - 1
		}
	}

}

//run GA with current setup must be initialized before
func (m *FindwordGA) Run() {

	//fmt.Printf("Genetic Algorithm Number %d starts\n", m.gaID)
	//fmt.Println("Search for string length: ", m.solutionLength)
	//fmt.Println("Search alphabet length: ", m.alphabetLength)
	//fmt.Println("with population size of: ", m.populationSize)
	//fmt.Println("start search ... ")

	startTime := time.Now()

	ex := m.exchangeGA

	for {
		//fmt.Println("Run ",m.evolutionCount)
		if ex.exchangeInterval != 0 && m.evolutionCount%ex.exchangeInterval == 0 {
			//fmt.Printf("GaID %d evolution will be exchanged %d\n",m.gaID,m.evolutionCount )
			m.exchangeIndividuals()
		}

		//m.solutionIndex = m.calcFitnessSingleThread()
		m.solutionIndex = m.calcFitnessMultiThread()

		//if solution was found OR other thread is done earlier
		//m.doneChannelMutex.Lock()
		if <-m.doneChannel == true || m.solutionIndex >= 0 {
			m.doneChannel <- true
			//m.doneChannelMutex.Unlock()
			break
		} else {
			m.doneChannel <- false
			//m.doneChannelMutex.Unlock()
		}

		//m.selectParentsSingleThread()
		m.selectParentsMultiThread()

		//m.crossoverPopulation()
		m.crossoverPopulationMultithread()

		//m.mutatePopulation()
		m.mutatePopulationMultiThread()

		m.evolutionCount++

	}

	m.Duration = time.Now().Sub(startTime)
	m.calcStatistics()

	//fmt.Println("time for search: ", m.Duration)
	//fmt.Printf("wordfindGA_ID %d: number of evolutions: %d\n", m.gaID, m.evolutionCount)
	//fmt.Println("Found solution at index: ", m.solutionIndex)
	//fmt.Print("solution searched : ",string(m.solution))
	//fmt.Print("solution found : ",string(m.population[m.solutionIndex]))
	//fmt.Printf("wordfindGA_ID %d: has %f average fitness last population\n", m.gaID, m.AverageFitnessInPercent)
	//fmt.Printf("wordfindGA_ID %d: has %f Standard Deviation in last population \n", m.gaID, m.StandardDeviationInPercent)

	if ex.waitgroup != nil {
		//if a connection has been set
		ex.waitgroup.Done()
	}

}

//generate a random string with specified number of chars
func (m *FindwordGA) randomString(numberOfChars int) []byte {

	randomString := make([]byte, numberOfChars)

	for i := 0; i < numberOfChars; i++ {
		randomString[i] = m.alphabet[m.random.Intn(m.alphabetLength)]
	}
	return randomString
}
