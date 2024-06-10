package main

import (
	"github.com/golesson/go-graceful-shutdown/tests"
)

func main() {

	//test with first way to achieve it
	// tests.With_native_libs()

	tests.With_go_chi()
}
