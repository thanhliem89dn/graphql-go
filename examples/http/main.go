package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chris-ramon/graphql-go"
	"github.com/chris-ramon/graphql-go/types"
)

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var data map[string]user

/*
   Create User object type with fields "id" and "name" by using GraphQLObjectTypeConfig:
       - Name: name of object type
       - Fields: a map of fields by using GraphQLFieldConfigMap
   Setup type of field use GraphQLFieldConfig
*/
var userType = types.NewGraphQLObjectType(
	types.GraphQLObjectTypeConfig{
		Name: "User",
		Fields: types.GraphQLFieldConfigMap{
			"id": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
	},
)

/*
   Create Query object type with fields "user" has type [userType] by using GraphQLObjectTypeConfig:
       - Name: name of object type
       - Fields: a map of fields by using GraphQLFieldConfigMap
   Setup type of field use GraphQLFieldConfig to define:
       - Type: type of field
       - Args: arguments to query with current field
       - Resolve: function to query data using params from [Args] and return value with current type
*/
var queryType = types.NewGraphQLObjectType(
	types.GraphQLObjectTypeConfig{
		Name: "Query",
		Fields: types.GraphQLFieldConfigMap{
			"user": &types.GraphQLFieldConfig{
				Type: userType,
				Args: types.GraphQLFieldConfigArgumentMap{
					"id": &types.GraphQLArgumentConfig{
						Type: types.GraphQLString,
					},
				},
				Resolve: func(p types.GQLFRParams) interface{} {
					idQuery, isOK := p.Args["id"].(string)
					if isOK {
						return data[idQuery]
					}
					return nil
				},
			},
		},
	})

var schema, _ = types.NewGraphQLSchema(
	types.GraphQLSchemaConfig{
		Query: queryType,
	},
)

func executeQuery(query string, schema types.GraphQLSchema) *types.GraphQLResult {
	graphqlParams := gql.GraphqlParams{
		Schema:        schema,
		RequestString: query,
	}
	resultChannel := make(chan *types.GraphQLResult)
	go gql.Graphql(graphqlParams, resultChannel)
	result := <-resultChannel
	if len(result.Errors) > 0 {
		fmt.Println("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func main() {
	_ = importJSONDataFromFile("data.json", &data)

	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query()["query"][0], schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Now server is running on port 8080")
	fmt.Println("Test with Get      : curl -g \"http://localhost:8080/graphql?query={user(id:%221%22){name}}\"")
	http.ListenAndServe(":8080", nil)
}

//Helper function to import json from file to map
func importJSONDataFromFile(fileName string, result interface{}) (isOK bool) {
	isOK = true
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Print("Error:", err)
		isOK = false
	}
	err = json.Unmarshal(content, result)
	if err != nil {
		isOK = false
		fmt.Print("Error:", err)
	}
	return
}
