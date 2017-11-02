"""
Helper functions for checking readiness and liveness probes in Kubernetes services.
"""
import requests


def http_probe(url):
    try:
        response = requests.get(url)
        return response.status_code == requests.codes.ok
    except:
        return False
