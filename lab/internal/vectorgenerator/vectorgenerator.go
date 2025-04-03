package vectorgenerator

import (
	"math/rand"
)

type VectorGeneratorImpl struct {
	n      int
	amax   int
	random *rand.Rand
}

func New(n int, amax int, seed int64) *VectorGeneratorImpl {
	return &VectorGeneratorImpl{
		n:      n,
		amax:   amax,
		random: rand.New(rand.NewSource(seed)),
	}
}

func (kg *VectorGeneratorImpl) GenerateVector() []int {
	vector := make([]int, kg.n)
	for i := 0; i < kg.n; i++ {
		vector[i] = kg.random.Intn(kg.amax) + 1
	}
	return vector
}

func (kg *VectorGeneratorImpl) GetAMax() int {
	return kg.amax
}

func (kg *VectorGeneratorImpl) GetLength() int {
	return kg.n
}
