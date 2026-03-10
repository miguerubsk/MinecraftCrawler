# Changelog

## [0.3.0](https://github.com/miguerubsk/MinecraftCrawler/compare/v0.2.0...v0.3.0) (2026-03-10)


### Features

* connect 'info' command to deep analysis pipeline ([0f3d2bb](https://github.com/miguerubsk/MinecraftCrawler/commit/0f3d2bb193796b76b2e857caf78b8ed438f1988e))
* expand data extraction to include MOTD, MapName and detailed server info [#22](https://github.com/miguerubsk/MinecraftCrawler/issues/22) ([508f1cc](https://github.com/miguerubsk/MinecraftCrawler/commit/508f1cc25439ff44550948fc5b05a8fafee10e9a))
* implement rich colored terminal UI for info command [#23](https://github.com/miguerubsk/MinecraftCrawler/issues/23) ([7fe7e2c](https://github.com/miguerubsk/MinecraftCrawler/commit/7fe7e2c3ee118f3c6ce8b66f5cc38460cfb852b8))
* implement robust target parsing and SRV resolution for info command [#21](https://github.com/miguerubsk/MinecraftCrawler/issues/21) ([8756e84](https://github.com/miguerubsk/MinecraftCrawler/commit/8756e84a83283196745f0a9fe6bddbb6f4c62843))
* implement robust target parsing and SRV resolution for info command [#21](https://github.com/miguerubsk/MinecraftCrawler/issues/21) ([e500070](https://github.com/miguerubsk/MinecraftCrawler/commit/e500070b83d7862cda7c4ea121531860dda9726e))
* implement SRV record resolution for info command [#21](https://github.com/miguerubsk/MinecraftCrawler/issues/21) ([161e52b](https://github.com/miguerubsk/MinecraftCrawler/commit/161e52b273ab946b69f6c00d2720f09a81666a97))
* **info:** probe default RCON port during deep analysis and report explicit RCON status ([ba8714e](https://github.com/miguerubsk/MinecraftCrawler/commit/ba8714e8d26a67e345c1aa21918a1882df065c97))
* initial 'info' command and SRV resolution [#21](https://github.com/miguerubsk/MinecraftCrawler/issues/21) ([3340eae](https://github.com/miguerubsk/MinecraftCrawler/commit/3340eae33de511b252a5b2329f8048d974d04a38))
* initialize skeleton for 'info' command [#21](https://github.com/miguerubsk/MinecraftCrawler/issues/21) ([6b930bc](https://github.com/miguerubsk/MinecraftCrawler/commit/6b930bc9a84923cf0f8c67e4d7ee15444379cf93))


### Bug Fixes

* added timeout ([9cff2fb](https://github.com/miguerubsk/MinecraftCrawler/commit/9cff2fbd64e724fa4fedcc311d613103b426e9af))
* **cli,tests:** move output flag to scan and correct query mock stat header padding ([64de5d4](https://github.com/miguerubsk/MinecraftCrawler/commit/64de5d40e5325cd0b50770b258ceba358fd57767))
* close UDP conn, use fresh sub-timeout for login probe [#23](https://github.com/miguerubsk/MinecraftCrawler/issues/23) ([53101bb](https://github.com/miguerubsk/MinecraftCrawler/commit/53101bb26fbf4ad56bab7983e3c10c33d378dbd6))
* correct SetDeadline order, extract separator const and hide raw query errors [#22](https://github.com/miguerubsk/MinecraftCrawler/issues/22) ([1682a99](https://github.com/miguerubsk/MinecraftCrawler/commit/1682a9921a6ba38a976e5470969f57d3c8f9ff03))
* escape backticks in ascii banner, print banner in info and format help examples [#23](https://github.com/miguerubsk/MinecraftCrawler/issues/23) ([a929b73](https://github.com/miguerubsk/MinecraftCrawler/commit/a929b7396e04a970a383b2f8b9910fcd89c78348))
* **formatter:** support §k and uppercase MOTD formatting codes ([adfec69](https://github.com/miguerubsk/MinecraftCrawler/commit/adfec695a9042f1fab18e6d515a8507a7e1b52a9))
* imports ([0b6f6ca](https://github.com/miguerubsk/MinecraftCrawler/commit/0b6f6ca1db619680e96f58a87a165ae530f99c3c))
* **info:** prevent format string injection in server detail rendering ([1df3ee7](https://github.com/miguerubsk/MinecraftCrawler/commit/1df3ee7d2e3ae63db727db2749141f89b28d4f0b))
* lint ([0d4e598](https://github.com/miguerubsk/MinecraftCrawler/commit/0d4e5981a802b2d798b2bd95384d0df336448bfb))
* more robust target parsing and export InfoCmd [#21](https://github.com/miguerubsk/MinecraftCrawler/issues/21) ([3e21497](https://github.com/miguerubsk/MinecraftCrawler/commit/3e2149715bdedfe78926d61f40cd4df0a033932d))
* motd parsing to support rich descriptions and clean up rcon deadlines [#22](https://github.com/miguerubsk/MinecraftCrawler/issues/22) ([3f56a4f](https://github.com/miguerubsk/MinecraftCrawler/commit/3f56a4f79896c5893523eadb7839b1aefe456274))
* motd recursive parsing, ansi color toggle support and prevent double separators [#23](https://github.com/miguerubsk/MinecraftCrawler/issues/23) ([58201a0](https://github.com/miguerubsk/MinecraftCrawler/commit/58201a016b0bdd8f835ac73a15e8b42e26f14e5e))
* motd reset conditional, output routing for info and check deadline error [#23](https://github.com/miguerubsk/MinecraftCrawler/issues/23) ([833502c](https://github.com/miguerubsk/MinecraftCrawler/commit/833502c7e58e9101a2d58d6edb8ac93cc998b40c))
* query payload slice offset [#22](https://github.com/miguerubsk/MinecraftCrawler/issues/22) ([458d850](https://github.com/miguerubsk/MinecraftCrawler/commit/458d850c661f63dfa7682421c6339579142aeb65))
* remove unused variables overallDeadline and out from linter errors ([98a1f4c](https://github.com/miguerubsk/MinecraftCrawler/commit/98a1f4cf6a299c83264acc0771d64bc251f28ff3))
* **storage,test:** persist map_name in sqlite and add parseMOTD internal unit tests ([9caf4e1](https://github.com/miguerubsk/MinecraftCrawler/commit/9caf4e16ccde987ce18e8a3af844721f6e411b47))
* technical refinements to deadlines, data assignment and spanish localization [#22](https://github.com/miguerubsk/MinecraftCrawler/issues/22) ([0a05536](https://github.com/miguerubsk/MinecraftCrawler/commit/0a05536f9856a97a97b5d5ab024bbfc6c7beafb6))

## [0.2.0](https://github.com/miguerubsk/MinecraftCrawler/compare/v0.1.2...v0.2.0) (2026-03-05)


### Features

* add session summary dashboard and refined timestamps to logs ([4bb1243](https://github.com/miguerubsk/MinecraftCrawler/commit/4bb1243ea8af2a247d30cd96b4f783f3ec3bf998))
* enhance CLI aesthetics and real-time scanner feedback ([81a2041](https://github.com/miguerubsk/MinecraftCrawler/commit/81a2041fdf39983c9bacce1591323c367b610545))


### Bug Fixes

* clear terminal line before logging to avoid overlap with masscan progress ([07dcec1](https://github.com/miguerubsk/MinecraftCrawler/commit/07dcec18fe5d75200bbc0a6e47c59580922f1a41))
* unify terminal and file output to ensure all logs are captured ([2d65de2](https://github.com/miguerubsk/MinecraftCrawler/commit/2d65de29e8c377fc4ec46573610e6472a2fb5681))

## [0.1.2](https://github.com/miguerubsk/MinecraftCrawler/compare/v0.1.1...v0.1.2) (2026-02-26)


### Bug Fixes

* restore codeql deleted by mistake ([e265f99](https://github.com/miguerubsk/MinecraftCrawler/commit/e265f99d7cc166abf7ef7b9ab21ffb4f1e3c64a2))

## [0.1.1](https://github.com/miguerubsk/MinecraftCrawler/compare/v0.1.0...v0.1.1) (2026-02-26)


### Bug Fixes

* restore codeql deleted by mistake ([e48debb](https://github.com/miguerubsk/MinecraftCrawler/commit/e48debb2dca552c90e38aadb58b05af1431a59d8))

## [0.1.0](https://github.com/miguerubsk/MinecraftCrawler/compare/v0.0.1...v0.1.0) (2026-02-25)


### Features

* stable release setup ([7a1a49e](https://github.com/miguerubsk/MinecraftCrawler/commit/7a1a49e1583531c25311c81c5f66cd7b04ef78b9))

## [0.0.1](https://github.com/miguerubsk/MinecraftCrawler/compare/v0.0.0...v0.0.1) (2026-02-25)


### Bug Fixes

* address copilot review and fix storage collisions ([57d461f](https://github.com/miguerubsk/MinecraftCrawler/commit/57d461f91118bcbe275a41be8602481e1b749123))
