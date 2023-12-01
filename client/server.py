from http.server import SimpleHTTPRequestHandler
import socketserver

# This server is just to test a localhost origin for google sign-in
# The backend folder has the real server

files = {
    "/": "index.html",
    "/login.js": "login.js",
    "/config.js": "config.js",
}

class CustomHandler(SimpleHTTPRequestHandler):
    def end_headers(self):
        self.send_header('Referrer-Policy', 'no-referrer-when-downgrade')
        self.send_header('Access-Control-Allow-Origin', 'http://localhost:3000')
        super().end_headers()
    
    def do_GET(self):
        if self.path in files:
            self.path = '/public/' + files[self.path]

        return super().do_GET()

port = 3000
httpd = socketserver.TCPServer(('localhost', port), CustomHandler)

print(f'Server listening on http://localhost:{port}')

try:
    httpd.serve_forever()
except KeyboardInterrupt:
    pass