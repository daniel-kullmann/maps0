from django.http import HttpResponse
from django.shortcuts import render
from django.views.decorators.csrf import ensure_csrf_cookie
from django import template
from . import models

import json

def get_csrf_token(request):
    json = '{"token": "' + request.COOKIES['csrftoken'] + '"}'
    response = HttpResponse(json, content_type='application/json')
    response['Access-Control-Allow-Origin'] = '*'
    return response

def get_all_settings(request):
    settings = models.Setting.objects.all()
    settings = {s.name: s.value for s in settings}
    response = HttpResponse(json.dumps(settings), content_type='application/json')
    response['Access-Control-Allow-Origin'] = '*'
    return response

def set_all_settings(request):
    for key in request.POST:
        if key == 'csrfmiddlewaretoken': continue
        models.Setting(name=key, value=request.POST[key]).save()
    response = HttpResponse(json.dumps([]), content_type='application/json')
    return response
