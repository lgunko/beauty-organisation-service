package start

import (
	"context"
	"github.com/lgunko/beauty-organisation-service/graph"
	"github.com/lgunko/beauty-organisation-service/graph/generated"
	"github.com/lgunko/beauty-reuse/env"
	"github.com/lgunko/beauty-reuse/graph/directives"
	"github.com/lgunko/beauty-reuse/server"
)

func Start() {
	env.SetUpEnv()

	database := server.GetDatabase()
	defer database.Client().Disconnect(context.Background())

	resolver := graph.NewResolver(database)
	server.StartServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver, Directives: generated.DirectiveRoot{
		ValidateLength: directives.ValidateLengthDirective(),
	}}))
}
