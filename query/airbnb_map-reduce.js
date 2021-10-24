db.listingsAndReviews.mapReduce(
        function(){
            emit("total",this.bedrooms);
        }, 
        function(key,vals){
            var count = 0;
            vals.forEach(function(v) {
                if (v !== undefined){
                  count +=v;
                }
            });
            
            return count;
        }, 
        {
            query: {"beds": {"$gte": 0}},
            out: "bed_counts"
        }
).find()