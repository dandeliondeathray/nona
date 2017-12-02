from setuptools import setup, find_packages
setup(
    name="nonacontrol",
    version="0.1",
    packages=find_packages(),
    py_modules=['control', 'nonacontrolclient'],

    author="Erik Edin",
    description="Control is a WebSocket/JSON interface to Nona.",
    license="Apache License 2.0",
    url="https://github.com/dandeliondeathray/nona/service/control"
)
