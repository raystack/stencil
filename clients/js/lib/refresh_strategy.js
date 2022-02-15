const fetch = require('node-fetch');

function checkStatus(res) {
  if (res.ok) {
    return res;
  }
  throw new Error('Unable to download descriptor file');
}

const joinPath = (url, path) => {
  if (url.endsWith('/')) {
    return url + path;
  }
  return [url, path].join('/');
};

const longPollingRefresh = () => async (url, options) =>
  fetch(url, options)
    .then(checkStatus)
    .then((res) => res.buffer());

const versionBasedRefresh = () => {
  let prevLatestVersion = 0;
  return async (url, options) => {
    const versionsURL = joinPath(url, 'versions');
    const versions = await fetch(versionsURL, options)
      .then(checkStatus)
      .then((res) => res.json())
      .then((data) => data.versions || []);
    const maxVersion = Math.max(...versions);
    if (!versions.length || maxVersion <= prevLatestVersion) {
      return null;
    }
    const versionedURL = joinPath(versionsURL, maxVersion.toString());
    const buffer = await fetch(versionedURL, options)
      .then(checkStatus)
      .then((res) => res.buffer());
    prevLatestVersion = maxVersion;
    return buffer;
  };
};

module.exports = {
  longPollingRefresh,
  versionBasedRefresh
};
