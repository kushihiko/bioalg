package main

import (
	"fmt"
	"lab/internal/knapsack"
	"lab/internal/population"
	"lab/internal/solution"
	"lab/internal/table"
	"lab/internal/vectorgenerator"
	"math"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	VectorLength             = 24
	VectorNumber             = 50 // 50
	TaskNumber               = 20 // 20
	Folder                   = "results/8/"
	VectorsTable             = "vectors.csv"
	TasksTable               = "tasks.csv"
	BruteForceSolutionsTable = "brutesolv.csv"
	GeneticSolutionTable     = "gensolv.csv"
)

const (
	PopulationSize = 750 // 750
	ChromosomeSize = VectorLength
	MutationRate   = 0.2
	CrossoverRate  = 0.85
	Module         = true
)

type geneticConfig struct {
	PopulationSize int
	ChromosomeSize int
	MutationRate   float64
	CrossoverRate  float64
	Module         bool
}

func main() {
	tasks := createTasks(VectorLength, VectorNumber, TaskNumber)
	SaveVectors(tasks, Folder+VectorsTable)
	SaveTasks(tasks, Folder+TasksTable)

	bruteSolutions := bruteForceSolveParallel(tasks)
	SaveBruteSolutions(bruteSolutions, Folder+BruteForceSolutionsTable)
	printBruteAverage(bruteSolutions)

	timeLimits := make([]float64, len(bruteSolutions))
	for i, bruteSolution := range bruteSolutions {
		timeLimits[i] = bruteSolution.TimeFirstSolution * 2
	}

	fmt.Println("Start genetic algorithm")

	genericSolutions := geneticSolveParallel(tasks, geneticConfig{
		PopulationSize: PopulationSize,
		ChromosomeSize: ChromosomeSize,
		MutationRate:   MutationRate,
		CrossoverRate:  CrossoverRate,
		Module:         Module,
	}, timeLimits)
	SaveGenericSolutions(genericSolutions, Folder+GeneticSolutionTable)
	printGenAverage(genericSolutions)
}

func createTasks(vectorLength, vectorNumber, taskNumber int) []knapsack.KnapsackTask {
	seed := time.Now().UnixNano()
	aMax := int(1 << int(math.Round(float64(vectorLength)/1.4)))

	vectorGen := vectorgenerator.New(vectorLength, aMax, seed)
	taskGen := knapsack.NewTaskGenerator(vectorGen, seed)
	tasks := taskGen.GenerateTasks(vectorNumber, taskNumber)

	return tasks
}

func bruteForceSolveParallel(tasks []knapsack.KnapsackTask) []solution.BruteSolution {
	var solutions []solution.BruteSolution
	solutions = make([]solution.BruteSolution, len(tasks))
	var wg sync.WaitGroup

	for i, task := range tasks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			solutions[i] = solution.FindSolutionByBruteForce(task)
		}()

	}

	wg.Wait()
	return solutions
}

func bruteForceSolve(tasks []knapsack.KnapsackTask) []solution.BruteSolution {
	var solutions []solution.BruteSolution
	solutions = make([]solution.BruteSolution, len(tasks))

	for i, task := range tasks {
		solutions[i] = solution.FindSolutionByBruteForce(task)
	}

	return solutions
}

func geneticSolveParallel(tasks []knapsack.KnapsackTask, config geneticConfig, timeLimits []float64) []solution.GeneticSolution {
	solver := solution.GeneticSolver{
		PopulationSize: config.PopulationSize,
		ChromosomeSize: config.ChromosomeSize,
		MutationRate:   config.MutationRate,
		CrossoverRate:  config.CrossoverRate,
	}
	solutions := make([]solution.GeneticSolution, len(tasks))
	var wg sync.WaitGroup

	for i, task := range tasks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			solutions[i] = solver.Solve(GeneticTaskImpl{Task: task, Module: config.Module}, timeLimits[i])
		}()
	}

	wg.Wait()
	return solutions
}

func geneticSolve(tasks []knapsack.KnapsackTask, config geneticConfig, timeLimits []float64) []solution.GeneticSolution {
	solver := solution.GeneticSolver{
		PopulationSize: config.PopulationSize,
		ChromosomeSize: config.ChromosomeSize,
		MutationRate:   config.MutationRate,
		CrossoverRate:  config.CrossoverRate,
	}
	solutions := make([]solution.GeneticSolution, len(tasks))

	for i, task := range tasks {
		solutions[i] = solver.Solve(GeneticTaskImpl{Task: task, Module: config.Module}, timeLimits[i])
	}

	return solutions
}

type GeneticTaskImpl struct {
	Task   knapsack.KnapsackTask
	Module bool
}

//func (gt GeneticTaskImpl) Fitness(x population.Chromosome) float64 {
//	sum := 0
//
//	for i, gene := range x.Genes {
//		if gene && i < gt.Task.VectorLength {
//			sum += gt.Task.Vector[i]
//		}
//	}
//
//	// return math.Abs(float64((sum - gt.Task.TargetWeight) % (gt.Task.Amax + 1)))
//	return math.Abs(float64(sum - gt.Task.TargetWeight))
//}

func (gt GeneticTaskImpl) Fitness(x population.Chromosome) int {
	sum := 0

	for i, gene := range x.Genes {
		if gene && i < gt.Task.VectorLength {
			sum += gt.Task.Vector[i]
		}
	}

	diff := sum - gt.Task.TargetWeight
	if diff < 0 {
		diff = -diff
	}

	return diff
}

func (gt GeneticTaskImpl) FitnessModule(x population.Chromosome) int {
	sum := 0
	module := gt.Task.Amax + 1

	for i, gene := range x.Genes {
		if gene && i < gt.Task.VectorLength {
			sum += gt.Task.Vector[i]
		}
	}

	diff := (sum % module) - (gt.Task.TargetWeight % module)
	if diff < 0 {
		diff = -diff
	}

	return diff
}

func (gt GeneticTaskImpl) GetTask() knapsack.KnapsackTask {
	return gt.Task
}

func (gt GeneticTaskImpl) IsModule() bool {
	return gt.Module
}

func SaveVectors(tasks []knapsack.KnapsackTask, filename string) {
	_ = os.Remove(filename)

	tbl, err := table.New(filename)
	defer func() {
		err = tbl.Close()
		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err)
	}

	err = tbl.Write([]string{"number", "vector", "amax"})
	if err != nil {
		panic(err)
	}

	for i, task := range tasks {
		if i == 0 || tasks[i-1].VectorNumber != task.VectorNumber {
			err = tbl.Write([]string{strconv.Itoa(task.VectorNumber), fmt.Sprint(task.Vector), strconv.Itoa(task.Amax)})
			if err != nil {
				panic(err)
			}
		}
	}

	err = tbl.Flush()
	if err != nil {
		panic(err)
	}
}

func SaveTasks(tasks []knapsack.KnapsackTask, filename string) {
	_ = os.Remove(filename)

	tbl, err := table.New(filename)
	defer func() {
		err = tbl.Close()
		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err)
	}

	err = tbl.Write([]string{"task number", "vector number", "target weight", "fraction"})
	if err != nil {
		panic(err)
	}

	for i, task := range tasks {
		err = tbl.Write([]string{strconv.Itoa(i + 1), strconv.Itoa(task.VectorNumber), strconv.Itoa(task.TargetWeight), fmt.Sprintf("%.3f", task.Fraction)})
		if err != nil {
			panic(err)
		}
	}

	err = tbl.Flush()
	if err != nil {
		panic(err)
	}
}

func SaveBruteSolutions(solutions []solution.BruteSolution, filename string) {
	_ = os.Remove(filename)

	tbl, err := table.New(filename)
	defer func() {
		err = tbl.Close()
		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err)
	}

	err = tbl.Write([]string{"task number", "first solution time", "all solutions time", "answers number"})
	if err != nil {
		panic(err)
	}

	for i, solution := range solutions {
		err = tbl.Write([]string{strconv.Itoa(i + 1), fmt.Sprintf("%.3f", solution.TimeFirstSolution), fmt.Sprintf("%.3f", solution.TimeAllSolutions), strconv.Itoa(len(solution.Solutions))})
		if err != nil {
			panic(err)
		}
	}

	err = tbl.Flush()
	if err != nil {
		panic(err)
	}
}

func printBruteAverage(solutions []solution.BruteSolution) {
	avgFirst := 0.0
	avgAll := 0.0

	for _, solution := range solutions {
		avgFirst += float64(solution.TimeFirstSolution) / float64(len(solutions))
		avgAll += float64(solution.TimeAllSolutions) / float64(len(solutions))
	}

	fmt.Println("avgFirst:", avgFirst)
	fmt.Println("avgAll:", avgAll)
}

func SaveGenericSolutions(solutions []solution.GeneticSolution, filename string) {
	_ = os.Remove(filename)

	tbl, err := table.New(filename)
	defer func() {
		err = tbl.Close()
		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err)
	}

	err = tbl.Write([]string{"task number", "solution time", "fitness minimum", "stop reason", "generations number"})
	if err != nil {
		panic(err)
	}

	for i, solution := range solutions {
		err = tbl.Write([]string{strconv.Itoa(i + 1), fmt.Sprintf("%.3f", solution.Time), strconv.Itoa(int(solution.Fitness)), strconv.Itoa(int(solution.StopReason)), strconv.Itoa(int(solution.GenerationNumber))})
		if err != nil {
			panic(err)
		}
	}

	err = tbl.Flush()
	if err != nil {
		panic(err)
	}
}

func printGenAverage(solutions []solution.GeneticSolution) {
	avgTime := 0.0
	avgGenNumber := 0.0
	solved := 0

	for _, item := range solutions {
		if item.StopReason == solution.Found {
			avgTime += item.Time / float64(len(solutions))
			avgGenNumber += float64(item.GenerationNumber) / float64(len(solutions))
			solved++
		}
	}

	fmt.Println("avgTime:", avgTime)
	fmt.Println("avgGenNumber:", avgGenNumber)
	fmt.Println("solved:", solved)
	fmt.Println("Total:", len(solutions))
	fmt.Println("Доля:", float64(solved)/float64(len(solutions)))
}
