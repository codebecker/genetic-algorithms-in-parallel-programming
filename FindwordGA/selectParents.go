package FindwordGA

import (
	"fmt"
)

//calculates the probability of each individual to be chosen for crossover by its fitness in a single thread
func (m *FindwordGA) selectParentsSingleThread() {
	//the fitter the individuals the more often their index will appear in the selectedParents array
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

//calculates the probability of each individual to be chosen for crossover by its fitness in multiple threads
func (m *FindwordGA) selectParentsMultiThread() {
	//the fitter the individuals the more often their index will appear in the selectedParents array
	fitnessDistribution := make([]int, m.numberOfThreads)
	numberFitnessTotal := 0

	for i := 0; i < m.numberOfThreads; i++ {
		go m.selectParentsCountThread(&fitnessDistribution[i], m.popBlockStart[i], m.popBlockEnd[i])
	}

	for i := 0; i < m.numberOfThreads; i++ {
		numberFitnessTotal += <-m.threadsDoneChannel
	}

	//check if fitness number is
	if numberFitnessTotal == 0 {
		m.selectedParents = make([]int, m.populationSize)
		for i := 0; i < len(m.selectedParents); i++ {
			m.selectedParents[i] = 1
		}

	} else {
		m.selectedParents = make([]int, numberFitnessTotal)
		offset := 0

		for i := 0; i < m.numberOfThreads; i++ {
			go m.selectParentsWriteThread(m.popBlockStart[i], m.popBlockEnd[i],
				offset)
			offset += fitnessDistribution[i]
		}

		for i := 0; i < m.numberOfThreads; i++ {
			<-m.threadsDoneChannel
		}

	}
}

//calculates the probability of each individual to be chosen for crossover by its fitness in multiple threads
func (m *FindwordGA) selectParentsCountThread(countSolution *int, start int, end int) {

	for i := start; i <= end; i++ {
		*countSolution += m.populationFitness[i]
	}
	m.threadsDoneChannel <- *countSolution
}

//calculates the probability of each individual to be chosen for crossover by its fitness in multiple threads
func (m *FindwordGA) selectParentsWriteThread(start int, end int, offset int) {
	for i := start; i <= end; i++ {
		for j := 0; j < m.populationFitness[i]; j++ {
			m.selectedParents[offset] = i
			offset++
		}
	}
	m.threadsDoneChannel <- 1
}

func (m *FindwordGA) selectParentsTest() {

	m.selectParentsSingleThread()
	// coppy solution
	singleResult := make([]int, len(m.selectedParents))
	for i := 0; i < len(m.selectedParents); i++ {
		singleResult[i] = m.selectedParents[i]
	}

	m.selectParentsMultiThread()

	if len(m.selectedParents) != len(singleResult) {
		fmt.Println("selectParents single and multi results are not equally long ")
		return
	} else {
		//compare content
		for i := 0; i < len(m.selectedParents); i++ {
			if singleResult[i] != m.selectedParents[i] {
				fmt.Println("selectParents single and multi results are different at ", i)
				return
			}
		}
	}
	//fmt.Println("selectParents single and multi results are SAME SAME")
}
