## [v0.1.1](https://github.com/odpf/stencil/compare/v0.1.0...v0.1.1) (2021-05-23)


### Features

* **server:** add newrelic support ([553fd38](https://github.com/odpf/stencil/commit/553fd38da38477a3e4e460958c12e0e0dff1ff97))


### Bug Fixes

* **server:** comparison of proto file options ([7f6ea6d](https://github.com/odpf/stencil/commit/7f6ea6dde78f4bcfb4f1a8cdd536d261e83b8b25))
* **server:** handle panics from rule checks ([1aa37ad](https://github.com/odpf/stencil/commit/1aa37ade94f428c5ef1147155ae4cf8d0b2682e3))


### Reverts

* Revert "chore: upload descriptors to new stencil service" ([ec1d4b4](https://github.com/odpf/stencil/commit/ec1d4b42c04775b0cd9668e146027c997c64cdb3))

##  (2021-05-11)


### ⚠ BREAKING CHANGES

* **server:** change in API contract. Removed x-scope-orgid header. Moved to namespace param.
New API endpoints starts with `/v1/namespaces/{namespace}`.

### Features

* **go-client:** add default values for options ([49b38a7](https://github.com/odpf/stencil/commit/49b38a777fa1c80c4350b970c7ed3b9b62f13aab))
* **go-client:** add ParseWithRefresh method to GO client ([d9a51e7](https://github.com/odpf/stencil/commit/d9a51e767b1b0f039fe6e78a28bc6b920f15e32a))
* **go-client:** add refresh method to client API ([904fdd4](https://github.com/odpf/stencil/commit/904fdd439680709892860ee858ff464a74e87439))
* **js-client:** add js stencil client ([556825a](https://github.com/odpf/stencil/commit/556825a379e908178a8a34ea91b5652273c0efc0))
* **server:** move x-scope-orgid header to namespace param ([d02ac65](https://github.com/odpf/stencil/commit/d02ac655e15fd928ae675603caa825c1bfb0c305))
* add authentication bearer token used for fetch requests ([da7a6e0](https://github.com/odpf/stencil/commit/da7a6e05297412a650bd4ef749ab615ca140524a))
* add backward compatability check to descriptor create API ([8d1ba88](https://github.com/odpf/stencil/commit/8d1ba886d847f3f7450d7e97ec6c1d645be5dfe6))
* add backward compatability rule checks ([7dfeb4a](https://github.com/odpf/stencil/commit/7dfeb4a6bbea1670007028e87998200087600eae))
* add descriptor APIs ([e8e17cc](https://github.com/odpf/stencil/commit/e8e17ccc8f8087c8c80099e2f01f9baccebcc152))
* add docker release in github action ([dcd2dc5](https://github.com/odpf/stencil/commit/dcd2dc54a1b151af8a12ad3caed4e6137ea1a85c))
* add javadoc comments ([a7ee4a4](https://github.com/odpf/stencil/commit/a7ee4a42ee516208665d8259b0286fa8689c161a))
* add metadata API ([d9499ff](https://github.com/odpf/stencil/commit/d9499ffe5d52f0c15d972e6fe2984f77161bfd0e))
* add skipRules field in descriptor upload API ([4071991](https://github.com/odpf/stencil/commit/4071991b491dcd62a7f3d5e9f12cefe4a9a39e65))
* add stencil go client ([b4c4ed5](https://github.com/odpf/stencil/commit/b4c4ed527816f3b29976e5306f9170c10ea2d525))


### Bug Fixes

* **server:** POST descriptor API to throw error if descriptor already exists ([240a176](https://github.com/odpf/stencil/commit/240a17662f138439c2233aa1d3e9f7fddb701f4f))
* gradle project versioning ([f0d3db7](https://github.com/odpf/stencil/commit/f0d3db72f2cf288e1be2cefe4dc76b1025e066b8))
* handle file not found errors ([df0fd3f](https://github.com/odpf/stencil/commit/df0fd3f5a5094fe5bb44762575254ea3f1e08e74))
* handle uploaded file read error ([a93f4e2](https://github.com/odpf/stencil/commit/a93f4e252cea6c6adb53ec19244bb9c94512d39b))
* hide meta.json in listing ([deebb9c](https://github.com/odpf/stencil/commit/deebb9c67846ec2dc38a037314e8a5e38f472b86))
* incorrect proto name where javapackage or protopackage is empty ([69de5e8](https://github.com/odpf/stencil/commit/69de5e88abeb978cba534050e6ac1153ffd7d8be))
* nil reader error if file is not found ([54ba89d](https://github.com/odpf/stencil/commit/54ba89dc0f5a7b7f2d8918d4dd3ed579e3f5cdd4))


### Reverts

* Revert "nexus-setup: add release step to publish to central maven" ([c748cd7](https://github.com/odpf/stencil/commit/c748cd7e4bea9ab1ef5115af333c2fb430ac5eaa))

