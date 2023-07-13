"use strict";(self.webpackChunkstencil=self.webpackChunkstencil||[]).push([[217],{3905:function(e,t,n){n.d(t,{Zo:function(){return p},kt:function(){return m}});var a=n(7294);function r(e,t,n){return t in e?Object.defineProperty(e,t,{value:n,enumerable:!0,configurable:!0,writable:!0}):e[t]=n,e}function l(e,t){var n=Object.keys(e);if(Object.getOwnPropertySymbols){var a=Object.getOwnPropertySymbols(e);t&&(a=a.filter((function(t){return Object.getOwnPropertyDescriptor(e,t).enumerable}))),n.push.apply(n,a)}return n}function i(e){for(var t=1;t<arguments.length;t++){var n=null!=arguments[t]?arguments[t]:{};t%2?l(Object(n),!0).forEach((function(t){r(e,t,n[t])})):Object.getOwnPropertyDescriptors?Object.defineProperties(e,Object.getOwnPropertyDescriptors(n)):l(Object(n)).forEach((function(t){Object.defineProperty(e,t,Object.getOwnPropertyDescriptor(n,t))}))}return e}function o(e,t){if(null==e)return{};var n,a,r=function(e,t){if(null==e)return{};var n,a,r={},l=Object.keys(e);for(a=0;a<l.length;a++)n=l[a],t.indexOf(n)>=0||(r[n]=e[n]);return r}(e,t);if(Object.getOwnPropertySymbols){var l=Object.getOwnPropertySymbols(e);for(a=0;a<l.length;a++)n=l[a],t.indexOf(n)>=0||Object.prototype.propertyIsEnumerable.call(e,n)&&(r[n]=e[n])}return r}var s=a.createContext({}),c=function(e){var t=a.useContext(s),n=t;return e&&(n="function"==typeof e?e(t):i(i({},t),e)),n},p=function(e){var t=c(e.components);return a.createElement(s.Provider,{value:t},e.children)},u={inlineCode:"code",wrapper:function(e){var t=e.children;return a.createElement(a.Fragment,{},t)}},d=a.forwardRef((function(e,t){var n=e.components,r=e.mdxType,l=e.originalType,s=e.parentName,p=o(e,["components","mdxType","originalType","parentName"]),d=c(n),m=r,k=d["".concat(s,".").concat(m)]||d[m]||u[m]||l;return n?a.createElement(k,i(i({ref:t},p),{},{components:n})):a.createElement(k,i({ref:t},p))}));function m(e,t){var n=arguments,r=t&&t.mdxType;if("string"==typeof e||r){var l=n.length,i=new Array(l);i[0]=d;var o={};for(var s in t)hasOwnProperty.call(t,s)&&(o[s]=t[s]);o.originalType=e,o.mdxType="string"==typeof e?e:r,i[1]=o;for(var c=2;c<l;c++)i[c]=n[c];return a.createElement.apply(null,i)}return a.createElement.apply(null,n)}d.displayName="MDXCreateElement"},9803:function(e,t,n){n.r(t),n.d(t,{assets:function(){return p},contentTitle:function(){return s},default:function(){return m},frontMatter:function(){return o},metadata:function(){return c},toc:function(){return u}});var a=n(7462),r=n(3366),l=(n(7294),n(3905)),i=["components"],o={},s="Installation",c={unversionedId:"installation",id:"installation",title:"Installation",description:"Stencil installation is simple. You can install Stencil on macOS, Windows, Linux, OpenBSD, FreeBSD, and on any machine. There are several approaches to installing Stencil.",source:"@site/docs/installation.md",sourceDirName:".",slug:"/installation",permalink:"/stencil/docs/installation",editUrl:"https://github.com/raystack/stencil/edit/master/docs/docs/installation.md",tags:[],version:"current",frontMatter:{},sidebar:"docsSidebar",previous:{title:"Usecases",permalink:"/stencil/docs/usecases"},next:{title:"Glossary",permalink:"/stencil/docs/glossary"}},p={},u=[{value:"Binary (Cross-platform)",id:"binary-cross-platform",level:3},{value:"MacOS",id:"macos",level:3},{value:"Linux",id:"linux",level:4},{value:"Windows",id:"windows",level:3},{value:"Docker",id:"docker",level:3},{value:"Building from source",id:"building-from-source",level:3},{value:"Verifying the installation",id:"verifying-the-installation",level:3}],d={toc:u};function m(e){var t=e.components,n=(0,r.Z)(e,i);return(0,l.kt)("wrapper",(0,a.Z)({},d,n,{components:t,mdxType:"MDXLayout"}),(0,l.kt)("h1",{id:"installation"},"Installation"),(0,l.kt)("p",null,"Stencil installation is simple. You can install Stencil on macOS, Windows, Linux, OpenBSD, FreeBSD, and on any machine. There are several approaches to installing Stencil."),(0,l.kt)("ol",null,(0,l.kt)("li",{parentName:"ol"},"Using a ",(0,l.kt)("a",{parentName:"li",href:"#binary-cross-platform"},"pre-compiled binary")),(0,l.kt)("li",{parentName:"ol"},"Installing with ",(0,l.kt)("a",{parentName:"li",href:"#MacOS"},"package manager")),(0,l.kt)("li",{parentName:"ol"},"Installing with ",(0,l.kt)("a",{parentName:"li",href:"#Docker"},"Docker")),(0,l.kt)("li",{parentName:"ol"},"Installing from ",(0,l.kt)("a",{parentName:"li",href:"#building-from-source"},"source"))),(0,l.kt)("h3",{id:"binary-cross-platform"},"Binary (Cross-platform)"),(0,l.kt)("p",null,"Download the appropriate version for your platform from ",(0,l.kt)("a",{parentName:"p",href:"https://github.com/raystack/stencil/releases"},"releases")," page. Once downloaded, the binary can be run from anywhere.\nYou don\u2019t need to install it into a global location. This works well for shared hosts and other systems where you don\u2019t have a privileged account.\nIdeally, you should install it somewhere in your ",(0,l.kt)("inlineCode",{parentName:"p"},"PATH")," for easy use. ",(0,l.kt)("inlineCode",{parentName:"p"},"/usr/local/bin")," is the most probable location."),(0,l.kt)("h3",{id:"macos"},"MacOS"),(0,l.kt)("p",null,(0,l.kt)("inlineCode",{parentName:"p"},"stencil")," is available via a Homebrew Tap, and as downloadable binary from the ",(0,l.kt)("a",{parentName:"p",href:"https://github.com/raystack/stencil/releases/latest"},"releases")," page:"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-sh"},"brew install raystack/tap/stencil\n")),(0,l.kt)("p",null,"To upgrade to the latest version:"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre"},"brew upgrade stencil\n")),(0,l.kt)("h4",{id:"linux"},"Linux"),(0,l.kt)("p",null,(0,l.kt)("inlineCode",{parentName:"p"},"stencil")," is available as downloadable binaries from the ",(0,l.kt)("a",{parentName:"p",href:"https://github.com/raystack/stencil/releases/latest"},"releases")," page. Download the ",(0,l.kt)("inlineCode",{parentName:"p"},".deb")," or ",(0,l.kt)("inlineCode",{parentName:"p"},".rpm")," from the releases page and install with ",(0,l.kt)("inlineCode",{parentName:"p"},"sudo dpkg -i")," and ",(0,l.kt)("inlineCode",{parentName:"p"},"sudo rpm -i")," respectively."),(0,l.kt)("h3",{id:"windows"},"Windows"),(0,l.kt)("p",null,(0,l.kt)("inlineCode",{parentName:"p"},"stencil")," is available via ",(0,l.kt)("a",{parentName:"p",href:"https://scoop.sh/"},"scoop"),", and as a downloadable binary from the ",(0,l.kt)("a",{parentName:"p",href:"https://github.com/raystack/stencil/releases/latest"},"releases")," page:"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre"},"scoop bucket add stencil https://github.com/raystack/scoop-bucket.git\n")),(0,l.kt)("p",null,"To upgrade to the latest version:"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre"},"scoop update stencil\n")),(0,l.kt)("h3",{id:"docker"},"Docker"),(0,l.kt)("p",null,"We provide ready to use Docker container images. To pull the latest image:"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre"},"docker pull raystack/stencil:latest\n")),(0,l.kt)("p",null,"To pull a specific version:"),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre"},"docker pull raystack/stencil:v0.3.3\n")),(0,l.kt)("h3",{id:"building-from-source"},"Building from source"),(0,l.kt)("p",null,"To compile from source, you will need ",(0,l.kt)("a",{parentName:"p",href:"https://golang.org/"},"Go")," installed and a copy of ",(0,l.kt)("a",{parentName:"p",href:"https://www.git-scm.com/"},"git")," in your ",(0,l.kt)("inlineCode",{parentName:"p"},"PATH"),"."),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-bash"},"# Clone the repo\n$ git clone git@github.com:raystack/stencil.git\n\n# Check all build comamnds available\n$ make help\n\n# Build stencil binary file\n$ make build\n\n# Check for installed stencil version\n$ ./stencil version\n")),(0,l.kt)("h3",{id:"verifying-the-installation"},"Verifying the installation"),(0,l.kt)("p",null,"To verify Stencil is properly installed, run ",(0,l.kt)("inlineCode",{parentName:"p"},"stencil --help")," on your system. You should see help output. If you are executing it from the command line, make sure it is on your ",(0,l.kt)("inlineCode",{parentName:"p"},"PATH")," or you may get an error about Stencil not being found."),(0,l.kt)("pre",null,(0,l.kt)("code",{parentName:"pre",className:"language-bash"},"$ stencil --help\n")))}m.isMDXComponent=!0}}]);