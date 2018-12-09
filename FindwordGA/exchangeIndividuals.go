package FindwordGA

import (
	"sync"
)

type exchangeGA struct {
	numberOfIndividualsToBeExchanged int
	numberOfSendRcvChannel           int
	exchangePercent                  float32
	sendChannel                      []chan [][]byte
	receiveChannel                   []chan [][]byte
	waitgroup                        *sync.WaitGroup
	exchangeInterval                 int
}

//setup to exchange individuals with other GA's running simultaneously
func (m *FindwordGA) InitNetwork(exchangePercent float32,
	initSendChannel []chan [][]byte,
	initReceiveChannels []chan [][]byte,
	wg *sync.WaitGroup,
	doneChannel chan bool,
	doneChannelMutex *sync.Mutex,
	exchangeInterval int) {

	m.doneChannel = doneChannel
	m.doneChannelMutex = doneChannelMutex

	m.exchangeGA.numberOfIndividualsToBeExchanged = int(float32(m.populationSize) * exchangePercent)
	m.exchangeGA.numberOfSendRcvChannel = len(initSendChannel)

	m.exchangeGA.sendChannel = initSendChannel
	m.exchangeGA.receiveChannel = initReceiveChannels
	m.exchangeGA.exchangeInterval = exchangeInterval
	m.exchangeGA.waitgroup = wg

}

// Send and receive best individuals from connected GAs either synchronous or asynchronous
func (m *FindwordGA) exchangeIndividuals() {
	ex := m.exchangeGA
	//init array for individuals
	fitnessCurrentlySearched := m.solutionLength
	var individualsToBeSend = make([][]byte, m.exchangeGA.numberOfIndividualsToBeExchanged)
	for i := 0; i < m.exchangeGA.numberOfIndividualsToBeExchanged; i++ {
		individualsToBeSend[i] = make([]byte, m.solutionLength)
	}

	//collect best individuals to send
	for i := 0; i < m.exchangeGA.numberOfIndividualsToBeExchanged && len(ex.sendChannel) != 0; {
		for j := 0; j < len(m.populationFitness); j++ {
			if m.populationFitness[j] == fitnessCurrentlySearched {

				//individualsToBeSend [i] = []byte("  wordfindGA_ID send this")
				//individualsToBeSend[i][0] = byte(m.gaID)
				for k := 0; k < m.solutionLength; k++ {
					individualsToBeSend[i][k] = m.population[j][k]
				}
				i++
				if i >= m.exchangeGA.numberOfIndividualsToBeExchanged {
					break
				}
			}
		}
		fitnessCurrentlySearched--
	}

	//send N best lokal individuals in channel
	for i := 0; i < len(ex.sendChannel); i++ {
		//fmt.Printf("wordfindGA_ID %d SENDS  %d individuals \n", m.gaID, len(individualsToBeSend))
		//gehe in jeden channel leeren wenn inhalt sonst schreiben
		select {
		case <-ex.sendChannel[i]:
			ex.sendChannel[i] <- individualsToBeSend
		default:
			ex.sendChannel[i] <- individualsToBeSend
		}
	}

	//get all available new individuals
	var receivedIndividuals = make([][]byte, 0)
	var receivedIndividualsCopyBuffer [][]byte
	for i := 0; i < len(ex.receiveChannel); i++ {
		//test each channel empty it has content and write data
		select {
		case rcv := <-ex.receiveChannel[i]:
			//fmt.Printf("wordfindGA_ID %d RCV  %d individuals \n", m.gaID, len(rcv))

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
