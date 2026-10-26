package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/graphql-go/graphql"
	"order-management/service"
)

func RegisterGraphQL(r *gin.Engine, svc *service.OrderService) {
	orderType := graphql.NewObject(graphql.ObjectConfig{Name: "Order", Fields: graphql.Fields{
		"id": &graphql.Field{Type: graphql.NewNonNull(graphql.ID)}, "customerId": &graphql.Field{Type: graphql.Int}, "productId": &graphql.Field{Type: graphql.Int},
		"quantity": &graphql.Field{Type: graphql.Int}, "unitPrice": &graphql.Field{Type: graphql.Float}, "totalAmount": &graphql.Field{Type: graphql.Float},
		"status": &graphql.Field{Type: graphql.String}, "createdAt": &graphql.Field{Type: graphql.String}, "updatedAt": &graphql.Field{Type: graphql.String},
	}})
	query := graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: graphql.Fields{
		"order": &graphql.Field{Type: orderType, Args: graphql.FieldConfigArgument{"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.ID)}}, Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			id, e := strconv.ParseInt(p.Args["id"].(string), 10, 64)
			if e != nil {
				return nil, e
			}
			o, e := svc.Get(p.Context, id)
			if e != nil {
				return nil, e
			}
			return o, nil
		}},
		"orders": &graphql.Field{Type: graphql.NewList(orderType), Resolve: func(p graphql.ResolveParams) (interface{}, error) { return svc.List(p.Context, 100, 0) }},
	}})
	schema, err := graphql.NewSchema(graphql.SchemaConfig{Query: query})
	if err != nil {
		panic(err)
	}
	r.POST("/graphql", func(c *gin.Context) {
		var req struct {
			Query string `json:"query"`
		}
		if c.ShouldBindJSON(&req) != nil || req.Query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query is required"})
			return
		}
		result := graphql.Do(graphql.Params{Schema: schema, RequestString: req.Query, Context: c.Request.Context()})
		c.JSON(http.StatusOK, result)
	})
}
