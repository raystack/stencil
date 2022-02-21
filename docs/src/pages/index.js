import React from 'react';
import Layout from '@theme/Layout';
import clsx from 'clsx';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Container from '../core/Container';
import GridBlock from '../core/GridBlock';
import useBaseUrl from '@docusaurus/useBaseUrl';

const Hero = () => {
  const { siteConfig } = useDocusaurusContext();
  return (
    <div className="homeHero">
      <div className="logo"><img src={useBaseUrl('img/pattern.svg')} /></div>
      <div className="container banner">
        <div className="row">
          <div className={clsx('col col--5')}>
            <div className="homeTitle">{siteConfig.tagline}</div>
            <small className="homeSubTitle">Manage, enforce, and evolve schemas for your event-driven applications.</small>
            <a className="button" href="docs/introduction">Documentation</a>
          </div>
          <div className={clsx('col col--1')}></div>
          <div className={clsx('col col--6')}>
            <div className="text--right"><img src={useBaseUrl('img/banner.svg')} /></div>
          </div>
        </div>
      </div>
    </div >
  );
};

export default function Home() {
  const { siteConfig } = useDocusaurusContext();
  return (
    <Layout
      title={siteConfig.tagline}
      description="Stencil is a schema registry that provides schema mangement and validation to ensure data
      compatibility across applications.">
      <Hero />
      <main>
        <Container className="textSection wrapper" background="light">
          <h1>Built for scale</h1>
          <p>
            Stencil is a schema registry that provides schema mangement and validation to ensure data
            compatibility across applications. It enables developers to create, manage and consume
            schemas dynamically, efficiently, and reliably, and provides a simple way to validate data
            against those schemas.
          </p>
          <GridBlock
            layout="threeColumn"
            contents={[
              {
                title: 'Version history',
                content: (
                  <div>
                    Stencil stores versioned history of proto descriptor file on specified namespace and name.
                  </div>
                ),
              },
              {
                title: 'Backward compatibility',
                content: (
                  <div>
                    Enforce backward compatability check on upload by default.
                  </div>
                ),
              },
              {
                title: 'Flexbility',
                content: (
                  <div>
                    Ability to skip some of the backward compatability checks while upload.
                  </div>
                ),
              },
              {
                title: 'Descriptor fetch',
                content: (
                  <div>
                    Ability to download proto descriptor files.
                  </div>
                ),
              },
              {
                title: 'Metadata',
                content: (
                  <div>
                    Provides metadata API to retrieve latest version number given a name and namespace.
                  </div>
                ),
              },
              {
                title: 'Clients',
                content: (
                  <div>
                    Stencil provides clients in GO, JAVA, JS languages to interact with Stencil server
                    and deserialize messages using dynamic schema.
                  </div>
                ),
              },
            ]}
          />
        </Container>
        <Container className="textSection wrapper" background="light">
          <h1>Trusted by</h1>
          <p>
            Meteor was originally created for the Gojek data processing platform,
            and it has been used, adapted and improved by other teams internally and externally.
          </p>
          <GridBlock className="logos"
            layout="fourColumn"
            contents={[
              {
                content: (
                  <img src={useBaseUrl('users/gojek.png')} />
                ),
              },
              {
                content: (
                  <img src={useBaseUrl('users/midtrans.png')} />
                ),
              },
              {
                content: (
                  <img src={useBaseUrl('users/mapan.png')} />
                ),
              },
              {
                content: (
                  <img src={useBaseUrl('users/moka.png')} />
                ),
              },
            ]}>
          </GridBlock>
        </Container>
      </main>
    </Layout >
  );
}
