# Techniques for enforcing rate limits
### OBJECTIVES

- maximize client satisfaction
- minimize the cost attrition on our end
----
# Rate Limiter Algorithms
- ❌ Token Bucket ```` could lead to a race condition````
- ❌ Leaky Bucket ```` not good for a distributed system````
- ❌ Fixed Window Counter ```` allows more request than necessary````
- ❌ Sliding Logs ```` not great for scalable APIs````
- ✅Sliding Window Counter ````PERFECT for us````
----
> if I understand your business correctly, 
> you serve large cooperation and big-name clients 
> so the only source of error that would trigger the
> need to use rate limits would be a developer/architecture error
----
>we need to make sure that our AWS/GCP bills don't 
> rack up, at the same time we want to provide the best
> service for our clientele

# decision
we will use the ````Sliding Window Counter```` method

![](sliding_window_ctr.png)

Here the window time is broken down into smaller buckets — and the size of each bucket depends on the rate-limit threshold. Each bucket stores the request count corresponding to the bucket range, which constantly keeps moving across time, while smoothing outbursts of traffic.