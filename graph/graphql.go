package graph

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"net/http"
	"testApplication/interfaces"
	"testApplication/models"
)

type graph struct {
	repo   interfaces.ClientRepo
	schema graphql.Schema
}

func NewGraph(repo interfaces.ClientRepo) (*graph, error) {
	graph := graph{
		repo: repo,
	}
	var clientType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Client",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.ID,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	})

	var queryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"client": &graphql.Field{
				Type:        clientType,
				Description: "Get client by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {

					if p.Args["id"] == nil {
						return nil, errors.New("id is empty")
					}
					id := p.Args["id"].(int)

					return graph.repo.GetClientById(context.TODO(), id)
				},
			},
			"clients": &graphql.Field{
				Type:        graphql.NewList(clientType),
				Description: "Get clients",
				Args: graphql.FieldConfigArgument{
					"limit": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"offset": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {

					offset := 0
					if p.Args["offset"] != nil {
						offset, _ = p.Args["offset"].(int)
					}
					limit := 0
					if p.Args["limit"] != nil {
						limit, _ = p.Args["limit"].(int)
					}

					return graph.repo.GetClients(context.TODO(), offset, limit)
				},
			},
		}})

	var mutationType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"create": &graphql.Field{
				Type: clientType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Description: "Add client",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					client := models.Client{
						Name: p.Args["name"].(string),
					}
					return graph.repo.CreateClient(context.TODO(), client)
				},
			},
			"update": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"name": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Description: "Update client by id",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					client := models.Client{
						Id:   p.Args["id"].(int),
						Name: p.Args["name"].(string),
					}
					err := graph.repo.UpdateClient(context.TODO(), client)
					if err != nil {
						return nil, err
					}
					return "Success", nil
				},
			},
			"delete": &graphql.Field{
				Type: graphql.String,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Description: "Delete client by id",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {

					if p.Args["id"] == nil {
						return nil, errors.New("id is empty")
					}
					id := p.Args["id"].(int)

					err := graph.repo.DeleteClient(context.TODO(), id)
					if err != nil {
						return nil, err
					}
					return "Success", nil
				},
			},
		},
	})

	var schema, _ = graphql.NewSchema(graphql.SchemaConfig{Query: queryType, Mutation: mutationType})

	graph.schema = schema

	return &graph, nil
}

type RequestParams struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
}

func (graph graph) GraphqlHandler(c *gin.Context) {
	var reqObj RequestParams
	if err := c.ShouldBindJSON(&reqObj); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := graphql.Do(graphql.Params{
		Context:        c,
		Schema:         graph.schema,
		RequestString:  reqObj.Query,
		VariableValues: reqObj.Variables,
		OperationName:  reqObj.Operation,
	})

	if len(result.Errors) > 0 {
		c.IndentedJSON(http.StatusBadRequest, result.Errors)
		return
	}
	c.IndentedJSON(http.StatusOK, result)
}
