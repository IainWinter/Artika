const express = require('express');
const helmet = require('helmet');
const path = require('path');

const app = express();

app.use(helmet({
    crossOriginOpenerPolicy: 'same-origin; include:same-origin-allow-popups',
    contentSecurityPolicy: {
        directives: {
            defaultSrc: ["'self'"],
            scriptSrc: ["'self'", "https://apis.google.com", "https://accounts.google.com/gsi/client", "'unsafe-inline'"],
            frameSrc: ["'self'", "https://accounts.google.com"],
        },
    },
    ReferrerPolicy: { policy: 'no-referrer-when-downgrade' },
}));

// Your other routes and middleware...

const port = 3000;

var options = {
    index: "google_test.html"
};

app.get('/', (req, res) => {
    res.sendFile(path.join(__dirname, 'public', 'google_test.html'));
});

app.set('hostname', 'localhost');

app.listen(port, () => {
    console.log(`Server is listening on port ${port}`);
});