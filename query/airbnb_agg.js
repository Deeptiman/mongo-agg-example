db.getCollection("listingsAndReviews").aggregate(

	// Pipeline
	[
		// Stage 1
		{
			$match: {
			    // enter query here
			    "beds": {"$gte": 0}
			}
		},
		
		// Stage 2
		{
			$group: {
			    _id: "total_bed",
			    vote_count: { 
			        "$sum": "$bedrooms"
			    }
			}
		},
	],
				
	// Options
	{

	}

	// Created with Studio 3T, the IDE for MongoDB - https://studio3t.com/

);
