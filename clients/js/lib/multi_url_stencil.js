const Stencil = require('./stencil');

/**
 * @typedef {import("./stencil").Options} Options
 */

/**
 * MultiURLStencil Client class
 */
class MultiURLStencil {
  /**
   *
   * @param {string[]} urls
   * @param {Options} options - Options for stencil client
   */
  constructor(urls, options) {
    this.urls = urls;
    this.options = options;
    this.clients = [];
  }

  async init() {
    this.clients = await Promise.all(
      this.urls.map((url) => Stencil.getInstance(url, this.options))
    );
  }

  /**
   * Clears any active timers if present
   */
  close() {
    this.clients.forEach((client) => client.close());
  }

  /**
   * @param {string} protoName
   * @returns {protobuf.Type}
   */
  getType(protoName) {
    let proto;
    for (let i = 0; i < this.clients.length; i += 1) {
      const client = this.clients[i];
      try {
        proto = client.getType(protoName);
        return proto;
      } catch (e) {
        // do nothing
      }
    }
    throw new Error(`no such type: ${protoName}`);
  }

  /**
   *
   * @param {string[]} urls
   * @param {Options} options - Options for stencil client
   */
  static async getInstance(urls, options) {
    const stencil = new MultiURLStencil(urls, options);
    await stencil.init();
    return stencil;
  }
}

module.exports = MultiURLStencil;
