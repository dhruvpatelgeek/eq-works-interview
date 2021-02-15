const http = require('http');
var _URL="http://localhost:5555/"

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function test_end_point(){
    for(let i=0;i<200;i++)
    {
        await sleep(10);
        http.get(_URL, (resp) => {
            let data = '';
            // A chunk of data has been received.
            resp.on('data', (chunk) => {
                data += chunk;
            });
            // The whole response has been received. Print out the result.
            resp.on('end', () => {
                console.log(data);
            });
        }).on("error", (err) => {
            console.log("Error: " + err.message);
        });
    }
}

test();

