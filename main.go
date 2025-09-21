package main

import (
	"context"

	"github.com/niluan304/auto-sign/tieba"
)

func main() {
	ctx := context.Background()

	must(tieba.Sign(ctx))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
