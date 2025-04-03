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
	Folder                   = "results/4_3/"
	VectorsTable             = "vectors.csv"
	TasksTable               = "tasks.csv"
	BruteForceSolutionsTable = "brutesolv.csv"
	GeneticSolutionTable     = "gensolv.csv"
)

const (
	PopulationSize = 200 // 750
	ChromosomeSize = VectorLength
	MutationRate   = 0.2
	CrossoverRate  = 0.85
)

type geneticConfig struct {
	PopulationSize int
	ChromosomeSize int
	MutationRate   float64
	CrossoverRate  float64
}

func main() {
	tasks := createTasks(VectorLength, VectorNumber, TaskNumber)
	SaveVectors(tasks, Folder+VectorsTable)
	SaveTasks(tasks, Folder+TasksTable)

	bruteSolutions := bruteForceSolveParallel(tasks)
	SaveBruteSolutions(bruteSolutions, Folder+BruteForceSolutionsTable)
	printBruteAverage(bruteSolutions)

	timeLimits := make([]int64, len(bruteSolutions))
	for i, bruteSolution := range bruteSolutions {
		timeLimits[i] = bruteSolution.TimeFirstSolution * 2
	}

	fmt.Println("Start genetic algorithm")

	genericSolutions := geneticSolveParallel(tasks, geneticConfig{
		PopulationSize: PopulationSize,
		ChromosomeSize: ChromosomeSize,
		MutationRate:   MutationRate,
		CrossoverRate:  CrossoverRate,
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

func geneticSolveParallel(tasks []knapsack.KnapsackTask, config geneticConfig, timeLimits []int64) []solution.GeneticSolution {
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
			solutions[i] = solver.Solve(GeneticTaskImpl{Task: task}, timeLimits[i])
		}()
	}

	wg.Wait()
	return solutions
}

func geneticSolve(tasks []knapsack.KnapsackTask, config geneticConfig, timeLimits []int64) []solution.GeneticSolution {
	solver := solution.GeneticSolver{
		PopulationSize: config.PopulationSize,
		ChromosomeSize: config.ChromosomeSize,
		MutationRate:   config.MutationRate,
		CrossoverRate:  config.CrossoverRate,
	}
	solutions := make([]solution.GeneticSolution, len(tasks))

	for i, task := range tasks {
		solutions[i] = solver.Solve(GeneticTaskImpl{Task: task}, timeLimits[i])
	}

	return solutions
}

type GeneticTaskImpl struct {
	Task knapsack.KnapsackTask
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

func (gt GeneticTaskImpl) GetTask() knapsack.KnapsackTask {
	return gt.Task
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

	// First raw
	var str []string
	str = append(str, "task number")
	for i := range tasks {
		str = append(str, strconv.Itoa(i+1))
	}
	err = tbl.Write(str)
	if err != nil {
		panic(err)
	}

	// Second raw
	str = make([]string, 0)
	str = append(str, "vector number")
	for _, task := range tasks {
		str = append(str, strconv.Itoa(task.VectorNumber))
	}
	err = tbl.Write(str)
	if err != nil {
		panic(err)
	}

	// Third raw
	str = make([]string, 0)
	str = append(str, "target weight")
	for _, task := range tasks {
		str = append(str, strconv.Itoa(task.TargetWeight))
	}
	err = tbl.Write(str)
	if err != nil {
		panic(err)
	}

	// Fourth raw
	str = make([]string, 0)
	str = append(str, "fraction")
	for _, task := range tasks {
		str = append(str, fmt.Sprintf("%.3f", task.Fraction))
	}
	err = tbl.Write(str)
	if err != nil {
		panic(err)
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

	// First raw
	var str []string
	str = append(str, "task number")
	for i := range solutions {
		str = append(str, strconv.Itoa(i+1))
	}
	err = tbl.Write(str)
	if err != nil {
		panic(err)
	}

	// Second raw
	str = make([]string, 0)
	str = append(str, "first solution time")
	for _, solution := range solutions {
		str = append(str, strconv.Itoa(int(solution.TimeFirstSolution)))
	}
	err = tbl.Write(str)
	if err != nil {
		panic(err)
	}

	// Third raw
	str = make([]string, 0)
	str = append(str, "all solutions time")
	for _, solution := range solutions {
		str = append(str, strconv.Itoa(int(solution.TimeAllSolutions)))
	}
	err = tbl.Write(str)
	if err != nil {
		panic(err)
	}

	// Fourth raw
	str = make([]string, 0)
	str = append(str, "answers number")
	for _, solution := range solutions {
		str = append(str, strconv.Itoa(len(solution.Solutions)))
	}
	err = tbl.Write(str)
	if err != nil {
		panic(err)
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

	// First raw
	var str []string
	str = append(str, "task number")
	for i := range solutions {
		str = append(str, strconv.Itoa(i+1))
	}
	err = tbl.Write(str)
	if err != nil {
		panic(err)
	}

	// Second raw
	str = make([]string, 0)
	str = append(str, "solution time")
	for _, solution := range solutions {
		str = append(str, strconv.Itoa(int(solution.Time)))
	}
	err = tbl.Write(str)
	if err != nil {
		panic(err)
	}

	// Third raw
	str = make([]string, 0)
	str = append(str, "fitness minimum")
	for _, solution := range solutions {
		str = append(str, strconv.Itoa(int(solution.Fitness)))
	}
	err = tbl.Write(str)
	if err != nil {
		panic(err)
	}

	// Fourth raw
	str = make([]string, 0)
	str = append(str, "stop reason")
	for _, solution := range solutions {
		str = append(str, strconv.Itoa(int(solution.StopReason)))
	}
	err = tbl.Write(str)
	if err != nil {
		panic(err)
	}

	err = tbl.Flush()
	if err != nil {
		panic(err)
	}

	// Fifth raw
	str = make([]string, 0)
	str = append(str, "generation number")
	for _, solution := range solutions {
		str = append(str, strconv.Itoa(int(solution.GenerationNumber)))
	}
	err = tbl.Write(str)
	if err != nil {
		panic(err)
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
			avgTime += float64(item.Time) / float64(len(solutions))
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
