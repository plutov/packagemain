## packagemain #27: Automated C4 Diagrams with Structurizr DSL

- https://packagemain.tech/p/software-architecture-diagrams-c4
- https://www.npmjs.com/package/serve
- https://docs.structurizr.com/dsl/language
- https://docs.structurizr.com/cli
- https://github.com/avisi-cloud/structurizr-site-generatr
- https://c4model.com/

Intro into C4 Diagrams: https://c4model.com/

Structurizr and its ecosystem: https://structurizr.com/

Write simple DSL - compile - serve

Write full DSL - compile - serve

Deploy with Github actions

```
brew tap avisi-cloud/tools
brew install structurizr-site-generatr

structurizr-site-generatr --help

mkdir ./build
structurizr-site-generatr generate-site --workspace-file diagram.dsl -o ./build

cd build
npm install -g serve
serve
```
