var map = null;
var gpx_files = [];

var colors = [
    ['red', 0],
    ['blue', 0],
    ['green', 0]
];

function bound_map_by_gpx_tracks() {
    if (gpx_files.length > 0) {
        var all_bounds = gpx_files.map(function(e) { return e.bounds; });
        var super_bounds = all_bounds.reduce(function(a,b) {return a.extend(b);});
        map.fitBounds(super_bounds);
    }
}

function drop_color(color) {
    for (var i=0; i<colors.length; i++) {
        if (colors[i][0] == color) {
            colors[i][1]--;
        }
    }
}

function pick_color() {
    var lowest = 0;
    for (var i=1; i<colors.length; i++) {
        if (colors[i][1] < colors[lowest][1]) {
            lowest = i;
        }
    }
    colors[lowest][1]++;
    return colors[lowest][0];
}

function show_gpx_file(show, file) {
    if (show) {
        $.ajax({
            type: 'GET',
            url: 'http://localhost:8000/api/gpx/get/' + file,
            //data: {...},
            success: function(data, textStatus, request) {
                var color = pick_color();
                var points = data.getElementsByTagName('trkpt');
                var markers = [];
                var line_coordinates = [];
                for (i=0; i<points.length; i++) {
                    var attr = points.item(i).attributes;
                    var lat = attr.getNamedItem("lat").value;
                    var lon = attr.getNamedItem("lon").value;
                    line_coordinates.push([lat,lon]);
                    if (i==0 || i == points.length-1) {
                        var r = L.marker([lat, lon]).addTo(map);
                        markers.push(r);
                    }
                }
                var poly_line = L.polyline(line_coordinates, {color: color});
                poly_line.addTo(map);
                gpx_files.push({
                    file: file,
                    markers: markers,
                    poly_line: poly_line,
                    bounds: poly_line.getBounds()
                });
            },
            error: function(request, textStatus, error) {
                $("#error-message").text(error);
                //showError('Oops, there was a problem retrieving the comments.');
            },
            dataType: 'xml'
        });
    } else {
        var gpx_file = gpx_files.find(function (e) { return e.file == file; });
        if (gpx_file) {
            if (gpx_file.markers) {
                gpx_file.markers.forEach(function(e) {
                    e.remove();
                });
            }
            if (gpx_file.poly_line) {
                gpx_file.poly_line.remove();
            }
            gpx_files = gpx_files.filter(function (e) { return e.file != file; });
        }
    }
}

function load_gpx_track_list() {
    $.ajax({
        type: 'GET',
        url: 'http://localhost:8000/api/gpx/',
        //data: {...},
        success: function(data, textStatus, request) {
            var text = "<ul>";
            data.forEach(function(e) {
                text += "<li>";
                text += "<input type=\"checkbox\" ";
                text += "onClick=\"show_gpx_file(this.checked,'" + e + "');\"";
                if (gpx_files.find(function (f) { return f.file == e })) {
                    text += " checked";
                }
                text += ">";
                text += e;
                text += "</a>";
                text += "</li>";
            });
            text += "</ul>";
            $("#gpx-list").html(text);
        },
        error: function(request, textStatus, error) {
            $("#error-message").text(error);
            //showError('Oops, there was a problem retrieving the comments.');
        },
        dataType: 'json'
    });
}


$(document).ready(function() {

    $("#map").css({
        width: $(window).width()-200,
        height: $(window).height()
    });

    map = L.map('map').setView([37.31915, -8.8033], 13);

    L.tileLayer('http://localhost:8000/api/tiles/{s}/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
	      crossOrigin: true
    }).addTo(map);

    load_gpx_track_list();

});
