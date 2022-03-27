package main

import (
	"fmt"
	"math"
)

type Recipe struct {
	Ingredients []Ingredient
	Scale       float64
	Name        string
	counted     map[string]int64
}

func (r Recipe) GetScale(newRecipe Recipe) map[string]float64 {
	countOrg := r.CountIngredients()
	countNew := newRecipe.CountIngredients()

	propMap := map[string]float64{}
	for i, v := range countOrg {
		value, found := countNew[i]
		if !found || value == 0 {
			for k := range propMap {
				propMap[k] = -1
			}
			return propMap
		}

		propMap[i] = float64(value) / float64(v)
	}

	return propMap
}

func (r Recipe) CanAdapt(newRecipe Recipe) bool {
	countOrg := r.CountIngredients()
	countNew := newRecipe.CountIngredients()

	scale := r.Scale

	prop := make([]float64, len(countOrg))
	index := 0
	for i, v := range countOrg {
		value, found := countNew[i]
		if !found || value == 0 {
			// fmt.Println("Not found or value is 0", found, value)
			return false
		}

		prop[index] = float64(value) / float64(v)
		index++
	}

	for i := range prop {
		for j := i + 1; j < len(prop); j++ {
			val := math.Abs(prop[i] - prop[j])
			if val > scale {
				return false
			}
		}
	}

	return true
}

func (r *Recipe) CountIngredients() map[string]int64 {
	counts := map[string]int64{}
	for _, v := range r.Ingredients {
		counts[v.Name]++
	}

	r.counted = counts

	return counts
}

func (r Recipe) FlatIngredients() (map[string]int64, map[string]Ingredient) {
	counts := map[string]int64{}
	ingre := map[string]Ingredient{}
	for _, v := range r.Ingredients {
		counts[v.Name]++
		ingre[v.Name] = v
	}

	return counts, ingre
}

func (r Recipe) PrettyPrint() string {
	res := fmt.Sprintf("Recipe Name: %s\n", r.Name)
	res += fmt.Sprintf("Ingredients: \n")
	countMap := map[Ingredient]int64{}

	for _, v := range r.Ingredients {
		countMap[v]++
		// res += fmt.Sprintf("%s | %s \n", v.Name, v.Nutrients.PrettyPrint())
	}
	for i, v := range countMap {
		res += fmt.Sprintf("%d %s | %s \n", v, i.Name, i.Nutrients.PrettyPrint())
	}

	res += fmt.Sprintf("Recipe Nutrients: %s \n", r.AllNutrients().PrettyPrint())
	return res
}

func (r Recipe) PrettyPrintNames() string {
	res := ""
	countMap := map[Ingredient]int64{}

	for _, v := range r.Ingredients {
		countMap[v]++
		// res += fmt.Sprintf("%s | %s \n", v.Name, v.Nutrients.PrettyPrint())
	}
	for i, v := range countMap {
		res += fmt.Sprintf("%d %s | ", v, i.Name)
	}

	return res
}

func (r Recipe) DeepClone() Recipe {
	new := Recipe{}
	new.Name = r.Name
	new.Scale = r.Scale
	new.Ingredients = make([]Ingredient, len(r.Ingredients))
	copy(new.Ingredients, r.Ingredients)
	new.CountIngredients()
	return new
}

func NewRecipe(name string, scale float64) Recipe {
	return Recipe{Name: name, Scale: scale, Ingredients: []Ingredient{}}
}

func (r Recipe) AddIngredient(ing Ingredient, times int) Recipe {
	clone := r.DeepClone()

	for i := 0; i < times; i++ {
		clone.Ingredients = append(clone.Ingredients, ing)
	}

	clone.CountIngredients()
	return clone
}

func (r Recipe) AI(ing Ingredient, times int) Recipe {
	return r.AddIngredient(ing, times)
}

// func (r Recipe) FromIngredients(ing []Ingredient) Recipe {
// 	new := r.DeepClone()
// 	new.Ingredients = ing
// 	return new
// }

func (r Recipe) AllNutrients() Nutrients {
	var carbs float64
	var protein float64
	var fat float64
	var calories float64

	for _, v := range r.Ingredients {
		carbs += v.Nutrients.Carbs
		protein += v.Nutrients.Proteins
		fat += v.Nutrients.Fat
		calories += v.Nutrients.Calories
	}

	return Nutrients{carbs, protein, fat, calories}
}
