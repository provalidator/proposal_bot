package main_test

import (
	"fmt"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	type ti struct {
		a time.Time
		b time.Time
	}
	var Ti []ti

	am := time.Now().UTC()
	bm := time.Now().Local()
	fmt.Println(am, bm)
	Ti = append(Ti, ti{am, bm})
	fmt.Println(Ti[0].a, Ti[0].b)
}
