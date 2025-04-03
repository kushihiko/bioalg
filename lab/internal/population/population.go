package population

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	TournamentSize = 5
)

type GAConfig struct {
	GeneMin float64
	GeneMax float64
}

type Chromosome struct {
	Genes []bool
}

func NewChromosome(chromosomeSize int) Chromosome {
	genes := make([]bool, chromosomeSize)
	for i := range genes {
		genes[i] = rand.Intn(2) == 1
	}
	return Chromosome{Genes: genes}
}

type Population struct {
	Chromosomes    []Chromosome
	PopulationSize int
	ChromosomeSize int
	MutationRate   float64
	CrossoverRate  float64
	Fitness        func(Chromosome) int
}

func NewPopulation(populationSize, chromosomeSize int, mutationRate, crossoverRate float64, fitness func(Chromosome) int) *Population {
	pop := &Population{
		Chromosomes:    make([]Chromosome, populationSize),
		MutationRate:   mutationRate,
		CrossoverRate:  crossoverRate,
		PopulationSize: populationSize,
		ChromosomeSize: chromosomeSize,
		Fitness:        fitness,
	}

	for i := range pop.Chromosomes {
		pop.Chromosomes[i] = NewChromosome(chromosomeSize)
	}
	return pop
}

func (p *Population) Print() {
	for i := range p.Chromosomes {
		for j := range p.Chromosomes[i].Genes {
			if p.Chromosomes[i].Genes[j] {
				fmt.Printf("1")
			} else {
				fmt.Printf("0")
			}
		}
		fmt.Println()
	}
}

func (p *Population) AverageFitness() float64 {
	var sum int

	for _, chrom := range p.Chromosomes {
		sum += p.Fitness(chrom)
	}

	return float64(sum) / float64(len(p.Chromosomes))
}

func (p *Population) GetBest() Chromosome {
	minFitness := math.MaxInt
	x := p.Chromosomes[0]

	for _, chrom := range p.Chromosomes {
		current := p.Fitness(chrom)
		if current < minFitness {
			minFitness = current
			x = chrom
		}
	}

	return x
}

func (p *Population) EvolvePopulation() {
	newChromosomes := make([]Chromosome, p.PopulationSize)

	for i := 0; i < p.PopulationSize/2; i++ {
		parent1 := p.tournamentSelection(TournamentSize)
		parent2 := p.tournamentSelection(TournamentSize)

		child1, child2 := p.crossover(parent1, parent2)

		p.mutate(child1)
		p.mutate(child2)

		newChromosomes[2*i] = child1
		newChromosomes[2*i+1] = child2
	}

	p.Chromosomes = newChromosomes
}

func (p *Population) crossover(parent1, parent2 int) (Chromosome, Chromosome) {
	if rand.Float64() > p.CrossoverRate {
		return p.Chromosomes[parent1], p.Chromosomes[parent2]
	}

	point := rand.Intn(p.ChromosomeSize)
	child1 := make([]bool, p.ChromosomeSize)
	child2 := make([]bool, p.ChromosomeSize)

	copy(child1[point:], p.Chromosomes[parent2].Genes[point:])
	copy(child1[:point], p.Chromosomes[parent1].Genes[:point])
	copy(child2[point:], p.Chromosomes[parent1].Genes[point:])
	copy(child2[:point], p.Chromosomes[parent2].Genes[:point])

	return Chromosome{Genes: child1}, Chromosome{Genes: child2}
}

func (p *Population) tournamentSelection(k int) int {
	nums := rand.Perm(p.PopulationSize)
	nums = nums[:k]

	best := nums[0]
	for i := 1; i < k; i++ {
		competitor := nums[i]
		if p.Fitness(p.Chromosomes[competitor]) < p.Fitness(p.Chromosomes[best]) {
			best = competitor
		}
	}

	return best
}

func (p *Population) mutate(chromosome Chromosome) {
	if rand.Float64() > p.MutationRate {
		return
	}

	index := rand.Intn(p.ChromosomeSize)
	chromosome.Genes[index] = !chromosome.Genes[index]
}
