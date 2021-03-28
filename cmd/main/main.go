package main

import (
	"github.com/lgunko/beauty-organisation-service/cmd/start"
	"github.com/lgunko/beauty-reuse/env"
	"os"
)

func main() {
	os.Setenv("EXECUTION_ENVIRONMENT", env.Prod.String())
	start.Start()
}
