package knapsack

import (
	"math/rand"
)

type KnapsackTask struct {
	Vector       []int
	VectorNumber int
	TargetWeight int
	Fraction     float64
	Amax         int
	VectorLength int
}

type VectorGenerator interface {
	GenerateVector() []int
	GetAMax() int
	GetLength() int
}

type TaskGenerator struct {
	vectorGenerator VectorGenerator
	rand            *rand.Rand
}

func NewTaskGenerator(vectorGenerator VectorGenerator, seed int64) *TaskGenerator {
	return &TaskGenerator{
		vectorGenerator: vectorGenerator,
		rand:            rand.New(rand.NewSource(seed)),
	}
}

func (tg *TaskGenerator) generateTask(vectorNumber, tasksNumber int) []KnapsackTask {
	vector := tg.vectorGenerator.GenerateVector()
	tasks := make([]KnapsackTask, tasksNumber)

	for i := range tasksNumber {

		n := len(vector)

		fraction := 0.1 + tg.rand.Float64()*0.4
		numItems := max(int(float64(n)*fraction), 1)

		indices := tg.rand.Perm(n)[:numItems]

		targetWeight := 0
		for _, idx := range indices {
			targetWeight += vector[idx]
		}

		dst := make([]int, len(vector))
		copy(dst, vector)
		tasks[i] = KnapsackTask{
			Vector:       dst,
			VectorNumber: vectorNumber,
			TargetWeight: targetWeight,
			Fraction:     fraction,
			Amax:         tg.vectorGenerator.GetAMax(),
			VectorLength: tg.vectorGenerator.GetLength(),
		}
	}

	return tasks
}

func (tg *TaskGenerator) GenerateTasks(vectorNumber, taskNumber int) []KnapsackTask {
	var tasks []KnapsackTask

	for i := 0; i < vectorNumber; i++ {
		tasks = append(tasks, tg.generateTask(i+1, taskNumber)...)
	}

	return tasks
}
