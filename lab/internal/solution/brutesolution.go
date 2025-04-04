package solution

import (
	"lab/internal/knapsack"
	"time"
)

type BruteSolution struct {
	Solutions         [][]int
	TimeFirstSolution float64
	TimeAllSolutions  float64
}

func FindSolutionByBruteForce(task knapsack.KnapsackTask) BruteSolution {
	startTime := time.Now()

	n := len(task.Vector)
	var solutions [][]int
	var timeFirstSolution float64
	var timeAllSolutions float64

	for mask := 0; mask < (1 << n); mask++ {
		weight := 0
		var subset []int

		for i := 0; i < n; i++ {
			if mask&(1<<i) > 0 {
				weight += task.Vector[i]
			}
		}

		if weight == task.TargetWeight {
			for i := 0; i < n; i++ {
				if mask&(1<<i) > 0 {
					subset = append(subset, task.Vector[i])
				}
			}

			solutions = append(solutions, subset)
			timeFirstSolution = time.Since(startTime).Seconds()
		}
	}

	timeAllSolutions = time.Since(startTime).Seconds()

	return BruteSolution{
		Solutions:         solutions,
		TimeFirstSolution: timeFirstSolution,
		TimeAllSolutions:  timeAllSolutions,
	}
}
