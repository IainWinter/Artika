from http.server import SimpleHTTPRequestHandler
import socketserver

class CustomHandler(SimpleHTTPRequestHandler):
    def end_headers(self):
        self.send_header('Referrer-Policy', 'no-referrer-when-downgrade')
        self.send_header('Access-Control-Allow-Origin', 'http://localhost:3000')
        super().end_headers()
    
    def do_GET(self):
        if self.path == '/':
            self.path = '/public/google_test.html'
        return super().do_GET()

port = 3000
httpd = socketserver.TCPServer(('localhost', port), CustomHandler)

print(f'Server listening on http://localhost:{port}')

httpd.serve_forever()