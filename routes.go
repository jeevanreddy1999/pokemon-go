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
	"time"

	"github.com/gin-gonic/gin"
)

func rateLimit(c *gin.Context) {
	ip := c.ClientIP()
	value := int(ips.Add(ip, 1))
	if value%50 == 0 {
		fmt.Printf("ip: %s, count: %d\n", ip, value)
	}
	if value >= 200 {
		if value%200 == 0 {
			fmt.Println("ip blocked")
		}
		c.Abort()
		c.String(http.StatusServiceUnavailable, "you were automatically banned :)")
	}
}

func getPokemon(c *gin.Context) {
	fmt.Println(time.Now())
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
	fmt.Println(time.Now())
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
