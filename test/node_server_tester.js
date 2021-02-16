var GLOBAL_RATE_LIMIT=3000; // maximum number of requests you can handle accorss users
var TIME_LIMIT=10; // how big is the sliding window
var RATE_LIMIT_TIME_WINDOW=100;// how big is the sliding window for th user
var LIMIT=10;
const delay = 100;

const fetch = require('node-fetch');
const _URL = "http://bribchat.com:5555";

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

const endpoints = ['/'];

async function test_end_point_(test_url){
    var res_ctr=0;
    var rate_limit_ctr=0;
    var error_ctr=0;
    try {
      const res = await fetch(_URL+test_url);
      console.log('Status Code:', res.status);
      res_ctr++;
      if (res.status==429)
      {
        rate_limit_ctr++;
      }
      
    } catch (err) {
      console.log(err.message); //can be console.error
      error_ctr++;
    }
    console.log("[TEST RESULT] "+test_url+ " [response "+res_ctr+"/"+LIMIT+"] [rate limited "+rate_limit_ctr+"/"+res_ctr+"]");
    console.log("Delay->"+delay+" LIMIT "+LIMIT);
  };


async function test(){
    for (let i=0;i<endpoints.length;i++){
       await test_end_point_(endpoints[i]);
       await sleep(delay)
    }
    console.log("[TESTS DONE]")
}

test();






