const lightCodeTheme = require('prism-react-renderer/themes/dracula');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

// With JSDoc @type annotations, IDEs can provide config autocompletion
/** @type {import('@docusaurus/types').DocusaurusConfig} */
(module.exports = {
  title: 'Stencil',
  tagline: 'Schema registry engineered for scale',
  url: 'https://raystack.github.io/',
  baseUrl: '/stencil/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/favicon.ico',
  organizationName: 'Raystack',
  projectName: 'stencil',

  presets: [
    [
      '@docusaurus/preset-classic',
      ({
        gtag: {
          trackingID: 'G-XXX',
        },
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          editUrl: 'https://github.com/raystack/stencil/edit/master/docs/',
          sidebarCollapsed: true,
        },
        blog: {
          showReadingTime: true,
          editUrl:
            'https://github.com/raystack/stencil/edit/master/docs/blog/',
        },
        theme: {
          customCss: [
            require.resolve('./src/css/theme.css'),
            require.resolve('./src/css/custom.css'),
            require.resolve('./src/css/icons.css'),
          ],
        },
      })
    ],
  ],

  themeConfig:
    ({
      colorMode: {
        defaultMode: 'light',
        respectPrefersColorScheme: true,
      },
      navbar: {
        title: 'Stencil',
        logo: { src: 'img/logo.svg', },
        hideOnScroll: true,
        items: [
          {
            type: 'doc',
            docId: 'introduction',
            position: 'left',
            label: 'Docs',
          },
          { to: '/blog', label: 'Blog', position: 'left' },
          { to: '/help', label: 'Help', position: 'left' },
          {
            href: 'https://bit.ly/2RzPbtn',
            position: 'right',
            className: 'header-slack-link',
          },
          {
            href: 'https://github.com/raystack/stencil',
            className: 'navbar-item-github',
            position: 'right',
          },
        ],
      },
      footer: {
        style: 'light',
        links: [
          {
            title: 'Products',
            items: [
              { label: 'Optimus', href: 'https://github.com/raystack/optimus' },
              { label: 'Firehose', href: 'https://github.com/raystack/firehose' },
              { label: 'Raccoon', href: 'https://github.com/raystack/raccoon' },
              { label: 'Dagger', href: 'https://raystack.github.io/dagger/' },
            ],
          },
          {
            title: 'Resources',
            items: [
              { label: 'Docs', to: '/docs/introduction' },
              { label: 'Blog', to: '/blog', },
              { label: 'Help', to: '/help', },
            ],
          },
          {
            title: 'Community',
            items: [
              { label: 'Slack', href: 'https://bit.ly/2RzPbtn' },
              { label: 'GitHub', href: 'https://github.com/raystack/stencil' }
            ],
          },
        ],
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
      announcementBar: {
        id: 'star-repo',
        content: '⭐️ If you like Stencil, give it a star on <a target="_blank" rel="noopener noreferrer" href="https://github.com/raystack/stencil">GitHub</a>! ⭐',
        backgroundColor: '#222',
        textColor: '#eee',
        isCloseable: true,
      },
    }),
});
