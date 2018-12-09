package FindwordGA

//calculate the probability a individual will be chosen for crossover with a density function
func (m *FindwordGA) selectParentsDensityFunction() {
	// density functions may have better performance with huge populations.
	//because with array density the parents array can reach size of lenghtOf(int) * population * (1+avgFitness) with avgFitness [0,1.0]
	//whereas the parentsDensityFuntion Array always has the same size of population * sizeOf(float32)

	//erstellen eines stufigen wahrscheinlichkeitsdiagrammes w(i)= i + SUM(0,1)
	//kandidaten mit einer wahrscheinlichkeit von 0 ergeben keine abstufung und der vorangegangene kanidat wird gewählt
	var numberFitnessTotal float32 = 0
	for i := 0; i < m.populationSize; i++ {
		numberFitnessTotal = numberFitnessTotal + float32(m.populationFitness[i])
	}

	//if total fitness number is zero each individual will be choosen with equal density
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

//crossover all individuals in the population and replace old generation with density function
func (m *FindwordGA) crossoverPopulationDensityFunction() {
	//crossover population using a buffer for new population
	//TODO add better method with channels for crossover individuals replace selected parents array
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
