import setuptools

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setuptools.setup(
    name="stencil-python-client",
    version="0.0.1",
    author="ODPF",
    author_email="odpf@gmail.com",
    description="Stencil Python client package provides a store to lookup protobuf descriptors and options to keep the protobuf descriptors upto date.",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/odpf/stencil",
    project_urls={
        "Bug Tracker": "https://github.com/odpf/stencil/issues",
    },
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: Apache Software License",
        "Operating System :: OS Independent",
    ],
    package_dir={"": "src"},
    packages=setuptools.find_packages(where="src"),
    python_requires=">=3.6",
)