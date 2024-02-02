
from http.server import HTTPServer, SimpleHTTPRequestHandler
import ssl,os

keypath = os.path.abspath("server.key")
crtpath = os.path.abspath("server.crt")

os.chdir('./Meecha_web')

port = 11333
httpd = HTTPServer(('0.0.0.0', port), SimpleHTTPRequestHandler)
httpd.socket = ssl.wrap_socket(httpd.socket, keyfile=keypath, certfile=crtpath, server_side=True)

print("Server running on https://0.0.0.0:" + str(port))

httpd.serve_forever()