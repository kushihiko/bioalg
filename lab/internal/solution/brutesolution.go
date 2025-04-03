package solution

import (
	"lab/internal/knapsack"
	"time"
)

type BruteSolution struct {
	Solutions         [][]int
	TimeFirstSolution int64
	TimeAllSolutions  int64
}

func FindSolutionByBruteForce(task knapsack.KnapsackTask) BruteSolution {
	startTime := time.Now()

	n := len(task.Vector)
	var solutions [][]int
	var timeFirstSolution int64
	var timeAllSolutions int64

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
			timeFirstSolution = time.Since(startTime).Milliseconds()
		}
	}

	timeAllSolutions = time.Since(startTime).Milliseconds()

	return BruteSolution{
		Solutions:         solutions,
		TimeFirstSolution: timeFirstSolution,
		TimeAllSolutions:  timeAllSolutions,
	}
}
