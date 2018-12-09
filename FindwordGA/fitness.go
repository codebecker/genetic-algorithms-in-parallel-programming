package FindwordGA

//calculates the fitness of each individual in the population in a single thread
func (m *FindwordGA) calcFitnessSingleThread() int {
	fittest := -1
	for i := 0; i < len(m.population); i++ {
		fitnessOfIndividual := 0
		for j := 0; j < m.solutionLength; j++ {
			if m.population[i][j] == m.solution[j] {
				fitnessOfIndividual++
			}
		}
		m.populationFitness[i] = fitnessOfIndividual

		//if a solution was found
		if fitnessOfIndividual == m.solutionLength && fittest < 0 {
			fittest = i
		}
	}
	return fittest
}

//calculates the fitness of each individual in the population in multiple threads
func (m *FindwordGA) calcFitnessMultiThread() int {
	solution := -1

	for i := 0; i < m.numberOfThreads; i++ {
		go m.calcFitnessOfIndividualsThread(m.threadsDoneChannel, m.popBlockStart[i], m.popBlockEnd[i])

	}

	for i := 0; i < m.numberOfThreads; i++ {
		tmp := <-m.threadsDoneChannel
		//if a solution was found and no other solution has been found before
		if tmp > -1 && solution == -1 {
			solution = tmp
		}
	}
	return solution
}

//calculates the fitness of a certain part of the population necessary for multithreading
func (m *FindwordGA) calcFitnessOfIndividualsThread(solutionChan chan int, start int, end int) int {
	var fittest = -1
	for i := start; i <= end; i++ {
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
