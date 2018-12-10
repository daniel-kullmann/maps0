from django.db import models

class Setting(models.Model):
    name = models.CharField(max_length=255, primary_key=True)
    value = models.CharField(max_length=255)

    allowed_names = [
        'base_tile_url',
    ]

