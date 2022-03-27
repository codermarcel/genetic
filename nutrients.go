package main

import (
	"fmt"
	"math"
)

type Nutrients struct {
	Carbs    float64
	Proteins float64
	Fat      float64
	Calories float64
}

func (n Nutrients) PrettyPrint() string {
	return fmt.Sprintf("Calories %.2f Proteins %.2f Carbs %.2f Fat %.2f", n.Calories, n.Proteins, n.Carbs, n.Fat)
}

func (n Nutrients) DiffPrettyPrint(target Nutrients) string {
	d := n.DiffNutrients(target)
	diffSum := d.SumUp()
	diffPer := n.Diff(target)
	return fmt.Sprintf("Has %s \nTarget %s \ndifference %s\ndifference sum %.2f\ndifference perc %.4f\n", n.PrettyPrint(), target.PrettyPrint(), d.PrettyPrint(), diffSum, diffPer)
}

func (n Nutrients) DiffPrettyOut(target Nutrients) (difNut Nutrients, diffSum float64, diffPer float64) {
	d := n.DiffNutrients(target)
	diffSum = d.SumUp()
	diffPer = n.Diff(target)

	return d, diffSum, diffPer
}

func (n Nutrients) SumUp() float64 {
	return math.Abs(n.Calories) + math.Abs(n.Proteins) + math.Abs(n.Carbs) + math.Abs(n.Fat)
}

func (n Nutrients) DiffNutrients(target Nutrients) Nutrients {
	var carbs, proteins, fat, calories float64
	carbs = n.Carbs - target.Carbs
	proteins = n.Proteins - target.Proteins
	fat = n.Fat - target.Fat
	calories = n.Calories - target.Calories

	return Nutrients{Carbs: carbs, Proteins: proteins, Fat: fat, Calories: calories}
}

func calcDif(given, target float64) float64 {
	if given > target {
		return target / given
	}

	return given / target
}

func (given Nutrients) Diff(target Nutrients) float64 {
	carb := calcDif(given.Carbs, target.Carbs)
	cal := calcDif(given.Calories, target.Calories)

	fat := calcDif(given.Fat, target.Fat)
	prot := calcDif(given.Proteins, target.Proteins)

	all := carb + cal + fat + prot
	res := all / 4

	return res
}

func (given Nutrients) Fitness(target Nutrients) float64 {
	carb := calcDif(given.Carbs, target.Carbs)
	cal := calcDif(given.Calories, target.Calories)

	fat := calcDif(given.Fat, target.Fat)
	prot := calcDif(given.Proteins, target.Proteins)

	penaltyRatio := 0.3

	all := carb + cal + fat + prot
	res := all / 4

	if carb < penaltyRatio {
		res = res / 3
	}
	if cal < penaltyRatio {
		res = res / 3
	}
	if fat < penaltyRatio {
		res = res / 3
	}
	if prot < penaltyRatio {
		res = res / 3
	}

	return res
}

func (n Nutrients) Procentual(dif float64) float64 {
	return dif
}

func (n Nutrients) Pro(c Nutrients) float64 {
	return n.Procentual(n.Diff(c))
}
