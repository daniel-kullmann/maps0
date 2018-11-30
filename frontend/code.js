var map = null;
var gpx_files = {};

function show_gpx_file(show, file) {
    if (show) {
        $.ajax({
            type: 'GET',
            url: 'http://localhost:8000/api/gpx/get/' + file,
            //data: {...},
            success: function(data, textStatus, request) {
                var points = data.getElementsByTagName('trkpt');
                var markers = [];
                for (i=0; i<points.length; i++) {
                    var attr = points.item(i).attributes;
                    var lat = attr.getNamedItem("lat").value;
                    var lon = attr.getNamedItem("lon").value;
                    var r = L.marker([lat, lon]).addTo(map);
                    markers.push(r);
                    gpx_files[file] = {
                        markers: markers
                    };
                }
            },
            error: function(request, textStatus, error) {
                $("#error-message").text(error);
                //showError('Oops, there was a problem retrieving the comments.');
            },
            dataType: 'xml'
        });
    } else {
        if (gpx_files[file]) {
            if (gpx_files[file].markers) {
                gpx_files[file].markers.forEach(function(e) {
                    e.remove();
                });
            }
        }
    }
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

    $.ajax({
        type: 'GET',
        url: 'http://localhost:8000/api/gpx/',
        //data: {...},
        success: function(data, textStatus, request) {
            var text = "<ul>";
            data.forEach(function(e) {
                text += "<li>";
                text += "<input type=\"checkbox\" ";
                text += "onClick=\"show_gpx_file(this.checked,'" + e + "');\">";
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

});
