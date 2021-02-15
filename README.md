# Table of contents

- [path 1 server-side rate-limiting NODE varient](#path-1-server-side-rate-limiting-node-varient)
  - [Implementation](#implementation)
    - [```` node-cache````](#-node-cache)
  - [Design](#design)
    - [OBJECTIVES](#objectives)
    - [Rate Limiter Algorithms](#rate-limiter-algorithms)
    - [decision](#decision)
- [path 2b back end track](#path-2b-back-end-track)
  - [Requirements](#requirements)
  - [Design](#design)
    - [Communication](#communication)
    - [Architecture](#architecture)

# path 1 server-side rate-limiting NODE varient



-------



## Implementation

We will use

### ```` node-cache````

![Logo](https://raw.githubusercontent.com/node-cache/node-cache/HEAD/logo/logo.png)

[![Node.js CI](https://github.com/node-cache/node-cache/workflows/Node.js%20CI/badge.svg?branch=master)](https://github.com/node-cache/node-cache/actions?query=workflow%3A"Node.js+CI"+branch%3A"master") ![Dependency status](https://img.shields.io/david/node-cache/node-cache) [![NPM package version](https://img.shields.io/npm/v/node-cache?label=npm%20package)](https://www.npmjs.com/package/node-cache) [![NPM monthly downloads](https://img.shields.io/npm/dm/node-cache)](https://www.npmjs.com/package/node-cache) [![GitHub issues](https://img.shields.io/github/issues/node-cache/node-cache)](https://github.com/node-cache/node-cache/issues) [![Coveralls Coverage](https://img.shields.io/coveralls/node-cache/node-cache.svg)](https://coveralls.io/github/node-cache/node-cache)



when a user makes a request we fist see how many requests have been made in the past ````TIME-WINDOW````

if it is below the limit will will serve the user and 



> add the key- value to the table using ````myCache.set( key, val, [ ttl ] )````when its ttl or time to live hits zero thr key will be automatically deleted



Else if it is greater than 



> will will respond with an```` **HTTP** 429 Too Many Requests **response status code** ````indicating the user has sent too many requests in a given amount of time ("**rate limiting**").





## Design



### OBJECTIVES

- maximize client satisfaction
- minimize the cost attrition on our end
----
### Rate Limiter Algorithms
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
----
inspiration from 
![](./ws-product-nodejs/swc.png)
A. H. Fahim Raouf and J. Abouei, "Cache Replacement Scheme Based on Sliding Window and TTL for Video on Demand," Electrical Engineering (ICEE), Iranian Conference on, Mashhad, 2018, pp. 499-504, doi: 10.1109/ICEE.2018.8472723.

### decision
we will use the ````Sliding Window Counter```` method

![](./ws-product-nodejs/sliding_window_ctr.png)

Here the window time is broken down into smaller buckets — and the size of each bucket depends on the rate-limit threshold. Each bucket stores the request count corresponding to the bucket range, which constantly keeps moving across time, while smoothing outbursts of traffic.

When the sum of the counters with timestamps in the past time-slot exceeds the request threshold, User 1 has exceeded the rate limit.

-------



# path 2b back end track

in order to show my distributed system prowess i am goinging desing a fault tolertent  satellite-based  storages in addiiton to what has been asked 

#### 

## Requirements

- fast 
- cheap
- low memory footprint

## Design



### Communication

i will ditch the application layer communication system ````http````

 in favour of transport layer communication   ````UDP````

- faster communication
- 100x lower cost

we will use 

Google's protocol buffers on top of this to make it as light a feather 

![](https://miro.medium.com/max/1400/1*2G7HXILlV5MUIHeNjiYZPA.png)

(https://www.researchgate.net/publication/311461272_Performance_evaluation_of_using_Protocol_Buffers_in_the_Internet_of_Things_communication)



### Architecture

![](/Users/dhruvpatel/Desktop/eq_careers/eq-works-interview/ws-product-golang/architecture.png)





The client will encode the message in GCP and send a udp bye stream to the server. 

we will used IEEE checksum for maintinaing data intergrity 

 