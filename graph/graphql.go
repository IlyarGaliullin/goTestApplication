package graph

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"net/http"
	"strings"
	"testApplication/interfaces"
	"testApplication/models"
	"testApplication/redis"
)

type graph struct {
	clientRepo interfaces.ClientRepo
	userRepo   interfaces.UserRepo
	schema     graphql.Schema
}

func getToken(c *gin.Context) (token string) {
	bearerToken := c.GetHeader("Authorization")

	if len(strings.Split(bearerToken, " ")) == 2 {
		token = strings.Split(bearerToken, " ")[1]
		return
	}
	return
}

func checkUserGrant(c *gin.Context, redisConn *redis.Connection, userRepo interfaces.UserRepo, table string, operation string) (int, bool, error) {
	token := getToken(c)

	userId, err := redisConn.CheckToken(c, token)

	if err != nil {
		return userId, false, err
	}

	grant, err := userRepo.CheckUserGrant(
		c,
		userId,
		table,
		operation,
	)

	if err != nil {
		return userId, false, err
	}

	return userId, grant, nil
}

func checkPermissions(requestFields []ast.Selection, allowedFields map[string]bool) []ast.Selection {
	authorizedFields := make([]ast.Selection, 0)

	for _, field := range requestFields {
		fieldName := field.(*ast.Field).Name.Value
		if allowed, ok := allowedFields[fieldName]; ok && allowed {
			authorizedFields = append(authorizedFields, field)

		}
	}

	return authorizedFields
}

func fillFieldsRecursive(parentField string, field *ast.Field, userId int, graph *graph) *ast.Field {

	newAST := field

	var newSelections []ast.Selection
	for _, v1 := range field.SelectionSet.Selections {

		if v1.(*ast.Field).SelectionSet != nil {
			newSelections = append(newSelections, fillFieldsRecursive(v1.(*ast.Field).Name.Value, v1.(*ast.Field), userId, graph))
			continue
		}
		newSelections = append(newSelections, v1)
	}

	permissions, err := graph.userRepo.GetFieldsPermissions(context.TODO(), userId, parentField)
	if err != nil {
		return nil
	}

	authorizedSelection := checkPermissions(newSelections, permissions)
	newAST.SelectionSet.Selections = authorizedSelection

	return newAST
}

func NewGraph(redisConn *redis.Connection, clientRepo interfaces.ClientRepo, userRepo interfaces.UserRepo) (*graph, error) {
	graph := graph{
		clientRepo: clientRepo,
		userRepo:   userRepo,
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

					userId, grant, err := checkUserGrant(p.Context.(*gin.Context), redisConn, userRepo, "clients", "read")

					if err != nil {
						return nil, err
					}
					if !grant {
						return nil, redis.ErrUnauthorized
					}

					var newASTs []*ast.Field

					newASTs = append(newASTs, fillFieldsRecursive(p.Info.FieldASTs[0].Name.Value, p.Info.FieldASTs[0], userId, &graph))

					p.Info.FieldASTs = newASTs

					if p.Args["id"] == nil {
						return nil, errors.New("id is empty")
					}
					id := p.Args["id"].(int)

					return graph.clientRepo.GetClientById(context.TODO(), id)
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

					userId, grant, err := checkUserGrant(p.Context.(*gin.Context), redisConn, userRepo, "clients", "read")

					if err != nil {
						return nil, err
					}
					if !grant {
						return nil, redis.ErrUnauthorized
					}

					var newASTs []*ast.Field

					newASTs = append(newASTs, fillFieldsRecursive(p.Info.FieldASTs[0].Name.Value, p.Info.FieldASTs[0], userId, &graph))

					p.Info.FieldASTs = newASTs
					offset := 0
					if p.Args["offset"] != nil {
						offset, _ = p.Args["offset"].(int)
					}
					limit := 0
					if p.Args["limit"] != nil {
						limit, _ = p.Args["limit"].(int)
					}

					return graph.clientRepo.GetClients(context.TODO(), offset, limit)
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

					userId, grant, err := checkUserGrant(p.Context.(*gin.Context), redisConn, userRepo, "clients", "create")

					if err != nil {
						return nil, err
					}
					if !grant {
						return nil, redis.ErrUnauthorized
					}

					var newASTs []*ast.Field

					newASTs = append(newASTs, fillFieldsRecursive(p.Info.FieldASTs[0].Name.Value, p.Info.FieldASTs[0], userId, &graph))

					p.Info.FieldASTs = newASTs

					client := models.Client{
						Name: p.Args["name"].(string),
					}
					return graph.clientRepo.CreateClient(context.TODO(), client)
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

					userId, grant, err := checkUserGrant(p.Context.(*gin.Context), redisConn, userRepo, "clients", "update")

					if err != nil {
						return nil, err
					}
					if !grant {
						return nil, redis.ErrUnauthorized
					}
					var newASTs []*ast.Field

					newASTs = append(newASTs, fillFieldsRecursive(p.Info.FieldASTs[0].Name.Value, p.Info.FieldASTs[0], userId, &graph))

					p.Info.FieldASTs = newASTs

					client := models.Client{
						Id:   p.Args["id"].(int),
						Name: p.Args["name"].(string),
					}
					err = graph.clientRepo.UpdateClient(context.TODO(), client)
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

					userId, grant, err := checkUserGrant(p.Context.(*gin.Context), redisConn, userRepo, "clients", "delete")

					if err != nil {
						return nil, err
					}
					if !grant {
						return nil, redis.ErrUnauthorized
					}
					var newASTs []*ast.Field

					newASTs = append(newASTs, fillFieldsRecursive(p.Info.FieldASTs[0].Name.Value, p.Info.FieldASTs[0], userId, &graph))

					p.Info.FieldASTs = newASTs

					if p.Args["id"] == nil {
						return nil, errors.New("id is empty")
					}
					id := p.Args["id"].(int)

					err = graph.clientRepo.DeleteClient(context.TODO(), id)
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

	if result.HasErrors() {
		if result.Errors[0].Error() == redis.ErrUnauthorized.Error() {
			c.IndentedJSON(http.StatusForbidden, result.Errors)
			return
		}
		c.IndentedJSON(http.StatusBadRequest, result.Errors)
		return
	}
	c.IndentedJSON(http.StatusOK, result)
}
