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
		bson.M{
			"$lookup": bson.M{
				"from": "comments",
				"let": bson.M{
					"movieId": bson.M{
						"$toObjectId": "$movie_details._id",
					},
				},
				"pipeline": []bson.M{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{
								"$eq": []interface{}{"$movie_id", "$$movieId"},
							},
						},
					},			
				},
				"as": "comments",
			},
		},
	}

	return ExecuteQuery(ctx, query, collection)
}

func bucketAggQuery(ctx context.Context, client *mongo.Client) interface{} {

	database := client.Database("sample_mflix")
	collection := database.Collection("movies")

	query := []bson.M{
		bson.M{
			"$bucket": bson.M{
				"groupBy": "$tomatoes.viewer.numReviews",
				"boundaries": []interface{}{500, 1200, 1500, 2000, 3500},
				"default": "Other",
				"output": bson.M{
					"count": bson.M{
						"$sum": 1,
					},
					// "movies": bson.M{
					// 	"$push": bson.M{
					// 		"title": "$title",
					// 		"year": "$year",
					// 		"poster": "$poster",
					// 	},
					// },
				},
			},
		},
	}

	return ExecuteQuery(ctx, query, collection)
}

func graphLookup(ctx context.Context, client *mongo.Client, email string) interface{} {

	database := client.Database("sample_mflix")
	collection := database.Collection("users")

	query := []bson.M{
		bson.M{
			"$match": bson.M{
				"email": bson.M{
					"$eq": email,
				},
			},
		},
		bson.M{
			"$graphLookup": bson.M{
				"from": "comments",
				"startWith": "$email",
				"connectFromField": "email",
				"connectToField": "email",
				"maxDepth": 3,
				"as": "comments",
			},
		},
		bson.M{
			"$unwind": "$comments",
		},		
		bson.M{
			"$group": bson.M{
				"_id": "$email",
				"comments": bson.M{
					"$push": bson.M{
						"movie_id": "$comments.movie_id",
						"comment": "$comments.text",
						"date": "$comments.date",
					},
				},
			},
		},
	}

	return ExecuteQuery(ctx, query, collection)
}

func ExecuteQuery(ctx context.Context, query []bson.M, collection *mongo.Collection) interface{} {

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