package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	var port string
	flag.StringVar(&port, "port", "80", "port for the service to listen to")
	flag.Parse()

	// Echo instance
	e := echo.New()
	e.HideBanner = true

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.Any("/", handle)
	//these are preset tests
	e.Any("/test0", test0)
	e.Any("/test1", test1)
	e.Any("/test2", test2)
	e.Any("/test3", test3)
	e.Any("/test4", test4)
	e.Any("/test5", test5)
	e.Any("/test6", test6)
	e.Any("/test7", test7)
	e.Any("/test8", test8)
	e.Any("/testa", testa)
	e.Any("/testb", testb)
	e.Any("/testc", testc)

	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}

func StartViaApi(req Request, winnerChan chan Problem) {
	StartGenetic(convertHttpRecipe(req.Recipe), convertHttpTarget(req.Target), req.HttpConfig.OkPercentage, req.HttpConfig.TimeoutSeconds, winnerChan, 0)
}

func StartGenetic(r Recipe, t Nutrients, okPercentage float64, sleepSeconds int, winnerChan chan Problem, concurrency int) {
	rand.Seed(time.Now().UTC().UnixNano())
	debug := true
	convertedRec := r
	convertedTarget := t

	if concurrency == 0 { //auto
		concurrency = 1
	}
	// concurrency := 1
	abortChan := make(chan bool, 1)
	resultChan := make(chan Problem, concurrency)
	populationSize := 1000

	done := make(chan bool, 1)
	wg := sync.WaitGroup{}

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			resultChan <- approximateResult(convertedTarget, convertedRec, debug, 10000, abortChan, populationSize, okPercentage)
		}()
	}

	winner := Problem{Fitness: 0}

	go func() {
		for v := range resultChan {
			if v.Fitness > winner.Fitness {
				winner = v.DeepClone()
			}
			if winner.Fitness == 1 || winner.Fitness >= okPercentage {
				fmt.Println(v.Done)
				close(done)
				winnerChan <- winner
				return
			}
		}

		close(done)
		winnerChan <- winner
	}()

	select {
	case <-time.After(time.Duration(sleepSeconds) * time.Second):
		fmt.Println("Timeout")
		close(abortChan)
		wg.Wait()
		close(resultChan)
	case <-done:
		fmt.Println("Done")
		close(abortChan)
	}
}

// Handler
func handle(c echo.Context) error {
	decodeJson := json.NewDecoder(c.Request().Body)

	var req Request
	err := decodeJson.Decode(&req)

	if err != nil {
		return c.JSON(http.StatusBadRequest, "Could not parse request: "+err.Error())
	}

	// start := time.Now()
	result := make(chan Problem, 0)
	StartViaApi(req, result)
	winner := <-result
	// elapsed := time.Since(start)

	return getResultManual(c, req.HttpConfig.TimeoutSeconds, 0, req.HttpConfig.OkPercentage,
		req.Recipe.Scale, convertHttpRecipe(req.Recipe), convertHttpTarget(req.Target), winner)
}

func convertRecipeToHttpRecipe(r Recipe) HttpRecipe {
	recipe := HttpRecipe{Scale: r.Scale, Name: r.Name}

	in := map[string]HttpIngredient{}

	for _, v := range r.Ingredients {
		nut := v.Nutrients
		value, found := in[v.Name]
		if !found {
			in[v.Name] = HttpIngredient{Calories: nut.Calories, Carbs: nut.Carbs, Fat: nut.Fat, Protein: nut.Proteins, Amount: 1, Name: v.Name}
		} else {
			value.Amount++
			in[v.Name] = value
		}
	}

	for _, v := range in {
		recipe.Ingredients = append(recipe.Ingredients, v)
	}

	return recipe

}
func convertHttpRecipe(hr HttpRecipe) Recipe {
	r := NewRecipe(hr.Name, hr.Scale)

	for _, v := range hr.Ingredients {
		r = r.AddIngredient(Ingredient{Name: v.Name, Nutrients: Nutrients{Calories: v.Calories, Proteins: v.Protein, Carbs: v.Carbs, Fat: v.Fat}}, v.Amount)
	}

	return r

}
func convertHttpTarget(v HttpTarget) Nutrients {
	return Nutrients{Calories: v.Calories, Proteins: v.Protein, Carbs: v.Carbs, Fat: v.Fat}
}

func parseQueryParams(c echo.Context, defaultTimeout, defaultConcurrency int, defaultPercentage, defaultScale float64) (timeout int, concurrency int, percentage float64, scale float64) {
	t := c.QueryParam("timeout_seconds")
	p := c.QueryParam("percentage")
	s := c.QueryParam("scale")
	con := c.QueryParam("concurrency")

	concurrency, err := strconv.Atoi(con)
	if err != nil {
		concurrency = defaultConcurrency
	}

	timeout, err = strconv.Atoi(t)
	if err != nil {
		timeout = defaultTimeout
	}

	percentage, err = strconv.ParseFloat(p, 64)
	if err != nil {
		percentage = defaultPercentage
	}

	scale, err = strconv.ParseFloat(s, 64)
	if err != nil {
		scale = defaultScale
	}

	return
}

func parse(c echo.Context) (timeout, concurrency int, okPercentage float64, scale float64) {
	return parseQueryParams(c, 3, 1, 0.99, 1)
}

type Output struct {
	ParametersUsed      string
	Timetaken           string
	Target              Nutrients
	WinnerNutrients     Nutrients
	DifferenceNutrients Nutrients
	DiffSum             float64
	DiffPercentage      float64
	WinnerScale         string
	Winner              HttpRecipe
	GivenRecipe         HttpRecipe
}

func getResult(c echo.Context, given Recipe, target Nutrients, winner Problem) error {
	timeout, concurrency, percentage, scale := parse(c)

	return getResultManual(c, timeout, concurrency, percentage, scale, given, target, winner)
}

func getResultManual(c echo.Context, timeout, concurrency int, percentage, scale float64, given Recipe, target Nutrients, winner Problem) error {
	winnerNut := winner.Given.AllNutrients()
	diffNut, diffSum, diffPer := winnerNut.DiffPrettyOut(winner.Target)

	givenR := convertRecipeToHttpRecipe(winner.Org)

	outScale := winner.Org.GetScale(winner.Given)
	build := ""
	i := 0
	for index, v := range outScale {
		if i+1 == len(outScale) {
			build += fmt.Sprintf("%s: %.2f", index, v)
		} else {
			build += fmt.Sprintf("%s: %.2f  |  ", index, v)
		}
		i++
	}

	out := Output{
		WinnerScale:         build,
		GivenRecipe:         givenR,
		Winner:              convertRecipeToHttpRecipe(winner.Given),
		ParametersUsed:      fmt.Sprintf("timeout_seconds %d  |  percentage %.4f  |  scale %.2f  |  concurrency %d", timeout, percentage, scale, concurrency),
		Target:              target,
		Timetaken:           fmt.Sprintf("%s", winner.Done),
		WinnerNutrients:     winnerNut,
		DifferenceNutrients: diffNut,
		DiffSum:             diffSum,
		DiffPercentage:      diffPer}

	return c.JSON(http.StatusOK, out)
}

func testCase(given Recipe, target Nutrients, c echo.Context) Problem {
	timeout, concurrency, percentage, scale := parse(c)

	given.Scale = scale

	start := time.Now()
	result := make(chan Problem, 0)
	StartGenetic(given, target, percentage, timeout, result, concurrency)
	winner := <-result
	elapsed := time.Since(start)
	winner.Done = elapsed

	return winner
}

type Request struct {
	Recipe     HttpRecipe `json:"recipe" form:"recipe" query:"recipe"`
	Target     HttpTarget `json:"target" form:"target" query:"target"`
	HttpConfig HttpConfig `json:"config" form:"config" query:"config"`
}

type HttpRecipe struct {
	Ingredients []HttpIngredient `json:"ingredients" form:"ingredients" query:"ingredients"`
	Scale       float64          `json:"scale" form:"scale" query:"scale"`
	Name        string           `json:"name" form:"name" query:"name"`
}

type HttpIngredient struct {
	Calories float64 `json:"calories" form:"calories" query:"calories"`
	Protein  float64 `json:"protein" form:"protein" query:"protein"`
	Carbs    float64 `json:"carbs" form:"carbs" query:"carbs"`
	Fat      float64 `json:"fat" form:"fat" query:"fat"`
	Amount   int     `json:"amount" form:"amount" query:"amount"`
	Name     string  `json:"name" form:"name" query:"name"`
}

type HttpTarget struct {
	Calories float64 `json:"calories" form:"calories" query:"calories"`
	Protein  float64 `json:"protein" form:"protein" query:"protein"`
	Carbs    float64 `json:"carbs" form:"carbs" query:"carbs"`
	Fat      float64 `json:"fat" form:"fat" query:"fat"`
}

type HttpConfig struct {
	TimeoutSeconds int     `json:"timeout_seconds" form:"timeout_seconds" query:"timeout_seconds"`
	OkPercentage   float64 `json:"percentage" form:"percentage" query:"percentage"`
}
