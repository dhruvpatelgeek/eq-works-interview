# path 2b back end track

in order to show my distributed system prowess i am goinging desing a fault tolertent  satellite-based  storages in addiiton to what has been asked 

#### 

## Requirements

- fast 
- cheap
- low memory footprint

## design

### communication

i will ditch the application layer communication system ````UDP```` instead of  ````http```` in favour of transport layer communication for 

- faster communication
- 100x lower cost

we will use 

Google's protocol buffers on top of this to make it as light a feather 

![](https://miro.medium.com/max/1400/1*2G7HXILlV5MUIHeNjiYZPA.png)

(https://www.researchgate.net/publication/311461272_Performance_evaluation_of_using_Protocol_Buffers_in_the_Internet_of_Things_communication)



### Architecture

![](https://github.com/dhruvpatelgeek/eq-works-interview/blob/master/ws-product-golang/architecture.png)





The client will encode the message in GCP and send a udp bye stream to the server. 

we will used IEEE checksum for maintinaing data intergrity 

 
