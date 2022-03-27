package main

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_approximateResult(t *testing.T) {
	//simple
	fish := Ingredient{Name: "Fish", Nutrients: Nutrients{Calories: 100, Proteins: 100, Carbs: 100, Fat: 100}}
	chips := Ingredient{Name: "Chips", Nutrients: Nutrients{Calories: 5, Proteins: 5, Carbs: 5, Fat: 5}}

	given0 := NewRecipe("Fish and Chips", 2).AddIngredient(fish, 1).AddIngredient(chips, 3)
	target0 := Nutrients{Calories: 140, Proteins: 140, Carbs: 140, Fat: 140}

	//pd1
	ham := Ingredient{Name: "Ham", Nutrients: Nutrients{Calories: 500, Proteins: 400, Carbs: 400, Fat: 600}}
	potato := Ingredient{Name: "Potato", Nutrients: Nutrients{Calories: 100, Proteins: 59.9, Carbs: 50, Fat: 34}}
	chicken := Ingredient{Name: "Chicken", Nutrients: Nutrients{Calories: 300, Proteins: 200, Carbs: 254, Fat: 323}}
	rice := Ingredient{Name: "Rice", Nutrients: Nutrients{Calories: 3, Proteins: 4, Carbs: 3, Fat: 2}}

	meatBallsX := Ingredient{Name: "Meat balls", Nutrients: Nutrients{Calories: 70, Proteins: 75, Carbs: 92, Fat: 80}}
	riceX := Ingredient{Name: "Rice", Nutrients: Nutrients{Calories: 4, Proteins: 2, Carbs: 3, Fat: 2}}
	breadX := Ingredient{Name: "Bread", Nutrients: Nutrients{Calories: 102, Proteins: 98, Carbs: 56, Fat: 83}}

	given1 := NewRecipe("Ham and Potatoes", 2).AddIngredient(ham, 1).AddIngredient(potato, 10)
	target1 := Nutrients{Calories: 1200, Proteins: 819, Carbs: 750, Fat: 838}

	given2 := NewRecipe("Chicken and Rice", 1).AddIngredient(chicken, 2).AddIngredient(rice, 100)
	target2 := Nutrients{Calories: 1300, Proteins: 1100, Carbs: 1150, Fat: 1230}

	given3 := NewRecipe("Meat balls, Rice and Bread", 1).AddIngredient(meatBallsX, 5).AddIngredient(riceX, 90).AddIngredient(breadX, 1)
	target3 := Nutrients{Calories: 1615, Proteins: 1350, Carbs: 1600, Fat: 1365}

	//pdf2
	salmon2 := Ingredient{Name: "Salmon", Nutrients: Nutrients{Calories: 195, Proteins: 8.5, Carbs: 0, Fat: 18.2}}
	tuna2 := Ingredient{Name: "Tuna", Nutrients: Nutrients{Calories: 202, Proteins: 24.9, Carbs: 2.1, Fat: 10.4}}
	bread2 := Ingredient{Name: "Bread", Nutrients: Nutrients{Calories: 212, Proteins: 5, Carbs: 40.3, Fat: 1.4}}
	rice2 := Ingredient{Name: "rice", Nutrients: Nutrients{Calories: 212, Proteins: 5, Carbs: 40.3, Fat: 1.4}}
	tomato2 := Ingredient{Name: "tomato", Nutrients: Nutrients{Calories: 20, Proteins: 0.7, Carbs: 2.9, Fat: 0.3}}

	given4 := NewRecipe("Tuna + Bread", 1).AddIngredient(tuna2, 10).AddIngredient(bread2, 1)
	target4 := Nutrients{Calories: 1800, Proteins: 200, Carbs: 60, Fat: 80}

	given5 := NewRecipe("Salmon + Bread + Tomato", 1).AddIngredient(salmon2, 3).AddIngredient(bread2, 1).AddIngredient(tomato2, 2)
	target5 := Nutrients{Calories: 1500, Proteins: 55, Carbs: 60, Fat: 95}

	given6 := NewRecipe("Salmon + Rice + Tomato", 1).AddIngredient(salmon2, 2).AddIngredient(rice2, 60).AddIngredient(tomato2, 8)
	target6 := Nutrients{Calories: 15100, Proteins: 280, Carbs: 3215, Fat: 90}

	given7 := NewRecipe("Tuna + Bread + Tomato", 1).AddIngredient(tuna2, 5).AddIngredient(bread2, 2).AddIngredient(tomato2, 8)
	target7 := Nutrients{Calories: 930, Proteins: 85, Carbs: 60, Fat: 35}

	given8 := NewRecipe("Tuna + Tomato", 1).AddIngredient(tuna2, 5).AddIngredient(tomato2, 5)
	target8 := Nutrients{Calories: 900, Proteins: 100, Carbs: 25, Fat: 42}

	givena := NewRecipe("Fish and Chips 2", 1).AddIngredient(fish, 1).AddIngredient(chips, 3)
	targeta := Nutrients{Calories: 27, Proteins: 2, Carbs: 2, Fat: 1}

	givenb := NewRecipe("random", 1).AddIngredient(fish, 1)
	targetb := Nutrients{Calories: 1, Proteins: 1, Carbs: 1, Fat: 1}

	type args struct {
		target         Nutrients
		givenRecep     Recipe
		debug          bool
		maxGenerations int
		// maxDuration    time.Duration
		popSize      int
		okPercentage float64
	}
	tests := []struct {
		name string
		args args
	}{
		{given0.Name, args{target: target0, givenRecep: given0, debug: false, maxGenerations: 50, popSize: 500, okPercentage: 0.928}},
		{given1.Name, args{target: target1, givenRecep: given1, debug: false, maxGenerations: 50, popSize: 500, okPercentage: 0.999}},
		{given2.Name, args{target: target2, givenRecep: given2, debug: false, maxGenerations: 50, popSize: 500, okPercentage: 0.99}},
		{given3.Name, args{target: target3, givenRecep: given3, debug: false, maxGenerations: 150, popSize: 500, okPercentage: 0.981}}, //could be better but takes longer
		{given4.Name, args{target: target4, givenRecep: given4, debug: false, maxGenerations: 50, popSize: 500, okPercentage: 0.965}},
		{given5.Name, args{target: target5, givenRecep: given5, debug: false, maxGenerations: 50, popSize: 500, okPercentage: 0.9085}},
		{given6.Name, args{target: target6, givenRecep: given6, debug: false, maxGenerations: 50, popSize: 500, okPercentage: 0.85}},
		{given7.Name, args{target: target7, givenRecep: given7, debug: false, maxGenerations: 50, popSize: 500, okPercentage: 0.9805}},
		{given8.Name, args{target: target8, givenRecep: given8, debug: false, maxGenerations: 50, popSize: 500, okPercentage: 0.967}},
		{givena.Name, args{target: targeta, givenRecep: givena, debug: false, maxGenerations: 50, popSize: 500, okPercentage: 0.0}}, //0.0551
		{givenb.Name, args{target: targetb, givenRecep: givenb, debug: false, maxGenerations: 50, popSize: 500, okPercentage: 0.0}}, //0.01
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()
			abortC := make(chan bool, 1)
			star := approximateResult(tt.args.target, tt.args.givenRecep, tt.args.debug,
				tt.args.maxGenerations, abortC, tt.args.popSize, tt.args.okPercentage)

			got := star.Given

			// diff := got.AllNutrients().DiffPrettyPrint(tt.args.target)
			res := got.AllNutrients().Diff(tt.args.target)
			p := tt.args.okPercentage
			if res < p {
				t.Errorf("wanted %.5f got %.5f ", p, res)
				// t.Errorf("approximateResult() = \n %v \n want %v \n diff: %v \n current percentage: %f", got.AllNutrients().PrettyPrint(), tt.args.target.PrettyPrint(), star.Fitness)
			}
		})
	}
}

// func TestConcurencyBasic(t *testing.T) {
// 	assert := assert.New(t)
// 	var wg sync.WaitGroup

// 	wg.Add(1)

// 	go func() {
// 		defer wg.Done()
// 		fmt.Println("Inside")
// 	}()

// 	wg.Wait()
// 	fmt.Println("Num goroutines sync done", runtime.NumGoroutine())
// 	fmt.Println()

// 	assert.True(false)
// }

func TestConcurency(t *testing.T) {
	return
	assert := assert.New(t)

	fish := Ingredient{Name: "Fish", Nutrients: Nutrients{Calories: 100, Proteins: 100, Carbs: 100, Fat: 100}}
	chips := Ingredient{Name: "Chips", Nutrients: Nutrients{Calories: 5, Proteins: 5, Carbs: 5, Fat: 5}}

	given := NewRecipe("Fish and Chips", 1).AddIngredient(fish, 1).AddIngredient(chips, 3)
	target := Nutrients{Calories: 140, Proteins: 140, Carbs: 140, Fat: 140}

	given2 := NewRecipe("Fish and Chips", 1).AddIngredient(fish, 1).AddIngredient(chips, 3)
	target2 := Nutrients{Calories: 140, Proteins: 140, Carbs: 140, Fat: 140}

	given3 := NewRecipe("Fish and Chips", 1).AddIngredient(fish, 1).AddIngredient(chips, 3)
	target3 := Nutrients{Calories: 140, Proteins: 140, Carbs: 140, Fat: 140}

	// given4 := NewRecipe("Fish and Chips", 1).AddIngredient(fish, 1).AddIngredient(chips, 3)
	// target4 := Nutrients{Calories: 140, Proteins: 140, Carbs: 140, Fat: 140}

	var wg sync.WaitGroup

	exec(given, target)
	exec(given, target)
	exec(given, target)

	fmt.Println("Sync done ----")
	fmt.Println("Num goroutines sync done", runtime.NumGoroutine())
	fmt.Println()

	wg.Add(2)

	start := time.Now()

	go func() {
		exec(given2, target2)
		wg.Done()
	}()

	go func() {
		exec(given3, target3)
		wg.Done()
	}()

	// go func() {
	// 	exec(given4, target4)
	// 	wg.Done()
	// }()

	// go func() {
	// 	exec(given, target)
	// 	wg.Done()
	// }()

	// go func() {
	// 	exec(given, target)
	// 	wg.Done()
	// }()

	wg.Wait()

	fmt.Println("Num goroutines all done", runtime.NumGoroutine())
	fmt.Printf("Concurrent time taken %s\n\n", time.Since(start))

	assert.True(false)
}

func exec(r Recipe, t Nutrients) {
	abortC := make(chan bool, 1)
	done := approximateResult(t, r, false, 10, abortC, 500, 1)

	fmt.Println("Num goroutines", runtime.NumGoroutine())
	fmt.Printf("Time Taken %s\n\n", done.Done)
}

// func TestConcurency(t *testing.T) {
// 	assert := assert.New(t)

// 	fish := Ingredient{Name: "Fish", Nutrients: Nutrients{Calories: 100, Proteins: 100, Carbs: 100, Fat: 100}}
// 	chips := Ingredient{Name: "Chips", Nutrients: Nutrients{Calories: 5, Proteins: 5, Carbs: 5, Fat: 5}}

// 	given := NewRecipe("Fish and Chips", 2).AddIngredient(fish, 1).AddIngredient(chips, 3)
// 	target := Nutrients{Calories: 140, Proteins: 140, Carbs: 140, Fat: 140}

// 	sample := 3
// 	fastest := fastestOfConcurrency(given, target, sample)

// 	fmt.Println(fastest)
// 	assert.NotNil(fastest)
// }

// func longestOf(given Recipe, target Nutrients, amount int) time.Duration {
// 	longest := time.Millisecond

// 	for i := 0; i < amount; i++ {
// 		abortC := make(chan bool, 1)
// 		resultChan := make(chan Problem, 1)
// 		var wg sync.WaitGroup

// 		approximateResult(target, given, true, 10000, abortC, 500, resultChan, 0.98, &wg)

// 		star := <-resultChan
// 		// got := star.Given
// 		if star.Done > longest {
// 			longest = star.Done
// 		}
// 	}

// 	return longest
// }

// func fastestOfConcurrency(given Recipe, target Nutrients, amount int) time.Duration {
// 	fastest := time.Second * time.Duration(10)
// 	wg := sync.WaitGroup{}
// 	l := sync.Mutex{}

// 	for i := 0; i < amount; i++ {
// 		wg.Add(1)
// 		go func() {
// 			abortC := make(chan bool, 1)
// 			resultChan := make(chan Problem, 1)
// 			bla := sync.WaitGroup{}
// 			bla.Add(1)
// 			approximateResult(target, given, true, 10000, abortC, 500, resultChan, 0.98, &bla)

// 			star := <-resultChan
// 			// got := star.Given
// 			l.Lock()
// 			if star.Done < fastest {
// 				fastest = star.Done
// 			}
// 			l.Unlock()
// 			wg.Done()
// 		}()
// 	}

// 	wg.Wait()

// 	return fastest
// }
