#!/usr/bin/env nodejs

var fs = require('fs');
var http = require('http');
var https = require('https');
var path = require('path');

function mkdirsSync(name) {
    var parts = name.split('/');
    var current;
    if (name[0] == '/') {
        current = '';
        parts = parts.slice(1);
    } else {
        current = '.';
    }
    parts.forEach((part) => {
        current += '/' + part;
        if (!fs.existsSync(current)) {
            fs.mkdirSync(current);
        }
    });

}

function nop() {
}

var CACHE_BASE_DIRECTORY = path.join(path.dirname(path.dirname(__filename)), 'maps_backend', 'tile_cache');
if (!fs.existsSync(CACHE_BASE_DIRECTORY)) {
    fs.mkdirSync(CACHE_BASE_DIRECTORY);
}

http.createServer(function (req, res) {
    var url = '/' + req.url.split('/').slice(2).join('/');
    var fileName = CACHE_BASE_DIRECTORY + url;
    if (fs.existsSync(fileName)) {
        fs.readFile(fileName, {}, (err, data) => {
            if (err) throw err;
            console.log(url + ': locally cached');
            res.writeHead(200, {'Content-Type': 'image/png'});
            res.write(data);
            res.end();
        });
    } else {
        var parts = url.split('/').slice(1);
        if (parts.length != 4) {
            res.writeHead(404, {'Content-Type': 'text/plain'});
            res.end('Nothing to see here');
            return;
        }

        console.log(url + ': fetch from osm');
        var serverName = parts[0] + ".tile.openstreetmap.org";
        var pathName = parts[1] + "/" + parts[2] + "/" + parts[3];
        const options = {
            hostname: serverName,
            port: 443,
            path: pathName,
            method: 'GET'
        };

        const req = https.request(options, (subres) => {
            console.log('statusCode:', subres.statusCode);
            console.log('headers:', subres.headers);

            if (subres.status != 200) {
                res.writeHead(status, {'Content-Type': 'text/plain'});
                res.end('Status: ' + subres.status);
                return;
            }

            mkdirsSync(path.dirname(fileName));
            // TODO Need to lock the file written to?
            // If w requests want the same file at exactly the same time, we
            // might get problems here..
            var fd = fs.openSync(fileName, 'w');
            res.writeHead(200, {'Content-Type': 'image/png'});

            subres.on('data', (data) => {
                fs.write(fd, data);
                res.write(data);
            });

            subres.on('end', () => {
                fs.close(fd, nop);
                res.end();
            });
        });

        req.on('error', (e) => {
            console.log("Error while fetching " + url + ":");
            console.error(e);
            res.writeHead(500, {'Content-Type': 'text/plain'});
            res.end('');
        });
        req.end();
    }
}).listen(9911, 'localhost');
console.log('Node Tile Server running at http://localhost:9911/');
