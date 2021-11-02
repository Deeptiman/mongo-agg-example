package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// #1
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

// #2
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

// #3
func geoNearAggQuery(ctx context.Context, client *mongo.Client, lat, long string) interface{} {

	database := client.Database("movie_details")
	collection := database.Collection("theater_list")

	var latC, longC float64
	if l, err := strconv.ParseFloat(lat, 32); err == nil {
		latC = l
	}

	if l, err := strconv.ParseFloat(long, 32); err == nil {
		longC = l
	}

	query := []bson.M{
		bson.M{
			"$geoNear": bson.M{
				"near": bson.M{
					"type": "Point",
					"coordinates": []interface{}{latC, longC},
				},
				"minDistance": 2,
				"maxDistance": 2000000,
				"distanceField": "distance",
				"includeLocs": "geo",
				"spherical": true,
			},
		},
		bson.M{
			"$project": bson.M{
				"location.geo": 0,
			},
		},
	}

	return ExecuteQuery(ctx, query, collection)
}

// #4
func graphLookup(ctx context.Context, client *mongo.Client, email string) interface{} {

	database := client.Database("movie_details")
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
				"from": "comment_list",
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
			"$lookup": bson.M{
				"from": "movie_list",
				"let": bson.M{
					"movieId": bson.M{
						"$toObjectId": "$comments.movie_id",
					},
				},
				"pipeline": []bson.M{
					bson.M{
						"$match": bson.M{
							"$expr": bson.M{
								"$eq": []interface{}{"$_id", "$$movieId"},
							},
						},
					},
					bson.M{
						"$project": bson.M{
							"title": "$title",
							"poster": "$poster",
						},
					},
				},
				"as": "movie_details",
			},
		},
		bson.M{
			"$unwind": "$movie_details",
		},
		bson.M{
			"$project": bson.M{
				"name": 0,
				"email": 0,
			},
		},
	}

	return ExecuteQuery(ctx, query, collection)
}

// #5
func commentAggQuery(ctx context.Context, client *mongo.Client, genre string) interface{}{

	database := client.Database("movie_details")
	collection := database.Collection("movie_list")

	fmt.Println("Genre -- ", genre)

	query := []bson.M{
		 // Stage - 1
		bson.M{
			"$unwind": "$genres",
		},
		// Stage - 2
		bson.M{
			"$match": bson.M{ 
				"genres": genre,
				"imdb.rating" : bson.M{
					"$gt": 5,
				},           
			},
		},
		// Stage - 3
		bson.M{
			"$lookup": bson.M{
				"from": "comment_list",
				"localField": "_id",
                "foreignField": "movie_id",
				// "let": bson.M{
				// 	"movieId": bson.M{
				// 		"$toObjectId": "$_id",
				// 	},
				// },
				// "pipeline": []bson.M{
				// 	bson.M{
				// 		"$match": bson.M{
				// 			"$expr": bson.M{
				// 				"$eq": []interface{}{"$movie_id", "$$movieId"},
				// 			},
				// 		},
				// 	},
				// 	bson.M{
				// 		"$sort" : bson.M{
				// 			"date": -1,
				// 		},					
				// 	},				
				// },
				"as": "comments",
			},
		},
		// Stage - 4
		bson.M{
			"$project": bson.M{
				"_id": -1,
				"title": -1,
				"comments": 1,
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