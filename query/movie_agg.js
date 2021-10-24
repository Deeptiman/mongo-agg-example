db.getCollection("movies").aggregate(

	// Pipeline
	[
		// Stage 1
		{
			$match: {
			    // enter query here
			    "imdb.rating": {"$lt": 7.5}
			}
		},

		// Stage 2
		{
			$group: {
			    _id: "total_vote",
			    vote_count: { 
			        "$sum": "$imdb.votes"
			    }
			}
		},
	],

	// Options
	{

	}

	// Created with Studio 3T, the IDE for MongoDB - https://studio3t.com/

);
