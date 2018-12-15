from django.http import HttpResponse

import json
import os
import os.path

DIR='../gpx_store_files'

def list_gpx(request):
    files = os.listdir(DIR)
    response = HttpResponse(json.dumps(files), content_type='application/json')
    response['Access-Control-Allow-Origin'] = '*'
    return response


def get_gpx(request, name):
    filename = os.path.join(DIR, name)
    if os.path.exists(filename):
        with open(filename, 'r') as fh:
            response = HttpResponse(fh, content_type='text/xml')
            response['Access-Control-Allow-Origin'] = '*'
            return response

def save_gpx(request):
    raw_data = json.loads(request.body)
    allowed_keys = ['name', 'date', 'description', 'track_points']
    data = {key: raw_data[key] for key in raw_data if key in allowed_keys}
    file_name = data['date'] + '-' + data['name']
    content = create_pgx_file(data['name'], data['date'], data['description'], data['track_points'])
    with open(os.path.join(DIR,file_name), 'w') as fh:
        fh.write(content)
    #return get_gpx(request, file_name)
    response = HttpResponse('{}', content_type='text/plain')
    response['Access-Control-Allow-Origin'] = '*'
    return response

def create_track_point_xml(track_point):
    return """      <trkpt lat="{latitude}" lon="{longitude}">
      </trkpt>""".format(latitude=track_point[0], longitude=track_point[1])

def create_pgx_file(name, date, description, track_points):
    """ TODO author, trk.type trk.src
    """
    if not name or not date or len(track_points) == 0:
        raise Exception('missing data: ' + name + ', ' + date + ', ' + str(track_points))
    track_points_xml = "\n".join([create_track_point_xml(track_point) for track_point in track_points])
    result = """<?xml version="1.0" encoding="UTF-8"?>
<gpx creator="maps0" version="1.1" xmlns="http://www.topografix.com/GPX/1/1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd" xmlns:gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3">
  <metadata>
    <name>{name}</name>
    <desc>{description}</desc>
    <time>{date}</time>
  </metadata>
  <trk>
    <name>{name}</name>
    <desc>{description}</desc>
    <trkseg>
{track_points_xml}
    </trkseg>
  </trk>
</gpx>""".format(name=name, date=date, description=description, track_points_xml=track_points_xml)
    return result
