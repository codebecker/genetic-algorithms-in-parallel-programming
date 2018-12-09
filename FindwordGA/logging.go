package FindwordGA

import "math"

//calculate the average fitness and standard derivation of the current population
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

//save the current population to a slice for analysis
func (m *FindwordGA) logPopulation() {
	//implement logging to disk with a suitable package

}
