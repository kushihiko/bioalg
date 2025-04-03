package solution

import (
	"lab/internal/knapsack"
	"lab/internal/population"
	"sync"
	"time"
)

type Status int

const (
	Found Status = iota
	TimeExpired
	NoProgress
)

type GeneticSolution struct {
	Solution         []int
	Time             int64
	Fitness          int
	StopReason       Status
	GenerationNumber int64
}

type GeneticSolver struct {
	PopulationSize int
	ChromosomeSize int
	MutationRate   float64
	CrossoverRate  float64
}

func NewGeneticSolver(populationSize, chromosomeSize int, mutationRate, crossoverRate float64) *GeneticSolver {
	return &GeneticSolver{
		PopulationSize: populationSize,
		ChromosomeSize: chromosomeSize,
		MutationRate:   mutationRate,
		CrossoverRate:  crossoverRate,
	}
}

type GeneticTask interface {
	Fitness(population.Chromosome) int
	GetTask() knapsack.KnapsackTask
}

var mu sync.Mutex

func (gs GeneticSolver) Solve(geneticTask GeneticTask, timeLimit int64) GeneticSolution {
	pop := population.NewPopulation(gs.PopulationSize, gs.ChromosomeSize, gs.MutationRate, gs.CrossoverRate, geneticTask.Fitness)

	start := time.Now()

	current := -1.0
	prev := -1.0
	beforePrev := -1.0

	bestFitness := -1

	var stopReason Status
	generationNumber := int64(1)

	for {
		pop.EvolvePopulation()
		bestFitness = geneticTask.Fitness(pop.GetBest())
		current = pop.AverageFitness()
		generationNumber++

		//fmt.Printf("FITNESS: %f\n", current)

		if bestFitness == 0 {
			stopReason = Found
			break
		}

		if time.Since(start).Milliseconds() > timeLimit {
			stopReason = TimeExpired
			break
		}

		if beforePrev == prev && current == prev {
			stopReason = NoProgress
			break
		}

		beforePrev = prev
		prev = pop.AverageFitness()
		_ = beforePrev
	}

	best := pop.GetBest()
	var solution []int
	for i, gene := range best.Genes {
		if gene {
			solution = append(solution, geneticTask.GetTask().Vector[i])
		}
	}

	return GeneticSolution{
		Solution:         solution,
		Time:             time.Since(start).Milliseconds(),
		Fitness:          bestFitness,
		StopReason:       stopReason,
		GenerationNumber: generationNumber,
	}
}

//func nextGeneration(pop *population.Population) *population.Population {
//	pop.Reproduction()
//	pop.Crossover()
//	pop.Mutate()
//
//	return pop
//}
