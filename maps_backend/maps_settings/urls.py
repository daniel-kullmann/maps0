from django.urls import path

from . import views

app_name='settings'
urlpatterns = [
    path('', views.get_all_settings, name='get_all_settings'),
    path('set_all_settings/', views.set_all_settings, name='set_all_settings'),
    path('token/', views.get_csrf_token, name='get_token'),
]
