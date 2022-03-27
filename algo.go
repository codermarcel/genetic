package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// MutationRate is the rate of mutation
// var MutationRate = 0.005

// PopSize is the size of the population
// var PopSize = 500

// type GeneticSolver struct {

// }

func approximateResult(tar Nutrients,
	giv Recipe,
	debug bool,
	maxGenerations int,
	abortChan chan bool,
	popSize int,
	okPercentage float64) Problem {

	var abort bool
	target := tar
	givenRecep := giv.DeepClone()

	probl := NewProblem(givenRecep, givenRecep, target)

	start := time.Now()
	algoID := start.UnixNano()

	population := createPopulation(target, probl, popSize)
	star := Problem{Fitness: 0}

	if debug {
		fmt.Printf("Start: %d\n---------\n", algoID)
	}

	found := false
	generation := 0
	for !found && !abort {
		select {
		case <-abortChan:
			abort = true
		default:
			//dont block here
		}

		if generation >= maxGenerations {
			abort = true
		}
		generation++
		bestOrganism := getBest(population)

		if bestOrganism.Fitness > star.Fitness {
			star = bestOrganism.DeepClone()
		}

		if debug {
			fmt.Printf("%d fitness: %2f | Names -- | %s | generation: %d \n", algoID, bestOrganism.Fitness, bestOrganism.Given.PrettyPrintNames(), generation)
		}

		if compare(bestOrganism, okPercentage) {
			found = true
		} else {
			maxFitness := bestOrganism.Fitness
			pool := createPool(population, maxFitness)

			if len(pool) < 2 {
				pool = createPopulation(target, probl, popSize)
			}

			population = naturalSelection(pool, population)
		}

	}
	elapsed := time.Since(start)

	if debug {
		fmt.Printf("Time taken for %d: %s\n---------\n\n\n", algoID, elapsed)
	}

	star.Done = elapsed

	return star
}

func compare(best Problem, percentage float64) bool {
	pro := best.Given.AllNutrients().Pro(best.Target)
	// fmt.Println("pro :", pro)
	return pro >= percentage
}

// creates an Organism
func createOrganismOrg(target Nutrients, p Problem) (organism Problem) {
	totalItems := rand.Int31n(int32(len(p.Given.Ingredients))) * 2
	ingredients := make([]Ingredient, totalItems)

	for i := int32(0); i < totalItems; i++ {
		someIndex := rand.Int31n(int32(len(p.Given.Ingredients)))
		ingredients[i] = p.Given.Ingredients[someIndex]

	}

	newP := p.DeepClone()
	newP.Given.Ingredients = ingredients
	newP.calcFitness()

	return newP
}

// creates the initial population
func createPopulation(target Nutrients, p Problem, popSize int) (population []Problem) {
	fu := createOrganismOrg

	population = make([]Problem, popSize)
	for i := 0; i < popSize; i++ {
		population[i] = fu(target, p)
	}
	return
}

// create the breeding pool that creates the next generation
func createPool(population []Problem, maxFitness float64) []Problem {
	pool := make([]Problem, 0)
	// create a pool for next generation
	for i := 0; i < len(population); i++ {
		population[i].calcFitness()
		num := int((population[i].Fitness / maxFitness) * 100)
		for n := 0; n < num; n++ {
			pool = append(pool, population[i])
		}
	}
	return pool
}

func naturalSelection(pool []Problem, population []Problem) []Problem {
	//if !oldMode {
	//	return naturalSelectionNew(pool, population)
	//}

	next := make([]Problem, len(population))

	for i := 0; i < len(population); i++ {
		r1, _ := rand.Intn(len(pool)), rand.Intn(len(pool))
		a := pool[r1]
		// b := pool[r2]

		// child := crossover(a, b)
		child := a.DeepClone()
		child.mutate()
		child.calcFitness()

		next[i] = child

	}

	return next
}

// perform natural selection to create the next generation
func naturalSelectionNew(pool []Problem, population []Problem) []Problem {
	next := make([]Problem, len(population))

	var wg sync.WaitGroup
	wg.Add(len(population))

	for i := 0; i < len(population); i++ {
		go func(i int) {
			defer wg.Done()
			r1, r2 := rand.Intn(len(pool)), rand.Intn(len(pool))
			a := pool[r1]
			b := pool[r2]

			// child := crossover(a, b)
			child := a.crossover(b)
			child.mutate()
			child.calcFitness()

			next[i] = child
		}(i)
	}

	wg.Wait()
	return next
}

// // perform natural selection to create the next generation
// func naturalSelection(pool []Problem, population []Problem) []Problem {
// 	next := make([]Problem, len(population))

// 	for i := 0; i < len(population); i++ {
// 		r1, r2 := rand.Intn(len(pool)), rand.Intn(len(pool))
// 		a := pool[r1]
// 		b := pool[r2]

// 		// child := crossover(a, b)
// 		child := a.crossover(b)
// 		child.mutate()
// 		child.calcFitness()

// 		next[i] = child
// 	}
// 	return next
// }

// Get the best organism
func getBest(population []Problem) Problem {
	best := 0.0
	index := 0
	for i := 0; i < len(population); i++ {
		if population[i].Fitness > best {
			index = i
			best = population[i].Fitness
		}
	}
	return population[index]
}
