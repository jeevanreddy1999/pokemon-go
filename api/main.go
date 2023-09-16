package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"pokemon-gin/types"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

func getPokemon(c *gin.Context) {
	min := 0
	max := 1000
	length, err := strconv.Atoi(c.Query("length"))
	if err != nil {
		fmt.Print(err)
	}
	randArr := []int{}
	for i := 0; i <= length; i++ {
		val := rand.Intn(max-min) + min
		randArr = append(randArr, val)
	}
	ch := make(chan types.Pokemon)
	pokeArr := []types.Pokemon{}
	wg := sync.WaitGroup{}
	for i := 1; i < len(randArr); i++ {
		wg.Add(1)
		go getPokemonByID(randArr[i], &wg, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for res := range ch {
		pokeArr = append(pokeArr, res)
	}

	c.IndentedJSON(http.StatusOK, pokeArr)
}

func getPokemonByID(i int, wg *sync.WaitGroup, ch chan<- types.Pokemon) {
	defer wg.Done()
	resp, err1 := http.Get("https://pokeapi.co/api/v2/pokemon/" + strconv.Itoa(i))
	if err1 != nil {
		fmt.Println(err1)
	}
	defer resp.Body.Close()

	bodyBytes, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Print(err2)
	}

	var poke types.Pokemon
	err3 := json.Unmarshal(bodyBytes, &poke)
	if err3 != nil {
		fmt.Print(err3)
	}
	ch <- poke
}

func routes(router *gin.RouterGroup) {
	router.GET("/pokemon", getPokemon)
}

func init() {
	app = gin.New()
	app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Use /api/pokemon?length=5 for getting 5 random pokemon")
	})
	router := app.Group("/api")
	routes(router)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	gin.SetMode(gin.ReleaseMode)
	app.ServeHTTP(w, r)
}
