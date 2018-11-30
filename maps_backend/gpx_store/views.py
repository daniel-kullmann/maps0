from django.http import HttpResponse

import json
import os
import os.path

DIR='gpx_store_files'

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
