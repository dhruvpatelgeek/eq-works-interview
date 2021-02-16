var GLOBAL_RATE_LIMIT=3000; // maximum number of requests you can handle accorss users
var TIME_LIMIT=10; // how big is the sliding window
var RATE_LIMIT_TIME_WINDOW=100;// how big is the sliding window for th user
var LIMIT=10;

const http = require('http');
const _URL = "http://localhost:5555";
const delay = 1;

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

const endpoints = ['/', '/events/hourly', '/events/daily', '/stats/hourly', '/stats/daily', '/poi'];

async function test_end_point_(test_url){
    var res_ctr=0;
    var rate_limit_ctr=0;
    var error_ctr=0;
    for(let i=0;i<LIMIT;i++)
    {
        await sleep(delay);
        http.get(_URL+test_url, (resp) => {
            let data = '';
            // A chunk of data has been received.
            resp.on('data', (chunk) => {
                data += chunk;
            });
            // The whole response has been received. Print out the result.
            resp.on('end', () => {
                res_ctr++;
               if(data=="RATE LIMITED")
                {
                    rate_limit_ctr++;
                }
            });
        }).on("error", (err) => {
            console.log("Error: " + err.message);
        });
    }

    console.log("[TEST RESULT] "+test_url+ " [response "+res_ctr+"/"+LIMIT+"] [rate limited "+rate_limit_ctr+"/"+res_ctr+"]");
    console.log("Delay->"+delay+" LIMIT "+LIMIT);
}

async function test(){
    for (let i=0;i<endpoints.length;i++){
       await test_end_point_(endpoints[i]);
        await sleep(10000);
    }
    console.log("[TESTS DONE]")
}

test();






