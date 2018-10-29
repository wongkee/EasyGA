package easyga

import (
	"errors"
	"math/rand"
)

type Parameters struct {
	CrossoverProbability float64
	MutationProbability  float64
	PopulationSize       int
	Genotype             uint8
	ChromosomeLength     int
	IterationsLimit      int
}

type GeneticAlgorithm struct {
	Params     Parameters
	Iteration  int
	CheckStopFunction func (ga *GeneticAlgorithm) bool
	Population population
}

func (ga *GeneticAlgorithm) Init(parameters Parameters) error {
	if err := checkParam(parameters); err != nil {
		return err
	}

	ga.Params = parameters
	ga.Population.Init(ga.Params.ChromosomeLength, ga.Params.PopulationSize, ga.Params.Genotype)
	ga.Iteration = 0
	return nil
}

func (ga *GeneticAlgorithm) Run() (best Chromosome, fitness float64, iteration int) {
	for !ga.checkStop() {
		// Initialization
		var nextPopulation population

		// Selection - Select parents from population
		parentsPair := ga.selectParents()

		// Crossover - perform crossover on parents creating population
		for i := 0; i < len(parentsPair); i++ {
			parents := parentsPair[i]
			var child1, child2 Chromosome

			if rand.Float64() < ga.Params.CrossoverProbability {
				child1, child2 = parents[0].Crossover(parents[1])
			} else {
				child1, child2 = parents[0], parents[1]
			}

			nextPopulation.chromosomes = append(nextPopulation.chromosomes, child1, child2)
		}

		// Mutation - perform mutation of population
		for i := range nextPopulation.chromosomes {
			if rand.Float64() < ga.Params.MutationProbability {
				go nextPopulation.chromosomes[i].Mutate(ga.Params.Genotype)
			}
		}

		ga.Population = nextPopulation
		ga.Iteration++
	}

	bestIndex, bestFitness := ga.Population.FindBest()

	best = ga.Population.chromosomes[bestIndex]
	fitness = bestFitness
	iteration = ga.Iteration

	return
}

func (ga *GeneticAlgorithm) tournament() (newPopulation population) {
	for i := 0; i < ga.Params.PopulationSize; i++ {
		chromosome1 := ga.Population.chromosomes[getRandomChromosomeIndex(&ga.Population)]
		chromosome2 := ga.Population.chromosomes[getRandomChromosomeIndex(&ga.Population)]

		if chromosome1.fitness > chromosome2.fitness {
			newPopulation.chromosomes = append(newPopulation.chromosomes, chromosome1)
		} else {
			newPopulation.chromosomes = append(newPopulation.chromosomes, chromosome2)
		}
	}

	return
}

func (ga *GeneticAlgorithm) selectParents() (parentsPair [][2]Chromosome) {
	selectedPopulation := ga.tournament()

	for i := 0; i < ga.Params.PopulationSize/2; i++ {
		parent1, parent2 := selectedPopulation.chromosomes[2*i], selectedPopulation.chromosomes[2*i+1]
		parentsPair = append(parentsPair, [2]Chromosome{parent1, parent2})
	}

	return parentsPair
}

func (ga *GeneticAlgorithm) checkStop() bool {
	if ga.CheckStopFunction != nil {
		return ga.CheckStopFunction(ga)
	}

	_, bestFitness := ga.Population.FindBest()
	maybeBest := int(ga.Params.Genotype-1) * ga.Params.ChromosomeLength

	if int(bestFitness) >= maybeBest || ga.Iteration >= ga.Params.IterationsLimit {
		return true
	}

	return false
}

func checkParam(param Parameters) error {
	if param.CrossoverProbability < 0 || param.CrossoverProbability > 1 {
		return errors.New("CrossoverProbability should be in [0, 1]")
	}
	if param.MutationProbability < 0 || param.MutationProbability > 1 {
		return errors.New("MutationProbability should be in [0, 1]")
	}
	if param.PopulationSize <= 2 {
		return errors.New("PopulationSize should > 2")
	}
	if param.Genotype <= 1 {
		return errors.New("Genotype should > 1")
	}
	if param.ChromosomeLength <= 0 {
		return errors.New("ChromosomeLength should > 0")
	}
	if param.IterationsLimit <= 0 {
		return errors.New("IterationsLimit should > 0")
	}
	return nil
}