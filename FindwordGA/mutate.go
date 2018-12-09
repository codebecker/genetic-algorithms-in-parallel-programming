package FindwordGA

//randomly chose individuals and change a single char
func (m *FindwordGA) mutatePopulation() {

	individualsLength := m.solutionLength
	alphabetLength := len(m.alphabet)
	for i := 0; i < m.populationSize; i++ {
		if m.random.Float32() < m.mutationRate {
			m.population[i][m.random.Intn(individualsLength)] = m.alphabet[m.random.Intn(alphabetLength)]
		}
	}

}

//randomly chose individuals and change a single char
func (m *FindwordGA) mutatePopulationMultiThread() {

	mutationRateThreads := make([]float32, m.numberOfThreads)

	/*
		var distributionSum float32
		for i := 0; i < m.numberOfThreads; i++ {
			mutationRateThreads[i] = m.random.Float32()
			distributionSum += mutationRateThreads[i]
		}
		for i := 0; i < m.numberOfThreads; i++ {
			//random density distribution between mutation threads
			mutationRateThreads[i] = mutationRateThreads[i]/ distributionSum*m.mutationRate*float32(m.numberOfThreads)
		}
	*/

	for i := 0; i < m.numberOfThreads; i++ {
		//same density between all mutation threads seems to be faster
		mutationRateThreads[i] = m.mutationRate
	}

	for i := 0; i < m.numberOfThreads; i++ {
		go m.mutatePopulationThread(m.popBlockStart[i], m.popBlockEnd[i], m.threadsDoneChannel, i, mutationRateThreads[i])
	}

	for i := 0; i < m.numberOfThreads; i++ {
		<-m.threadsDoneChannel
	}
}

//randomly chose individuals and change a single char
func (m *FindwordGA) mutatePopulationThread(start int, end int, threadsDoneChannel chan int, threadNumber int, threadsMutationRate float32) {
	random := m.randomForThreads[threadNumber]
	for i := start; i <= end; i++ {
		if random.Float32() < threadsMutationRate {
			m.population[i][random.Intn(m.solutionLength)] = m.alphabet[random.Intn(len(m.alphabet))]
		}
	}
	threadsDoneChannel <- 1
}
