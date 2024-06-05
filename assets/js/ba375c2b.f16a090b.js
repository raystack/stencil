"use strict";(self.webpackChunkstencil=self.webpackChunkstencil||[]).push([[930],{3905:function(e,t,r){r.d(t,{Zo:function(){return p},kt:function(){return m}});var n=r(7294);function a(e,t,r){return t in e?Object.defineProperty(e,t,{value:r,enumerable:!0,configurable:!0,writable:!0}):e[t]=r,e}function s(e,t){var r=Object.keys(e);if(Object.getOwnPropertySymbols){var n=Object.getOwnPropertySymbols(e);t&&(n=n.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),r.push.apply(r,n)}return r}function o(e){for(var t=1;t<arguments.length;t++){var r=null!=arguments[t]?arguments[t]:{};t%2?s(Object(r),!0).forEach((function(t){a(e,t,r[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(r)):s(Object(r)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(r,t))}))}return e}function i(e,t){if(null==e)return{};var r,n,a=function(e,t){if(null==e)return{};var r,n,a={},s=Object.keys(e);for(n=0;n<s.length;n++)r=s[n],t.indexOf(r)>=0||(a[r]=e[r]);return a}(e,t);if(Object.getOwnPropertySymbols){var s=Object.getOwnPropertySymbols(e);for(n=0;n<s.length;n++)r=s[n],t.indexOf(r)>=0||Object.prototype.propertyIsEnumerable.call(e,r)&&(a[r]=e[r])}return a}var l=n.createContext({}),c=function(e){var t=n.useContext(l),r=t;return e&&(r="function"==typeof e?e(t):o(o({},t),e)),r},p=function(e){var t=c(e.components);return n.createElement(l.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return n.createElement(n.Fragment,{},t)}},g=n.forwardRef((function(e,t){var r=e.components,a=e.mdxType,s=e.originalType,l=e.parentName,p=i(e,["components","mdxType","originalType","parentName"]),g=c(r),m=a,f=g["".concat(l,".").concat(m)]||g[m]||u[m]||s;return r?n.createElement(f,o(o({ref:t},p),{},{components:r})):n.createElement(f,o({ref:t},p))}));function m(e,t){var r=arguments,a=t&&t.mdxType;if("string"==typeof e||a){var s=r.length,o=new Array(s);o[0]=g;var i={};for(var l in t)hasOwnProperty.call(t,l)&&(i[l]=t[l]);i.originalType=e,i.mdxType="string"==typeof e?e:a,o[1]=i;for(var c=2;c<s;c++)o[c]=r[c];return n.createElement.apply(null,o)}return n.createElement.apply(null,r)}g.displayName="MDXCreateElement"},2379:function(e,t,r){r.r(t),r.d(t,{assets:function(){return p},contentTitle:function(){return l},default:function(){return m},frontMatter:function(){return i},metadata:function(){return c},toc:function(){return u}});var n=r(7462),a=r(3366),s=(r(7294),r(3905)),o=["components"],i={},l="Go",c={unversionedId:"clients/go",id:"clients/go",title:"Go",description:"Go Reference",source:"@site/docs/clients/go.md",sourceDirName:"clients",slug:"/clients/go",permalink:"/stencil/docs/clients/go",editUrl:"https://github.com/raystack/stencil/edit/master/docs/docs/clients/go.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Overview",permalink:"/stencil/docs/clients/overview"},next:{title:"Java",permalink:"/stencil/docs/clients/java"}},p={},u=[{value:"Requirements",id:"requirements",level:2},{value:"Installation",id:"installation",level:2},{value:"Usage",id:"usage",level:2},{value:"Creating a client",id:"creating-a-client",level:3},{value:"Get Descriptor",id:"get-descriptor",level:3},{value:"Parse protobuf message.",id:"parse-protobuf-message",level:3},{value:"Serialize data.",id:"serialize-data",level:3},{value:"Enable auto refresh of schemas",id:"enable-auto-refresh-of-schemas",level:3},{value:"Using VersionBasedRefresh strategy",id:"using-versionbasedrefresh-strategy",level:3}],g={toc:u};function m(e){var t=e.components,r=(0,a.Z)(e,o);return(0,s.kt)("wrapper",(0,n.Z)({},g,r,{components:t,mdxType:"MDXLayout"}),(0,s.kt)("h1",{id:"go"},"Go"),(0,s.kt)("p",null,(0,s.kt)("a",{parentName:"p",href:"https://pkg.go.dev/github.com/raystack/stencil/clients/go"},(0,s.kt)("img",{parentName:"a",src:"https://pkg.go.dev/badge/github.com/raystack/stencil/clients/go.svg",alt:"Go Reference"}))),(0,s.kt)("p",null,"Stencil go client package provides a store to lookup protobuf descriptors and options to keep the protobuf descriptors upto date."),(0,s.kt)("p",null,"It has following features"),(0,s.kt)("ul",null,(0,s.kt)("li",{parentName:"ul"},"Deserialize protobuf messages directly by specifying protobuf message name"),(0,s.kt)("li",{parentName:"ul"},"Serialize data by specifying protobuf message name"),(0,s.kt)("li",{parentName:"ul"},"Ability to refresh protobuf descriptors in specified intervals"),(0,s.kt)("li",{parentName:"ul"},"Support to download descriptors from multiple urls")),(0,s.kt)("h2",{id:"requirements"},"Requirements"),(0,s.kt)("ul",null,(0,s.kt)("li",{parentName:"ul"},"go 1.16")),(0,s.kt)("h2",{id:"installation"},"Installation"),(0,s.kt)("p",null,"Use ",(0,s.kt)("inlineCode",{parentName:"p"},"go get")),(0,s.kt)("pre",null,(0,s.kt)("code",{parentName:"pre"},"go get github.com/raystack/stencil/clients/go\n")),(0,s.kt)("p",null,"Then import the stencil package into your own code as mentioned below"),(0,s.kt)("pre",null,(0,s.kt)("code",{parentName:"pre",className:"language-go"},'import stencil "github.com/raystack/stencil/clients/go"\n')),(0,s.kt)("h2",{id:"usage"},"Usage"),(0,s.kt)("h3",{id:"creating-a-client"},"Creating a client"),(0,s.kt)("pre",null,(0,s.kt)("code",{parentName:"pre",className:"language-go"},'import stencil "github.com/raystack/stencil/clients/go"\n\nurl := "http://localhost:8000/v1beta1/namespaces/{test-namespace}/schemas/{schema-name}"\nclient, err := stencil.NewClient([]string{url}, stencil.Options{})\n')),(0,s.kt)("h3",{id:"get-descriptor"},"Get Descriptor"),(0,s.kt)("pre",null,(0,s.kt)("code",{parentName:"pre",className:"language-go"},'import stencil "github.com/raystack/stencil/clients/go"\n\nurl := "http://localhost:8000/v1beta1/namespaces/{test-namespace}/schemas/{schema-name}"\nclient, err := stencil.NewClient([]string{url}, stencil.Options{})\nif err != nil {\n    return\n}\ndesc, err := client.GetDescriptor("google.protobuf.DescriptorProto")\n')),(0,s.kt)("h3",{id:"parse-protobuf-message"},"Parse protobuf message."),(0,s.kt)("pre",null,(0,s.kt)("code",{parentName:"pre",className:"language-go"},'import stencil "github.com/raystack/stencil/clients/go"\n\nurl := "http://localhost:8000/v1beta1/namespaces/{test-namespace}/schemas/{schema-name}"\nclient, err := stencil.NewClient([]string{url}, stencil.Options{})\nif err != nil {\n    return\n}\ndata := []byte("")\nparsedMsg, err := client.Parse("google.protobuf.DescriptorProto", data)\n')),(0,s.kt)("h3",{id:"serialize-data"},"Serialize data."),(0,s.kt)("pre",null,(0,s.kt)("code",{parentName:"pre",className:"language-go"},'import stencil "github.com/raystack/stencil/clients/go"\n\nurl := "http://url/to/proto/descriptorset/file"\nclient, err := stencil.NewClient([]string{url}, stencil.Options{})\nif err != nil {\n    return\n}\ndata := map[string]interface{}{}\nserializedMsg, err := client.Serialize("google.protobuf.DescriptorProto", data)\n')),(0,s.kt)("h3",{id:"enable-auto-refresh-of-schemas"},"Enable auto refresh of schemas"),(0,s.kt)("pre",null,(0,s.kt)("code",{parentName:"pre",className:"language-go"},'import stencil "github.com/raystack/stencil/clients/go"\n\nurl := "http://localhost:8000/v1beta1/namespaces/{test-namespace}/schemas/{schema-name}"\n// Configured to refresh schema every 12 hours\nclient, err := stencil.NewClient([]string{url}, stencil.Options{AutoRefresh: true, RefreshInterval: time.Hours * 12})\nif err != nil {\n    return\n}\ndesc, err := client.GetDescriptor("google.protobuf.DescriptorProto")\n')),(0,s.kt)("h3",{id:"using-versionbasedrefresh-strategy"},"Using VersionBasedRefresh strategy"),(0,s.kt)("pre",null,(0,s.kt)("code",{parentName:"pre",className:"language-go"},'import stencil "github.com/raystack/stencil/clients/go"\n\nurl := "http://localhost:8000/v1beta1/namespaces/{test-namespace}/schemas/{schema-name}"\n// Configured to refresh schema every 12 hours\nclient, err := stencil.NewClient([]string{url}, stencil.Options{AutoRefresh: true, RefreshInterval: time.Hours * 12, RefreshStrategy: stencil.VersionBasedRefresh})\nif err != nil {\n    return\n}\ndesc, err := client.GetDescriptor("google.protobuf.DescriptorProto")\n')),(0,s.kt)("p",null,"Refer to ",(0,s.kt)("a",{parentName:"p",href:"https://pkg.go.dev/github.com/raystack/stencil/clients/go"},"go documentation")," for all available methods and options."))}m.isMDXComponent=!0}}]);