const protobuf = require('protobufjs');
const descriptor = require('protobufjs/ext/descriptor');
const fetch = require('node-fetch');

/**
 * Stencil Client options
 * @typedef {Object} Options
 * @property {boolean} [shouldRefresh] - Boolean flag to enable or disable descriptor auto refresh
 * @property {number} [refreshInterval] - interval duration in seconds for refreshing the descriptors
 * @property {Object} [HTTPOptions] - HTTP Options for passing extra information while sending a request. Available options are https://www.npmjs.com/package/node-fetch#options
 * @property {Object.<string, string>} [HTTPOptions.headers] - headers to add while downloading descriptor file
 * @property {Object.<string, string>} [HTTPOptions.timeout] - req/res timeout in ms, it resets on redirect. 0 to disable (OS limit applies).
 */

function checkStatus(res) {
  if (res.ok) {
    return res;
  }
  throw new Error('Unable to download descriptor file');
}
/**
 * Stencil Client class
 */
class Stencil {
  /**
   *
   * @param {string} url
   * @param {Options} options - Options for stencil client
   */
  constructor(url, options) {
    this.url = url;
    this.options = options;
    this.HTTPoptions = { method: 'GET', ...options.HTTPOptions };
    this.root = {};
  }

  async init() {
    if (this.options.shouldRefresh) {
      this.timer = setInterval(
        () =>
          this.load().catch((e) => {
            // eslint-disable-next-line no-console
            console.error(`refresh failed: ${e}`);
          }),
        this.options.refreshInterval * 1000
      );
    }
    await this.load();
  }

  /**
   * Clears any active timers if present
   */
  close() {
    clearInterval(this.timer);
  }

  async load() {
    const buffer = await fetch(this.url, this.HTTPoptions)
      .then(checkStatus)
      .then((res) => res.buffer());
    const decodedDescriptor = descriptor.FileDescriptorSet.decode(buffer);
    this.root = protobuf.Root.fromDescriptor(decodedDescriptor);
    this.root.resolveAll();
  }

  /**
   * @param {string} protoName
   * @returns {protobuf.Type}
   */
  getType(protoName) {
    return this.root.lookupType(protoName);
  }

  /**
   *
   * @param {string} url
   * @param {Options} options - Options for stencil client
   * @returns {Stencil}
   */
  static async getInstance(url, options) {
    const stencil = new Stencil(url, options);
    await stencil.init();
    return stencil;
  }
}

module.exports = Stencil;
