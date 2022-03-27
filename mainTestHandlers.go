package main

import (
	"github.com/labstack/echo"
)

func test0(c echo.Context) error {
	fish := Ingredient{Name: "Fish", Nutrients: Nutrients{Calories: 100, Proteins: 100, Carbs: 100, Fat: 100}}
	chips := Ingredient{Name: "Chips", Nutrients: Nutrients{Calories: 5, Proteins: 5, Carbs: 5, Fat: 5}}
	given := NewRecipe("Fish and Chips", 2).AddIngredient(fish, 1).AddIngredient(chips, 3)
	target := Nutrients{Calories: 140, Proteins: 140, Carbs: 140, Fat: 140}

	winner := testCase(given, target, c)
	return getResult(c, given, target, winner)
}

func test1(c echo.Context) error {
	ham := Ingredient{Name: "Ham", Nutrients: Nutrients{Calories: 500, Proteins: 400, Carbs: 400, Fat: 600}}
	potato := Ingredient{Name: "Potato", Nutrients: Nutrients{Calories: 100, Proteins: 59.9, Carbs: 50, Fat: 34}}
	given := NewRecipe("Ham and Potatoes", 2).AddIngredient(ham, 1).AddIngredient(potato, 10)
	target := Nutrients{Calories: 1200, Proteins: 819, Carbs: 750, Fat: 838}

	winner := testCase(given, target, c)
	return getResult(c, given, target, winner)
}

func test2(c echo.Context) error {
	chicken := Ingredient{Name: "Chicken", Nutrients: Nutrients{Calories: 300, Proteins: 200, Carbs: 254, Fat: 323}}
	rice := Ingredient{Name: "Rice", Nutrients: Nutrients{Calories: 3, Proteins: 4, Carbs: 3, Fat: 2}}
	given := NewRecipe("Chicken and Rice", 2).AddIngredient(chicken, 2).AddIngredient(rice, 100)
	target := Nutrients{Calories: 1300, Proteins: 1100, Carbs: 1150, Fat: 1230}

	winner := testCase(given, target, c)
	return getResult(c, given, target, winner)
}

func test3(c echo.Context) error {
	meatBallsX := Ingredient{Name: "Meat balls", Nutrients: Nutrients{Calories: 70, Proteins: 75, Carbs: 92, Fat: 80}}
	riceX := Ingredient{Name: "Rice", Nutrients: Nutrients{Calories: 4, Proteins: 2, Carbs: 3, Fat: 2}}
	breadX := Ingredient{Name: "Bread", Nutrients: Nutrients{Calories: 102, Proteins: 98, Carbs: 56, Fat: 83}}
	given := NewRecipe("Meat balls, Rice and Bread", 1).AddIngredient(meatBallsX, 5).AddIngredient(riceX, 90).AddIngredient(breadX, 1)
	target := Nutrients{Calories: 1615, Proteins: 1350, Carbs: 1600, Fat: 1365}

	winner := testCase(given, target, c)
	return getResult(c, given, target, winner)
}

func test4(c echo.Context) error {
	tuna2 := Ingredient{Name: "Tuna", Nutrients: Nutrients{Calories: 202, Proteins: 24.9, Carbs: 2.1, Fat: 10.4}}
	bread2 := Ingredient{Name: "Bread", Nutrients: Nutrients{Calories: 212, Proteins: 5, Carbs: 40.3, Fat: 1.4}}
	given := NewRecipe("Tuna + Bread", 1).AddIngredient(tuna2, 10).AddIngredient(bread2, 1)
	target := Nutrients{Calories: 1800, Proteins: 200, Carbs: 60, Fat: 80}

	winner := testCase(given, target, c)
	return getResult(c, given, target, winner)
}

func test5(c echo.Context) error {
	salmon2 := Ingredient{Name: "Salmon", Nutrients: Nutrients{Calories: 195, Proteins: 8.5, Carbs: 0, Fat: 18.2}}
	bread2 := Ingredient{Name: "Bread", Nutrients: Nutrients{Calories: 212, Proteins: 5, Carbs: 40.3, Fat: 1.4}}
	tomato2 := Ingredient{Name: "Tomato", Nutrients: Nutrients{Calories: 20, Proteins: 0.7, Carbs: 2.9, Fat: 0.3}}
	given := NewRecipe("Salmon + Bread + Tomato", 1).AddIngredient(salmon2, 3).AddIngredient(bread2, 1).AddIngredient(tomato2, 2)
	target := Nutrients{Calories: 1500, Proteins: 55, Carbs: 60, Fat: 95}

	winner := testCase(given, target, c)
	return getResult(c, given, target, winner)
}

func test6(c echo.Context) error {
	salmon2 := Ingredient{Name: "Salmon", Nutrients: Nutrients{Calories: 195, Proteins: 8.5, Carbs: 0, Fat: 18.2}}
	rice2 := Ingredient{Name: "Rice", Nutrients: Nutrients{Calories: 212, Proteins: 5, Carbs: 40.3, Fat: 1.4}}
	tomato2 := Ingredient{Name: "Tomato", Nutrients: Nutrients{Calories: 20, Proteins: 0.7, Carbs: 2.9, Fat: 0.3}}
	given := NewRecipe("Salmon + Rice + Tomato", 1).AddIngredient(salmon2, 2).AddIngredient(rice2, 60).AddIngredient(tomato2, 8)
	target := Nutrients{Calories: 15100, Proteins: 280, Carbs: 3215, Fat: 90}

	winner := testCase(given, target, c)
	return getResult(c, given, target, winner)
}

func test7(c echo.Context) error {
	tuna2 := Ingredient{Name: "Tuna", Nutrients: Nutrients{Calories: 202, Proteins: 24.9, Carbs: 2.1, Fat: 10.4}}
	bread2 := Ingredient{Name: "Bread", Nutrients: Nutrients{Calories: 212, Proteins: 5, Carbs: 40.3, Fat: 1.4}}
	tomato2 := Ingredient{Name: "Tomato", Nutrients: Nutrients{Calories: 20, Proteins: 0.7, Carbs: 2.9, Fat: 0.3}}
	given := NewRecipe("Tuna + Bread + Tomato", 1).AddIngredient(tuna2, 5).AddIngredient(bread2, 2).AddIngredient(tomato2, 8)
	target := Nutrients{Calories: 930, Proteins: 85, Carbs: 60, Fat: 35}

	winner := testCase(given, target, c)
	return getResult(c, given, target, winner)
}

func test8(c echo.Context) error {
	tuna2 := Ingredient{Name: "Tuna", Nutrients: Nutrients{Calories: 202, Proteins: 24.9, Carbs: 2.1, Fat: 10.4}}
	tomato2 := Ingredient{Name: "Tomato", Nutrients: Nutrients{Calories: 20, Proteins: 0.7, Carbs: 2.9, Fat: 0.3}}
	given := NewRecipe("Tuna + Tomato", 1).AddIngredient(tuna2, 5).AddIngredient(tomato2, 5)
	target := Nutrients{Calories: 900, Proteins: 100, Carbs: 25, Fat: 42}

	winner := testCase(given, target, c)
	return getResult(c, given, target, winner)
}

func testa(c echo.Context) error {
	fish := Ingredient{Name: "Fish", Nutrients: Nutrients{Calories: 140, Proteins: 140, Carbs: 140, Fat: 140}}
	chips := Ingredient{Name: "Chips", Nutrients: Nutrients{Calories: 5, Proteins: 5, Carbs: 5, Fat: 5}}
	given := NewRecipe("Fish and Chips", 1).AddIngredient(fish, 1).AddIngredient(chips, 3)
	target := Nutrients{Calories: 27, Proteins: 2, Carbs: 2, Fat: 1}

	winner := testCase(given, target, c)
	return getResult(c, given, target, winner)
}

func testb(c echo.Context) error {
	i2 := Ingredient{Name: "SomeIngredient", Nutrients: Nutrients{Calories: 1, Proteins: 2, Carbs: 3, Fat: 4}}
	given := NewRecipe("Recipe number 241", 1).AddIngredient(i2, 1)
	target := Nutrients{Calories: 1, Proteins: 1, Carbs: 1, Fat: 1}

	winner := testCase(given, target, c)
	return getResult(c, given, target, winner)
}

func testc(c echo.Context) error {
	ham := Ingredient{Name: "Ham", Nutrients: Nutrients{Calories: 500, Proteins: 400, Carbs: 400, Fat: 600}}
	potato := Ingredient{Name: "Potato", Nutrients: Nutrients{Calories: 100, Proteins: 59.9, Carbs: 50, Fat: 34}}
	chicken := Ingredient{Name: "Chicken", Nutrients: Nutrients{Calories: 300, Proteins: 200, Carbs: 254, Fat: 323}}
	rice := Ingredient{Name: "Rice", Nutrients: Nutrients{Calories: 3, Proteins: 4, Carbs: 3, Fat: 2}}
	meatBallsX := Ingredient{Name: "Meat balls", Nutrients: Nutrients{Calories: 70, Proteins: 75, Carbs: 92, Fat: 80}}
	riceX := Ingredient{Name: "Rice X", Nutrients: Nutrients{Calories: 4, Proteins: 2, Carbs: 3, Fat: 2}}
	breadX := Ingredient{Name: "Bread X", Nutrients: Nutrients{Calories: 102, Proteins: 98, Carbs: 56, Fat: 83}}
	salmon2 := Ingredient{Name: "Salmon", Nutrients: Nutrients{Calories: 195, Proteins: 8.5, Carbs: 0, Fat: 18.2}}
	tuna2 := Ingredient{Name: "Tuna 2", Nutrients: Nutrients{Calories: 202, Proteins: 24.9, Carbs: 2.1, Fat: 10.4}}
	bread2 := Ingredient{Name: "Bread 2", Nutrients: Nutrients{Calories: 212, Proteins: 5, Carbs: 40.3, Fat: 1.4}}
	rice2 := Ingredient{Name: "Rice 2", Nutrients: Nutrients{Calories: 212, Proteins: 5, Carbs: 40.3, Fat: 1.4}}
	tomato2 := Ingredient{Name: "tomato 2", Nutrients: Nutrients{Calories: 20, Proteins: 0.7, Carbs: 2.9, Fat: 0.3}}

	given := NewRecipe("Recipe number 241", 1).AI(ham, 3).AI(potato, 10).AI(chicken, 3).
		AI(rice, 70).AI(meatBallsX, 4).AI(riceX, 60).AI(breadX, 3).AI(salmon2, 5).AI(tuna2, 5).AI(bread2, 10).AI(rice2, 22).AI(tomato2, 9)

	target := Nutrients{Calories: 80000, Proteins: 17300, Carbs: 13333, Fat: 13000}

	winner := testCase(given, target, c)
	return getResult(c, given, target, winner)
}
