package FindwordGA

//crossover all individuals in the population an replace old generation
func (m *FindwordGA) crossoverPopulation() {
	var areEqual = true
	m.populationBuffer = make([][]byte, m.populationSize)
	individual1 := make([]byte, m.solutionLength)
	individual2 := make([]byte, m.solutionLength)
	buffer := make([]byte, m.solutionLength)

	for i := 0; i < m.populationSize; i++ {
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

//crossover all individuals in the population an replace old generation in multiple threads
func (m *FindwordGA) crossoverPopulationMultithread() {
	m.populationBuffer = make([][]byte, m.populationSize)

	for i := 0; i < m.numberOfThreads; i++ {
		go m.crossoverPopulationThread(i, m.threadsDoneChannel, m.popBlockStart[i], m.popBlockEnd[i])
	}
	for i := 0; i < m.numberOfThreads; i++ {
		<-m.threadsDoneChannel
	}
	m.population = m.populationBuffer
}

func (m *FindwordGA) crossoverPopulationThread(threadNumber int, threadDone chan int, start int, end int) {
	var areEqual = true
	individual1 := make([]byte, m.solutionLength)
	individual2 := make([]byte, m.solutionLength)
	buffer := make([]byte, m.solutionLength)
	random := m.randomForThreads[threadNumber]

	for i := start; i <= end; i++ {

		areEqual = true
		buffer = make([]byte, m.solutionLength)
		individual1 = m.population[m.selectedParents[random.Intn(len(m.selectedParents))]]
		individual2 = m.population[m.selectedParents[random.Intn(len(m.selectedParents))]]

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
	threadDone <- 1
}
