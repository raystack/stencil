## [v0.1.4](https://github.com/odpf/stencil/compare/v0.1.3...v0.1.4) (2021-12-07)


### ⚠ BREAKING CHANGES

* **server:** remove cloud storage support as backend store

### Features

* add caching ([#92](https://github.com/odpf/stencil/issues/92)) ([c6f4f34](https://github.com/odpf/stencil/commit/c6f4f34b9db3ee689f49b6717da9ffc7b0d8310c))
* add cli commands ([#53](https://github.com/odpf/stencil/issues/53)) ([78502c1](https://github.com/odpf/stencil/commit/78502c1936b0295e5b7aee408dc09c438120dd73))
* add CNCF spec compatible apis ([#87](https://github.com/odpf/stencil/issues/87)) ([958d3b9](https://github.com/odpf/stencil/commit/958d3b911b6b6b6e6e6cf74a9ff69ba81f25f8db))
* add db abstraction to support multiple dbs ([#59](https://github.com/odpf/stencil/issues/59)) ([03104ce](https://github.com/odpf/stencil/commit/03104ce8e4c25cd251054784331fcd90e1ebdebe))
* add dev command grouping ([ad2189f](https://github.com/odpf/stencil/commit/ad2189fb1862d11d2c7f2d9a4d8bde81bf44ad05))
* add merge descriptor feature providing backward compatibility ([#58](https://github.com/odpf/stencil/issues/58)) ([90a8931](https://github.com/odpf/stencil/commit/90a8931a46adff3091ad92897b7d781f4dd3ae18))
* add protobuf compatibility checks ([#88](https://github.com/odpf/stencil/issues/88)) ([62a5555](https://github.com/odpf/stencil/commit/62a55556843f3780451ee6a0cc2c94fcc93af88f))
* add search API ([#72](https://github.com/odpf/stencil/issues/72)) ([1cf4b58](https://github.com/odpf/stencil/commit/1cf4b5891bcdca0e936eca988df948a58f436a2c))
* add support for visualisation of proto descriptor file dependencies ([#71](https://github.com/odpf/stencil/issues/71)) ([fe5bbf0](https://github.com/odpf/stencil/commit/fe5bbf0e365a33d4539de28316f90a9d9d09a528))
* add type API ([#82](https://github.com/odpf/stencil/issues/82)) ([b68d488](https://github.com/odpf/stencil/commit/b68d4886ecd86c7e7011eff53717ca6b8b88becf))
* embed database migrations ([#55](https://github.com/odpf/stencil/issues/55)) ([dd5e8c4](https://github.com/odpf/stencil/commit/dd5e8c409068d5d67fb060973d8ad5bac8c150f8))
* print cmd for uploaded snapshots ([#60](https://github.com/odpf/stencil/issues/60)) ([5aca993](https://github.com/odpf/stencil/commit/5aca99346d486ad4126a2c671060881227c69123))
* support different content types based on schema type ([#91](https://github.com/odpf/stencil/issues/91)) ([1fc405e](https://github.com/odpf/stencil/commit/1fc405e4edcdf61c7a03b4920fff3da4bcfd4cbb))
* **cmd:** add custom help and grouping for commands ([#57](https://github.com/odpf/stencil/issues/57)) ([837b40f](https://github.com/odpf/stencil/commit/837b40f7c75054b17ffc773a1b243ce0fde8aa37))
* **go-client:** add Serialize() method to serialize data to protobuf ([#62](https://github.com/odpf/stencil/issues/62)) ([f81115d](https://github.com/odpf/stencil/commit/f81115d4561a52c3e2a0a9db8206b1864be953f3))
* provide flag to load server config ([#54](https://github.com/odpf/stencil/issues/54)) ([d2342aa](https://github.com/odpf/stencil/commit/d2342aa5d5ce5a359daa380f3b6da49dabdd916a))
* support gRPC and http apis ([#41](https://github.com/odpf/stencil/issues/41)) ([90a5f25](https://github.com/odpf/stencil/commit/90a5f25d0ec01705f9725db2ff6226d598d3c71d))
* **server:** add docker-compose for dev setup ([a8ab18c](https://github.com/odpf/stencil/commit/a8ab18ccaf6112a43ea69c6ddb6e28cbd97edd3f))
* **server:** add migrate command ([d986f36](https://github.com/odpf/stencil/commit/d986f36858eed2f8de7282e0b85552a5f74da8b8))
* **server:** add postgres data store ([1cda77d](https://github.com/odpf/stencil/commit/1cda77da30cbd7c949d4728ff809aa5865f22b0c))


### Bug Fixes

* **server:** handle no data on download api ([079cc82](https://github.com/odpf/stencil/commit/079cc825444be95e21531f6c68ef3d362f308984))
* **server:** upgrade gin version ([fc5a2fa](https://github.com/odpf/stencil/commit/fc5a2fa795e2c07fdecf17597dc7c0f83114726a))


### Code Refactoring

* **server:** remove cloud storage support as backend store ([d19aad8](https://github.com/odpf/stencil/commit/d19aad884ce02a2b7cc1ea914d88b88e0144708c))

## [v0.1.3](https://github.com/odpf/stencil/compare/v0.1.2...v0.1.3) (2021-06-29)


### Features

* **java-client:** replace java-statsd-client dependency with java-dogstatsd-client ([fed8f35](https://github.com/odpf/stencil/commit/fed8f3512465417187dd01ed262a2929c87491a9))
* **server:** add dryrun flag for upload API ([7e697ab](https://github.com/odpf/stencil/commit/7e697abbc25a12e674f43ccf7fd2c5d286bea72d))


### Bug Fixes

* **java-client:** expose protobuf lib as transitive dependency for consumers ([22b410d](https://github.com/odpf/stencil/commit/22b410da54cc75c2527d42548f0595a96dffeb64))
* **java-client:** expose statsd lib as transitive dependency for consumers ([9bd2108](https://github.com/odpf/stencil/commit/9bd21082490a9a080f85e859850daef738a60026))
* **java-client:** specify type in stencilUtils ([f33a868](https://github.com/odpf/stencil/commit/f33a868d95f8f369cf402a342e7c5d480a25af86))
* **server:** check message and enum names for enum and message field types ([ce328d8](https://github.com/odpf/stencil/commit/ce328d829b76a83e2fc93679127030220a10cd33))

## [v0.1.2](https://github.com/odpf/stencil/compare/v0.1.1...v0.1.2) (2021-06-03)


### Features

* **server:** make graceful shutdown period configurable ([4a11332](https://github.com/odpf/stencil/commit/4a1133201ea865d6940e1e9007d5a8cabd70f241))


### Bug Fixes

* **server:** newrelic config not working with envs ([e29f048](https://github.com/odpf/stencil/commit/e29f0488aeb972f9511cc93a9eaeb580360fbd99))

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

