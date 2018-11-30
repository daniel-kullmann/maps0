from django.urls import path

from . import views

urlpatterns = [
    path('<s>/<z>/<x>/<y>.png', views.get_tile, name="get_tile"),
]
