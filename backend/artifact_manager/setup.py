from setuptools import setup, find_packages

setup(
    name='odh/data-science-pipelines-artifact-manager',
    version='1.0.0',
    url='https://github.com/opendatahub-io/data-science-pipelines.git',
    author='Open Data Hub',
    author_email='users@lists.opendatahub.io',
    description='Artifact Manager for Data Science Pipelines',
    packages=find_packages(),
    install_requires=['awscli >= 1.25.78']
)
