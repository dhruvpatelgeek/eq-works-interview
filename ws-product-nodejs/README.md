# Rate limiter

### configerable variables
in ````ws-product-nodejs/RateLimit.js````

there are three variables to control the rate limit

````GLOBAL_RATE_LIMIT````

>maximum number of requests you can handle across users

````TIME_LIMIT```` 
>how long do should the request window be

````RATE_LIMIT_TIME_WINDOW````
>how big is the sliding window for the user

### Example
say you have a GCP instance that can 
handle no more than
````5```` total requests 
in ````10 secs````

so set

````TIME_LIMIT=10````

````RATE_LIMIT_TIME_WINDOW=5````


### future features
> individual client rate limiting

I have configured my assignment in such a way that 
if you want to limit rates on a per-client basis you can do so 
by simply changing the ````RATE_LIMIT_TIME_WINDOW```` for that client 
it is fairly trivial to implement here

