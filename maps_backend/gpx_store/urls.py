from django.urls import path

from . import views

urlpatterns = [
    path('', views.list_gpx, name="list_gpx"),
    path('get/<name>', views.get_gpx, name="get_gpx"),
]
