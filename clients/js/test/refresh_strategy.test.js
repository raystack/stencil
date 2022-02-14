const fetchMock = require('node-fetch');
const {
  longPollingRefresh,
  versionBasedRefresh
} = require('../lib/refresh_strategy');

jest.mock('node-fetch');

describe('RefreshStrategy', () => {
  const url = 'http://localhost:8000/v1beta1/namespaces/a/schemas/b';
  describe('LongPollingRefresh', () => {
    test('should throw error if download fails', async () => {
      fetchMock.mockResolvedValue({ ok: false });
      const refresh = longPollingRefresh();
      expect(refresh(url, {})).rejects.toThrow();
    });
  });
  describe('VersionBasedRefresh', () => {
    beforeEach(() => {
      fetchMock.mockReset();
    });
    test('should throw error if download fails', async () => {
      fetchMock.mockResolvedValue({ ok: false });
      const refresh = versionBasedRefresh();
      expect(refresh(url, {})).rejects.toThrow(
        'Unable to download descriptor file'
      );
    });
    test('should download schema using versioned URL', async () => {
      const fn = jest.fn();
      fn.mockResolvedValue({ versions: [1] });
      const bufferFn = jest.fn();
      bufferFn.mockResolvedValue('some data');
      fetchMock
        .mockResolvedValueOnce({ ok: true, json: fn })
        .mockResolvedValueOnce({ ok: true, buffer: bufferFn });
      const refresh = versionBasedRefresh();
      const val = await refresh(url, {});
      expect(val).toEqual('some data');
      expect(fetchMock).toHaveBeenCalledWith(`${url}/versions`, {});
    });
    test('should not download schema if versions not changed', async () => {
      const trailingUrl = `${url}/`;
      const fn = jest.fn();
      fn.mockResolvedValue({ versions: [1] });
      const bufferFn = jest.fn();
      bufferFn.mockResolvedValue('some data');
      fetchMock
        .mockResolvedValueOnce({ ok: true, json: fn })
        .mockResolvedValueOnce({ ok: true, buffer: bufferFn })
        .mockResolvedValueOnce({ ok: true, json: fn });
      const refresh = versionBasedRefresh();
      let val = await refresh(trailingUrl, {});
      expect(val).toEqual('some data');
      expect(fetchMock).toHaveBeenCalledWith(`${url}/versions`, {});
      val = await refresh(trailingUrl, {});
      expect(val).toBe(null);
    });
    test('should not download schema if versions not changed', async () => {
      const fn = jest.fn();
      fn.mockResolvedValueOnce({ versions: [1] }).mockResolvedValueOnce({
        versions: [1, 2]
      });
      const bufferFn = jest.fn();
      bufferFn.mockResolvedValue('some data');
      fetchMock
        .mockResolvedValueOnce({ ok: true, json: fn })
        .mockResolvedValueOnce({ ok: true, buffer: bufferFn })
        .mockResolvedValueOnce({ ok: true, json: fn })
        .mockResolvedValueOnce({ ok: true, buffer: bufferFn });
      const refresh = versionBasedRefresh();
      let val = await refresh(url, {});
      expect(fetchMock).toHaveBeenCalledWith(`${url}/versions`, {});
      expect(val).toEqual('some data');
      expect(fetchMock).toHaveBeenCalledWith(`${url}/versions/1`, {});
      val = await refresh(url, {});
      expect(val).toEqual('some data');
      expect(fetchMock).toHaveBeenCalledWith(`${url}/versions/2`, {});
    });
  });
});
