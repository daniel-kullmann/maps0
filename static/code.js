$(document).ready(function() {
    var map = L.map('map').setView([37.31915, -8.8033], 13);

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        useCache: true,
        cacheMaxAge: 30 * 24 * 3600 * 1000, // 30 days as milliseconds
	      crossOrigin: true
    }).addTo(map);

    var popup = L.popup();

    function onMapClick(e) {
        popup
            .setLatLng(e.latlng)
            .setContent("You clicked the map at " + e.latlng.toString())
            .openOn(map);
    }

    map.on('click', onMapClick);


    $.ajax({
        type: 'GET',
        url: '1.gpx',
        //data: {...},
        success: function(data, textStatus, request) {
            var points = data.getElementsByTagName('trkpt');
            for (i=0; i<points.length; i++) {
                var attr = points.item(i).attributes;
                var lat = attr.getNamedItem("lat").value;
                var lon = attr.getNamedItem("lon").value;
                var r = L.marker([lat, lon]).addTo(map);
                console.log(r);
            }
        },
        error: function(request, textStatus, error) {
            console.log(error);
            //showError('Oops, there was a problem retrieving the comments.');
        },
        dataType: 'xml'
    });
});
