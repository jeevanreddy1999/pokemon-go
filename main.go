package main

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

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, "Hello")
	})
	router.GET("/pokemon", getPokemon)
	router.Run(":8080")
	fmt.Println("Server started on 8080 port")
}
