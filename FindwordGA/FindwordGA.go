package FindwordGA

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

type FindwordGA struct {
	gaID                             int
	populationSize                   int
	mutationRate                     float64
	solution                         []byte
	alphabet                         []byte
	alphabetLength                   int
	log                              bool
	population                       [][]byte
	populationBuffer                 [][]byte
	populationHistory                [][][]byte
	selectedParents                  []int
	selectedParentsDensityFunction   []float32
	solutionIndex                    int
	solutionLength                   int
	populationFitness                []int
	evolutionCount                   int
	AverageFitnessInPercent          float32
	StandardDeviationInPercent       float32
	Duration                         time.Duration
	numberOfThreads                  int
	sendChannel                      []chan [][]byte
	receiveChannel                   []chan [][]byte
	doneChannelMutex                 *sync.Mutex
	doneChannel                      chan bool //size 1
	numberOfIndividualsToBeExchanged int
	exchangeInterval                 int
	numberOfSendRcvChannel           int
	randomSourceSeed                 int64
	random                           *rand.Rand
	waitgroup                        *sync.WaitGroup

	averageFitnessHist     []float32
	standardDerivationHist []float32
}

func (m *FindwordGA) SetGaID(gaID int) {
	m.gaID = gaID
}

func (m *FindwordGA) InitPopulation(initSolution []byte, initAlphabet []byte, initPopulationSize int, initMutationRate float64, logPop bool, initSeed int64, initNumberOfThreads int) {

	/////////////////////////////////////////////////////////////////////
	//Implementation of word finding genetic algorithm///////////////////
	/////////////////////////////////////////////////////////////////////

	//TODO Store all or some populations in a list/map/array ? watch memory!
	//TODO Produce some kind of a logfile
	//TODO struct for all measure values

	m.mutationRate = initMutationRate
	m.solution = initSolution
	m.alphabet = initAlphabet
	m.log = logPop
	m.solutionLength = len(m.solution)
	m.alphabetLength = len(m.alphabet)
	m.populationSize = initPopulationSize
	m.population = make([][]byte, initPopulationSize)
	m.populationBuffer = make([][]byte, initPopulationSize)
	m.populationFitness = make([]int, initPopulationSize)
	m.selectedParents = make([]int, initPopulationSize)
	m.randomSourceSeed = initSeed
	m.selectedParentsDensityFunction = make([]float32, m.populationSize)
	m.numberOfThreads = initNumberOfThreads

	//if it is running alone
	m.doneChannel = make(chan bool, 1)
	m.doneChannelMutex = &sync.Mutex{}
	m.doneChannel <- false

	if initSeed != 0 {
		m.randomSourceSeed = initSeed
	} else {
		m.randomSourceSeed = time.Now().UnixNano()
	}
	var seedSource = rand.NewSource(m.randomSourceSeed)
	m.random = rand.New(seedSource)

	//filling the population with random strings
	for i := 0; i < m.populationSize; i++ {
		m.population[i] = m.randomString(m.solutionLength)
	}
}

func (m *FindwordGA) InitNetwork(exchangePercent float32,
	initSendChannel []chan [][]byte,
	initReceiveChannels []chan [][]byte,
	wg *sync.WaitGroup,
	doneChannel chan bool,
	doneChannelMutex *sync.Mutex,
	exchangeInterval int) {

	m.numberOfIndividualsToBeExchanged = int(float32(m.populationSize) * exchangePercent)
	m.numberOfSendRcvChannel = len(initSendChannel)

	m.sendChannel = initSendChannel
	m.receiveChannel = initReceiveChannels
	m.waitgroup = wg
	m.doneChannel = doneChannel
	m.doneChannelMutex = doneChannelMutex
	m.exchangeInterval = exchangeInterval

}

func (m *FindwordGA) Run() {

	fmt.Printf("Genetic Algorithm Number %d starts\n", m.gaID)
	//fmt.Println("Search for string length: ", m.solutionLength)
	//fmt.Println("Search alphabet length: ", m.alphabetLength)
	//fmt.Println("with population size of: ", m.populationSize)
	//fmt.Println("start search ... ")

	startTime := time.Now()

	for {
		m.calcStatistics()
		if m.exchangeInterval != 0 && m.evolutionCount%m.exchangeInterval == 0 {
			//fmt.Printf("mokeyID %d evolution will be exchanged %d\n",m.gaID,m.evolutionCount )
			m.exchangeIndividuals()
		}

		m.solutionIndex = m.calcFitnessOfIndividuals()

		if m.log == true {
			m.averageFitnessHist = append(m.averageFitnessHist, m.AverageFitnessInPercent)
			m.standardDerivationHist = append(m.standardDerivationHist, m.StandardDeviationInPercent)
			/*
				for i := 0; i < len(m.populationHistory); i++ {
					//fmt.Println(string(m.populationHistory[i][m.solutionIndex]))
					m.averageFitnessHist = append(m.averageFitnessHist, m.AverageFitnessInPercent)

				}*/
		}

		m.doneChannelMutex.Lock()
		//solution was found OR other thread is done earlier

		//TODO DEADLOCK HERE? MULTIPLE THREADS EACH PRODUCING AND CONSUMING
		if true == <-m.doneChannel || m.solutionIndex >= 0 {
			m.doneChannel <- true
			m.doneChannelMutex.Unlock()
			break
		} else {
			m.doneChannel <- false
			m.doneChannelMutex.Unlock()
		}

		//m.selectParentsDensityFunction()
		//m.crossoverPopulationDensityFunction()
		m.selectParents()
		m.crossoverPopulation()

		m.mutatePopulation()
		m.evolutionCount = m.evolutionCount + 1

	}

	m.Duration = time.Now().Sub(startTime)

	//fmt.Println("time for search: ", m.Duration)
	fmt.Printf("wordfindGA_ID %d: number of evolutions: %d\n", m.gaID, m.evolutionCount)
	//fmt.Println("Found solution at index: ", m.solutionIndex)
	//fmt.Print("solution searched : ",string(m.solution))
	//fmt.Print("solution found : ",string(m.population[m.solutionIndex]))
	fmt.Printf("wordfindGA_ID %d: has %f average fitness last population\n", m.gaID, m.AverageFitnessInPercent)
	fmt.Printf("wordfindGA_ID %d: has %f Standard Deviation in last population \n", m.gaID, m.StandardDeviationInPercent)

	if m.waitgroup != nil {
		//if a connection has been set
		m.waitgroup.Done()
	}

}

func (m *FindwordGA) exchangeIndividuals() {
	//// Send and receive best individuals from channels a- / synchronously

	//TODO may make channels send channels with bigger size to send to multiple FindwordGA but then synchronus mode made cant be achieved to easy anymore
	//TODO exchange every X evolutions IMPORTANT

	//init array for individuals
	fitnessCurrentlySearched := m.solutionLength
	var individualsToBeSend = make([][]byte, m.numberOfIndividualsToBeExchanged)
	for i := 0; i < m.numberOfIndividualsToBeExchanged; i++ {
		individualsToBeSend[i] = make([]byte, m.solutionLength)
	}

	//collect best individuals to send
	for i := 0; i < m.numberOfIndividualsToBeExchanged && len(m.sendChannel) != 0; {
		for j := 0; j < len(m.populationFitness); j++ {
			if m.populationFitness[j] == fitnessCurrentlySearched {

				//individualsToBeSend [i] = []byte("  wordfindGA_ID send this")
				//individualsToBeSend[i][0] = byte(m.gaID)
				for k := 0; k < m.solutionLength; k++ {
					individualsToBeSend[i][k] = m.population[j][k]
				}
				i++
				if i >= m.numberOfIndividualsToBeExchanged {
					break
				}
			}
		}
		fitnessCurrentlySearched--
	}

	//send N best lokal individuals in channel
	for i := 0; i < len(m.sendChannel); i++ {
		//gehe in jeden channel leeren wenn inhalt sonst schreiben
		select {
		case <-m.sendChannel[i]:
			m.sendChannel[i] <- individualsToBeSend
		default:
			m.sendChannel[i] <- individualsToBeSend
		}
	}

	//get all available new individuals
	var receivedIndividuals = make([][]byte, 0)
	var receivedIndividualsCopyBuffer [][]byte
	for i := 0; i < len(m.receiveChannel); i++ {
		//test each channel empty it has content and write data
		select {
		case rcv := <-m.receiveChannel[i]:
			//append fucks stuff up so init manually

			//fmt.Printf("wordfindGA_ID %d RCV from %d \n", m.gaID, rcv[0][0])

			receivedIndividualsCopyBuffer = make([][]byte, len(receivedIndividuals)+len(rcv))

			for i := 0; i < len(receivedIndividualsCopyBuffer); i++ {
				receivedIndividualsCopyBuffer[i] = make([]byte, m.solutionLength)
			}

			//copy old individuals
			for i := 0; i < len(receivedIndividuals); i++ {
				for j := 0; j < m.solutionLength; j++ {
					receivedIndividualsCopyBuffer[i][j] = receivedIndividuals[i][j]
				}
			}
			//past new individuals with offset
			for i := 0; i < len(rcv); i++ {
				for j := 0; j < m.solutionLength; j++ {
					receivedIndividualsCopyBuffer[len(receivedIndividuals)+i][j] = rcv[i][j]
				}
			}
			receivedIndividuals = receivedIndividualsCopyBuffer

		default:
		}
	}

	//fmt.Printf("wordfindGA_ID %d received number %d of individuals \n", m.gaID, len(receivedIndividuals))

	//replace worst lokal individuals with new individuals
	fitnessCurrentlySearched = 0

	for i := 0; i < len(receivedIndividuals); {

		for j := 0; j < len(m.populationFitness); j++ {

			if m.populationFitness[j] == fitnessCurrentlySearched {
				//deep copy

				for k := 0; k < m.solutionLength; k++ {
					receivedIndividuals[i][k] = m.population[j][k]
					m.population[j][k] = receivedIndividuals[i][k]
				}
				i++
				if i >= len(receivedIndividuals) {
					break
				}
			}
		}
		fitnessCurrentlySearched++
	}
}

func (m *FindwordGA) randomString(numberOfChars int) []byte {

	randomString := make([]byte, numberOfChars)

	for i := 0; i < numberOfChars; i++ {
		randomString[i] = m.alphabet[m.random.Intn(m.alphabetLength)]
	}
	return randomString
}

func (m *FindwordGA) calcFitnessOfIndividuals() int {
	solutionChan := make(chan int, m.numberOfThreads)
	solution := -1

	for i := 0; i < m.numberOfThreads; i++ {
		//when this is the last thread and there is a rest from dividing throug numberOfThreads
		if i == m.numberOfThreads-1 && m.populationSize%m.numberOfThreads != 0 {
			go m.calcFitnessOfIndividualsParallel(solutionChan, i*(m.populationSize/m.numberOfThreads), (i+1)*int(m.populationSize/m.numberOfThreads)+(m.populationSize%m.numberOfThreads))
		} else {
			go m.calcFitnessOfIndividualsParallel(solutionChan, i*(m.populationSize/m.numberOfThreads), (i+1)*m.populationSize/m.numberOfThreads)
		}
	}

	for i := 0; i < m.numberOfThreads; i++ {
		tmp := <-solutionChan
		//if a solution was found and no other solution has been found before
		if tmp > -1 && solution == -1 {
			solution = tmp
		}
	}
	return solution
}

func (m *FindwordGA) calcFitnessOfIndividualsParallel(solutionChan chan int, startIndividual int, endIndividual int) int {
	var fittest = -1
	for i := startIndividual; i < endIndividual; i++ {
		fitnessOfIndividual := 0
		for j := 0; j < m.solutionLength; j++ {
			if m.population[i][j] == m.solution[j] {
				fitnessOfIndividual++
			}
		}
		m.populationFitness[i] = fitnessOfIndividual
		//if a solution was found AND ther has been no solution before
		if fitnessOfIndividual == m.solutionLength && fittest < 0 {
			fittest = i
		}
	}
	//return solution if found else return -1
	solutionChan <- fittest
	return fittest
}

func (m *FindwordGA) selectParents() {
	//the fiter the individuals the more often their index will appear in the selectedParents array
	numberFitnessTotal := 0
	for i := 0; i < m.populationSize; i++ {
		numberFitnessTotal = numberFitnessTotal + m.populationFitness[i]
	}
	//check if fitness number is
	if numberFitnessTotal == 0 {
		m.selectedParents = make([]int, m.populationSize)
		for i := 0; i < len(m.selectedParents); i++ {
			m.selectedParents[i] = 1
		}
	} else {
		m.selectedParents = make([]int, numberFitnessTotal)
		fitnessCounter := 0
		for i := 0; i < m.populationSize; i++ {
			for j := 0; j < m.populationFitness[i]; j++ {
				m.selectedParents[fitnessCounter] = i
				fitnessCounter++
			}
		}
	}
}

func (m *FindwordGA) crossoverPopulation() {
	var areEqual = true
	m.populationBuffer = make([][]byte, m.populationSize)
	individual1 := make([]byte, m.solutionLength)
	individual2 := make([]byte, m.solutionLength)
	buffer := make([]byte, m.solutionLength)

	for i := 0; i < m.populationSize; i++ {
		m.populationBuffer[i] = make([]byte, m.solutionLength)
		areEqual = true
		buffer = make([]byte, m.solutionLength)
		individual1 = m.population[m.selectedParents[m.random.Intn(len(m.selectedParents))]]
		individual2 = m.population[m.selectedParents[m.random.Intn(len(m.selectedParents))]]

		for j := 0; j < m.solutionLength; j++ {
			buffer[j] = individual1[j]
			if buffer[j] != individual2[j] {
				areEqual = false
			}
		}
		if areEqual == true {
			i--
		} else {
			for j := 0; j < len(individual1)/2; j++ {
				buffer[j] = individual2[j]
			}
			m.populationBuffer[i] = buffer
		}
	}

	m.population = m.populationBuffer
}

// density functions may have better performance with huge populations.
//because with array density the parents array can reach size of lenghtOf(int) * population * (1+avgFitness) with avgFitness [0,1.0]
//whereas the parentsDensityFuntion Array always has the same size of population * sizeOf(float32)
func (m *FindwordGA) selectParentsDensityFunction() {
	//erstellen eines stufigen wahrscheinlichkeitsdiagrammes w(i)= i + SUM(0,1)
	//kandaten mit einer wahrscheinlichkeit von 0 ergeben keine abstufung und der vorangegangene kanidat wird gewählt
	var numberFitnessTotal float32 = 0
	for i := 0; i < m.populationSize; i++ {
		numberFitnessTotal = numberFitnessTotal + float32(m.populationFitness[i])
	}

	//if total fitness number is zero each individual will be chosen with equal density
	if numberFitnessTotal == 0 {
		for i := 0; i < len(m.selectedParents); i++ {
			m.selectedParentsDensityFunction[i] = 1 / float32(m.populationSize)
		}
	} else {

		m.selectedParentsDensityFunction[0] = float32(m.populationFitness[0]) / numberFitnessTotal

		for i := 1; i < m.populationSize; i++ {
			m.selectedParentsDensityFunction[i] = m.selectedParentsDensityFunction[i-1] + float32(m.populationFitness[i])/numberFitnessTotal
		}
	}
}

func (m *FindwordGA) crossoverPopulationDensityFunction() {
	//crossover population using a buffer for new population
	//TODO add better method for crossover individuals replace selected parents
	var areEqual = true
	var chosenMate int
	var chosenMate2 int
	individual1 := make([]byte, m.solutionLength)
	individual2 := make([]byte, m.solutionLength)
	m.populationBuffer = make([][]byte, m.populationSize)

	for i := 0; i < m.populationSize; i++ {
		areEqual = true
		//chose a mate from the parents array which is a density stair function
		randomNumber := m.random.Float32()
		for j := 0; j < m.populationSize; j++ {
			if randomNumber <= m.selectedParentsDensityFunction[j] {
				chosenMate = j
				break
			}
		}

		randomNumber = m.random.Float32()
		for j := 0; j < m.populationSize; j++ {
			if randomNumber <= m.selectedParentsDensityFunction[j] {
				chosenMate2 = j
				break
			}
		}
		var buffer = make([]byte, m.solutionLength)
		individual1 = m.population[chosenMate2]
		individual2 = m.population[chosenMate]

		m.populationBuffer[i] = make([]byte, m.solutionLength)
		//write current individual to populationBuffer and check if it is equal to it's mate		)
		for j := 0; j < m.solutionLength; j++ {
			buffer[j] = individual1[j]
			if buffer[j] != individual2[j] {
				areEqual = false
			}
		}

		if areEqual == true {
			//both individuals are equal skip crossover and repeat
			i--
		} else {
			//note equal pair them into the bufferPopulation
			for j := 0; j < m.solutionLength/2; j++ {
				buffer[j] = individual2[j]
			}
			m.populationBuffer[i] = buffer
		}
	}

	m.population = m.populationBuffer
}

func (m *FindwordGA) mutatePopulation() {
	//TODO Umverteilung auf genau n aus population und keine doppelt ???

	individualsLength := m.solutionLength
	for i := 0; i < m.populationSize; i++ {
		if m.random.Float64() < m.mutationRate {
			m.population[i][m.random.Intn(individualsLength)] = m.alphabet[m.random.Intn(len(m.alphabet))]
		}
	}
}

func (m *FindwordGA) savePop() {
	deepCpy := make([][]byte, m.populationSize)
	for i := 0; i < m.populationSize; i++ {
		deepCpy[i] = make([]byte, m.solutionLength)
	}

	for i := 0; i < len(deepCpy); i++ {
		for j := 0; j < m.solutionLength; j++ {
			deepCpy[i][j] = m.population[i][j]
		}
	}
	m.populationHistory = append(m.populationHistory, deepCpy)
}

func (m *FindwordGA) calcStatistics() {

	//calc average fitness
	var avgFit float32 = 0
	for _, individualFit := range m.populationFitness {
		avgFit += float32(individualFit)
	}
	avgFit = avgFit / float32(len(m.populationFitness))
	m.AverageFitnessInPercent = avgFit / float32(m.solutionLength)

	//calc standard derivation form average fitness
	var avgDeviationFit float32 = 0
	for _, num := range m.populationFitness {
		avgDeviationFit += (float32(num) - avgFit) * (float32(num) - avgFit)
	}
	avgDeviationFit = avgDeviationFit / float32(len(m.populationFitness))

	m.StandardDeviationInPercent = float32(math.Sqrt(float64(avgDeviationFit))) / float32(m.solutionLength)
}
