db.movies.mapReduce(
        function(){
            emit(1,this.imdb.votes);
        }, 
        function(key,vals){
            var count = 0;
            vals.forEach(function(v) {
                count +=v;
            });
            return count;
        }, 
        {
            query: {"imdb.rating": {"$lt": 7.5}},
            out: "vote_counts"
        }
).find()