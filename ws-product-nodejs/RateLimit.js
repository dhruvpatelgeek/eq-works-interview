//variables------------------------------------
var GLOBAL_RATE_LIMIT=3000; // maximum number of requests you can handle accorss users
var TIME_LIMIT=10; // how big is the sliding window
var RATE_LIMIT_TIME_WINDOW=5;// how big is the sliding window for th user

/*
WE WILL LIMIT A USER TO 5 (RATE_LIMIT_TIME_WINDOW) REQUESTS 
OVER 10(TIME_LIMIT) SECONDS OF TIME

WE WILL LIMIT ALL THE USERS TO 3000(GLOBAL_RATE_LIMIT) REQUESTS
OVER 10(TIME_LIMIT) SECONDS OF TIME
*/


// the smaller the more perfomance intensive 
// the larger the poor exprience for to client
// i would leave this up to my manager to decide 
var PURGE_PERIOD=1;
var debug =true;
//---------------------------------------------

//DEPENDENCIES---------------------------------
const NodeCache = require( "node-cache" );
const myCache = new NodeCache( { stdTTL: TIME_LIMIT, checkperiod: PURGE_PERIOD,maxKeys:GLOBAL_RATE_LIMIT} );
//---------------------------------------------

//ERROR FUNCTION-------------------------------
class rateErr extends Error {};

//---------------------------------------------
//HELPER FUNCTION------------------------------
/**
 *
 * @param key
 * @returns the the array of requests that the user has made
 */
function getCallStack(key){
    var callStack=[];
    let ctr=0;
    let flag=true;
    while(flag){
        let saltedKey=key+ctr.toString();
        let value=myCache.get(saltedKey);
        if ( value == undefined ){
            flag=false;
        } else {
            callStack.push(value);
            ctr++;
        }
    }
    if (debug)
        console.log(callStack.length.toString());
    return callStack;
}

/**
 *
 * @param key
 * @returns {boolean} if true server the user else rate limit the user
 */
function putIntoStack(key){
    var callStack=getCallStack(key);
    if(callStack.length>RATE_LIMIT_TIME_WINDOW-1)
    {
        if(debug)
            console.log("TOO MANY REQUESTS FROM",key);
        return false
    }
    try {
        let saltedKey=key+callStack.length;
        var success=myCache.set(saltedKey,{date:Date.now()})
        return success;
    } catch (error) {
        console.error("[CRITICAL] GLOBAL RATE LIMIT HIT ERRCODE->"+error.errorcode);
    }

    return false;
}
//---------------------------------------------

//MAIN FUNCTION---------------------------------
/**
 *  middleware function for rate limiting
 * @constructor
 */
function RateLimit (){

    this.middlewareRateLimiter=function (req, res, next) {
       
        if(typeof(req)==typeof(undefined))
		{
            if (debug){
                console.log("request not defined")
            }
            var errObj=new rateErr;
			next(errObj);
		}
        else{
            var result=putIntoStack(req.socket.remoteAddress);
            // use a client UID here to prevent DDOS attacks
            // for not i will simply use the ip address here for sake of simplicity

            console.log(req.socket.remoteAddress);
           if(result){
                if (debug){
                    console.log("added");
                }
                res.status(201);
                next();
           } else {
                if (debug){
                    console.log("not added");
                    res.status(429);
		            res.send("RATE LIMITED");
                }
           }
        }
    }
};
//---------------------------------------------

RateLimit.Error = rateErr;
module.exports = RateLimit;