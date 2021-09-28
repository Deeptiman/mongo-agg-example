package main

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func filterGenreAgg(ctx context.Context, client *mongo.Client, input_genre string) interface{} {

	database := client.Database("sample_mflix")
	collection := database.Collection("movies")

	fmt.Println("Genre ---> ", input_genre)

	query := []bson.M{
		bson.M{
			"$project": bson.M{
				"movie_details": "$$ROOT",
				"filter_genre": bson.M{
					"$filter": bson.M{
						"input": "$genres",
						"as":    "gnr",
						"cond": bson.M{
							"$eq": []interface{}{"$$gnr", input_genre},
						},
					},
				},
			},
		},
		bson.M{
			"$unwind": "$filter_genre",
		},
		bson.M{
			"$match": bson.M{
				"movie_details.imdb.rating": bson.M{
					"$lt": 5,
				},
			},
		},
	}

	qu, _ := json.MarshalIndent(query, ", ", " ")

	fmt.Println("Queries -->> ", string(qu))

	cursor, err := collection.Aggregate(ctx, query)
	if err != nil {
		fmt.Println("Mongo Error - ", err.Error())
		return nil
	}

	var info []bson.M
	if err = cursor.All(ctx, &info); err != nil {
		fmt.Println("Cursor Error - ", err.Error())
		return nil
	}

	fmt.Println("Total Docs -- ", len(info))

	return info

}