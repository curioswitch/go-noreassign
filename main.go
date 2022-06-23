package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/anuraaga/reassign/internal/analyzer"
)

func main() {
	singlechecker.Main(analyzer.New())
}
