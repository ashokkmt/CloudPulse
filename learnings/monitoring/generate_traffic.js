const http = require('http');

// The endpoint to hit. 
// You can change this to /api/broken for Experiment 2!
const URL = 'http://localhost:8000/api/broken';

function makeRequest() {
    return new Promise((resolve) => {
        const req = http.get(URL, (res) => {
            // Consume data to free up memory
            res.on('data', () => {});
            res.on('end', resolve);
        });

        req.on('error', () => {
            // Ignore errors for the flood script
            resolve();
        });
        
        req.setTimeout(2000, () => {
            req.destroy();
            resolve();
        });
    });
}

async function floodTraffic() {
    console.log(`Starting traffic flood to ${URL}...`);
    console.log('Press Ctrl+C to stop.');
    
    try {
        while (true) {
            const requests = [];
            // Spawn multiple concurrent requests
            for (let i = 0; i < 20; i++) {
                requests.push(makeRequest());
            }
            
            await Promise.all(requests);
            
            // Sleep briefly to avoid completely freezing your computer!
            await new Promise(r => setTimeout(r, 100));
        }
    } catch (e) {
        console.error(e);
    }
}

floodTraffic();
