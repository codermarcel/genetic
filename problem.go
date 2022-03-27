package main

import (
	"math"
	"math/rand"
	"time"
)

type Problem struct {
	Org     Recipe
	Given   Recipe
	Target  Nutrients
	Fitness float64
	Done    time.Duration //time when it is done
}

func NewProblem(r Recipe, org Recipe, target Nutrients) Problem {
	var problem Problem
	problem.Org = org.DeepClone()
	clR := r.DeepClone()
	problem.Given = clR
	problem.Target = target
	problem.calcFitness()

	return problem
}

func (p Problem) Info() string {
	return p.Given.AllNutrients().DiffPrettyPrint(p.Target)
}

func (p Problem) DeepClone() Problem {
	var problem Problem
	problem.Org = p.Org.DeepClone()
	problem.Given = p.Given.DeepClone()
	problem.Target = p.Target
	problem.Fitness = p.Fitness
	problem.Done = p.Done

	return problem
}

// calculates the fitness of the Organism
func (d *Problem) calcFitness() {
	n := d.Given.AllNutrients()
	// proz := n.Pro(d.Target)
	proz := n.Fitness(d.Target)

	if math.IsInf(proz, 1) || math.IsInf(proz, -1) {
		d.Fitness = 0
		return
	}

	if !d.isValid(*d) {
		d.Fitness = 0
	} else {
		d.Fitness = proz
	}

}

func (d *Problem) isValid(newProblem Problem) bool {
	return d.Org.CanAdapt(newProblem.Given)
}

func (d Problem) buildFlatRecipe(counts map[string]int64, ingre map[string]Ingredient) Recipe {
	rec := NewRecipe(d.Given.Name, d.Given.Scale)
	for i, v := range ingre {
		rec = rec.AddIngredient(v, int(counts[i]))
	}
	return rec
}

func (d *Problem) canReplace(old string, new string) bool {
	counts, ingre := d.Given.FlatIngredients()
	counts[old]--
	counts[new]++
	re := d.buildFlatRecipe(counts, ingre)
	cloned := d.DeepClone()
	cloned.Given = re
	return d.isValid(cloned)
}

func (d *Problem) canAdd(name string) bool {
	counts, ingre := d.Given.FlatIngredients()
	counts[name]++
	re := d.buildFlatRecipe(counts, ingre)
	cloned := d.DeepClone()
	cloned.Given = re
	return d.isValid(cloned)
}

func (d *Problem) replace(i, new int32) {
	d.Given.Ingredients[i] = d.Org.Ingredients[new]
}

func (d *Problem) randomReplace(i int) {
	if i < 1 {
		return
	}

	choice := rand.Int31n(int32(len(d.Org.Ingredients)))
	repC := rand.Int31n(int32(len(d.Given.Ingredients)))
	o := d.Given.Ingredients[repC]
	n := d.Org.Ingredients[choice]

	if d.canReplace(o.Name, n.Name) {
		d.replace(repC, choice)
	} else {
		d.randomReplace(i - 1)
	}
}

func (d *Problem) add(i Ingredient) {
	d.Given.Ingredients = append(d.Given.Ingredients, i)
}

func (d *Problem) randomAdd(i int) {
	if i < 1 {
		return
	}

	choice := rand.Int31n(int32(len(d.Org.Ingredients)))

	if d.canAdd(d.Org.Ingredients[choice].Name) {
		d.Given.Ingredients = append(d.Given.Ingredients, d.Org.Ingredients[choice])
	} else {
		d.randomAdd(i - 1)
	}
}

func (d *Problem) mutate() {
	// return
	rate := 0.95
	rate2 := 0.95
	rate3 := 0.99

	if len(d.Given.CountIngredients()) == 0 {
		choice := rand.Int31n(int32(len(d.Org.Ingredients)))
		d.Given.Ingredients = append(d.Given.Ingredients, d.Org.Ingredients[choice])
		return
	}

	if rand.Float64() > rate {
		d.randomReplace(10)
	}

	if rand.Float64() > rate2 {
		d.randomAdd(10)
	}

	if rand.Float64() > rate3 {
		for i := 0; i < len(d.Org.Ingredients); i++ {
			d.randomAdd(3)
		}
	}
}

// crosses over 2 Organisms
func (p Problem) crossover(d2 Problem) Problem {
	return p.DeepClone()
	//new := p.DeepClone()
	//
	//new.Given.Ingredients = []Ingredient{}
	//
	//for i := 0; i < len(p.Given.Ingredients); i++ {
	//	if i%2 == 0 {
	//		if new.canAdd(p.Given.Ingredients[i].Name) {
	//			new.Given.Ingredients = append(new.Given.Ingredients, p.Given.Ingredients[i])
	//		}
	//	}
	//}
	//
	//for i := 0; i < len(d2.Given.Ingredients); i++ {
	//	if i%2 == 0 {
	//		if new.canAdd(d2.Given.Ingredients[i].Name) {
	//			new.Given.Ingredients = append(new.Given.Ingredients, d2.Given.Ingredients[i])
	//		}
	//	}
	//}
	//
	//return new
}
