const http = require('http');

// Function to create a server on a specified port
const startServer = (port) => {
    const server = http.createServer((req, res) => {
        console.log(`Received request on port ${port}`);
        res.writeHead(200, { 'Content-Type': 'text/html' });
        res.end(`<html>
                    <head><title>Server ${port}</title></head>
                    <body>
                        <h1>Hello from the web server running on port ${port}</h1>
                    </body>
                 </html>`);
    });

    server.listen(port, () => {
        console.log(`Server is running on http://localhost:${port}`);
    });
};

// Start servers on ports 8080, 8081, and 8082
[8080, 8081, 8082].forEach((port) => startServer(port));
