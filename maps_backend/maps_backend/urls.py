from django.contrib import admin
from django.urls import path, include

urlpatterns = [
    path('api/admin/', admin.site.urls),
    path('api/tiles/', include('tiles.urls')),
    path('api/gpx/', include('gpx_store.urls')),
]
