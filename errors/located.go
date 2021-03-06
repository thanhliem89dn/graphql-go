package graphqlerrors

import (
	"github.com/chris-ramon/graphql-go/language/ast"
)

func NewLocatedError(err interface{}, nodes []ast.Node) *GraphQLError {
	message := "An unknown error occurred."
	if err, ok := err.(error); ok {
		message = err.Error()
	}
	if err, ok := err.(string); ok {
		message = err
	}
	stack := message
	return NewGraphQLError(
		message,
		nodes,
		stack,
		nil,
		[]int{},
	)
}

func FieldASTsToNodeASTs(fieldASTs []*ast.Field) []ast.Node {
	nodes := []ast.Node{}
	for _, fieldAST := range fieldASTs {
		nodes = append(nodes, fieldAST)
	}
	return nodes
}
