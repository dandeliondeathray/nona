from setuptools import setup, find_packages
setup(
    name="nonastaging",
    version="0.1",
    packages=find_packages(),
    scripts=['nonastaging.py'],

    author="Erik Edin",
    description="Nonastaging is a WebSocket/JSON interface to Nona.",
    license="Apache License 2.0",
    url="https://github.com/dandeliondeathray/nona/test/nonastaging"
)
