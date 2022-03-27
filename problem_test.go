package main

import (
	"testing"
	"time"
)

func TestProblem_isValid(t *testing.T) {
	fish := Ingredient{Name: "Fish", Nutrients: Nutrients{Calories: 100, Proteins: 100, Carbs: 100, Fat: 100}}
	chips := Ingredient{Name: "Chips", Nutrients: Nutrients{Calories: 5, Proteins: 5, Carbs: 5, Fat: 5}}
	tom := Ingredient{Name: "Tomato", Nutrients: Nutrients{Calories: 5, Proteins: 5, Carbs: 5, Fat: 5}}

	type fields struct {
		Org     Recipe
		Given   Recipe
		Target  Nutrients
		Fitness float64
		Done    time.Duration
	}

	f1 := fields{Org: NewRecipe("test", 1).AddIngredient(fish, 1).AddIngredient(chips, 3).AddIngredient(tom, 5)}

	type args struct {
		scale  float64
		recipe Recipe
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"test1", f1, args{1, NewRecipe("test", 1).AI(fish, 2).AI(chips, 3).AI(tom, 5)}, true},
		{"test2", f1, args{1, NewRecipe("test", 1).AI(fish, 1).AI(chips, 5).AI(tom, 4)}, true},
		{"test3", f1, args{1, NewRecipe("test", 1).AI(fish, 3).AI(chips, 9).AI(tom, 15)}, true},
		{"test4", f1, args{1, NewRecipe("test", 1).AI(fish, 2).AI(chips, 6).AI(tom, 5)}, true},
		{"test5", f1, args{1, NewRecipe("test", 1).AI(fish, 2).AI(chips, 2).AI(tom, 5)}, false},
		{"test6", f1, args{1, NewRecipe("test", 1).AI(fish, 1).AI(chips, 3).AI(tom, 0)}, false},
		{"test7", f1, args{1, NewRecipe("test", 1).AI(fish, 1).AI(chips, 2).AI(tom, 10)}, false},
		{"test8", f1, args{1, NewRecipe("test", 1).AI(fish, 10).AI(chips, 3).AI(tom, 5)}, false},
		{"test9", f1, args{1, NewRecipe("test", 1).AI(fish, 3).AI(chips, 9).AI(tom, 5)}, false},
		{"test10", f1, args{1, NewRecipe("test", 1).AI(fish, 0).AI(chips, 0).AI(tom, 5)}, false},
		{"test11", f1, args{1, NewRecipe("test", 1).AI(fish, 13).AI(chips, 3).AI(tom, 15)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Problem{
				Org:     tt.fields.Org,
				Given:   tt.fields.Given,
				Target:  tt.fields.Target,
				Fitness: tt.fields.Fitness,
				Done:    tt.fields.Done,
			}
			p := NewProblem(tt.args.recipe, tt.fields.Org, tt.fields.Target)

			if got := d.isValid(p); got != tt.want {
				t.Errorf("Problem.isValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
