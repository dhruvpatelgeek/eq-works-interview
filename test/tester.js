const https = require('https');
var _URL="http://localhost:5555/"
https.get(_URL, (resp) => {
    let data = '';
    // A chunk of data has been received.
    resp.on('data', (chunk) => {
        data += chunk;
    });
    // The whole response has been received. Print out the result.
    resp.on('end', () => {
        console.log(JSON.parse(data).explanation);
    });
}).on("error", (err) => {
    console.log("Error: " + err.message);
});