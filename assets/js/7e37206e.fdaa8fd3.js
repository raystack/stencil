/*! For license information please see 7e37206e.fdaa8fd3.js.LICENSE.txt */
(self.webpackChunkstencil=self.webpackChunkstencil||[]).push([[740],{2579:function(e,t,n){"use strict";var r=n(4184),i=n.n(r),l=n(7294),o=function(e){var t,n=i()(e.className,{darkBackground:"dark"===e.background,highlightBackground:"highlight"===e.background,lightBackground:"light"===e.background,paddingAll:e.padding.indexOf("all")>=0,paddingBottom:e.padding.indexOf("bottom")>=0,paddingLeft:e.padding.indexOf("left")>=0,paddingRight:e.padding.indexOf("right")>=0,paddingTop:e.padding.indexOf("top")>=0});return t=e.wrapper?l.createElement("div",{className:"container"},e.children):e.children,l.createElement("div",{className:n,id:e.id},t)};o.defaultProps={background:null,padding:[],wrapper:!0},t.Z=o},9260:function(e,t,n){"use strict";var r=n(4578),i=n(4184),l=n.n(i),o=n(7294),s=function(e){function t(){return e.apply(this,arguments)||this}(0,r.Z)(t,e);var n=t.prototype;return n.renderBlock=function(e){var t=Object.assign({},{imageAlign:"left"},e),n=l()("blockElement",this.props.className,{alignCenter:"center"===this.props.align,alignRight:"right"===this.props.align,fourByGridBlock:"fourColumn"===this.props.layout,threeByGridBlock:"threeColumn"===this.props.layout,twoByGridBlock:"twoColumn"===this.props.layout});return o.createElement("div",{className:n,key:t.title},o.createElement("div",{className:"blockContent"},this.renderBlockTitle(t.title),t.content))},n.renderBlockTitle=function(e){return e?o.createElement("h2",null,e):null},n.render=function(){return o.createElement("div",{className:"gridBlock"},this.props.contents.map(this.renderBlock,this))},t}(o.Component);s.defaultProps={align:"left",contents:[],layout:"twoColumn"},t.Z=s},5401:function(e,t,n){"use strict";n.r(t),n.d(t,{default:function(){return a}});var r=n(7294),i=n(2773),l=n(2263),o=n(2579),s=n(9260);function a(){(0,l.Z)().siteConfig;return r.createElement(i.Z,{title:"Help",description:"Comunity help"},r.createElement("main",null,r.createElement(o.Z,{className:"textSection wrapper",background:"light"},r.createElement("h1",null,"Need help?"),r.createElement("p",null,"Need a bit of help? We're here for you. Check out our current issues, GitHub discussions, or get support through Slack."),r.createElement(s.Z,{layout:"threeColumn",contents:[{title:"Slack",content:r.createElement("div",null,"The Stencil team has an open source slack workspace to discuss development and support. Most of the Stencil discussions happen in #stencil channel.",r.createElement("br",null),r.createElement("a",{href:"https://raystack-community.slack.com/"}," Join us on Slack "))},{title:"GitHub Issues",content:r.createElement("div",null,"Have a general issue or bug that you've found? We'd love to hear about it in our GitHub issues. This can be feature requests too!",r.createElement("br",null)," ",r.createElement("a",{target:"_blank",href:"https://github.com/raystack/stencil/issues"}," Go to issues "))},{title:"GitHub Discussions",content:r.createElement("div",null,"For help and questions about best practices, join our GitHub discussions. Browse and ask questions.",r.createElement("br",null),r.createElement("a",{target:"_blank",href:"https://github.com/raystack/stencil/discussions"}," Go to discussions "))}]}))))}},4184:function(e,t){var n;!function(){"use strict";var r={}.hasOwnProperty;function i(){for(var e=[],t=0;t<arguments.length;t++){var n=arguments[t];if(n){var l=typeof n;if("string"===l||"number"===l)e.push(n);else if(Array.isArray(n)){if(n.length){var o=i.apply(null,n);o&&e.push(o)}}else if("object"===l)if(n.toString===Object.prototype.toString)for(var s in n)r.call(n,s)&&n[s]&&e.push(s);else e.push(n.toString())}}return e.join(" ")}e.exports?(i.default=i,e.exports=i):void 0===(n=function(){return i}.apply(t,[]))||(e.exports=n)}()}}]);