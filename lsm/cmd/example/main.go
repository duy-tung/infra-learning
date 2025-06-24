package main

import (
	"fmt"
	"math/rand"
	"time"

	"lsm/lsmtree"
)

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	tree, err := lsmtree.New("data", 1000)
	if err != nil {
		panic(err)
	}

	var sampleKey string
	for i := 0; i < 5000; i++ {
		key := fmt.Sprintf("key%05d_%s", i, randomString(5))
		val := randomString(20)
		if i == 2000 {
			sampleKey = key
		}
		if err := tree.Put(key, val); err != nil {
			panic(err)
		}
	}

	value, ok, err := tree.Get(sampleKey)
	if err != nil {
		panic(err)
	}
	fmt.Println("read existing?", ok, value)
}
