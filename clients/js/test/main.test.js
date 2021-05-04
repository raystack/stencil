const fetchMock = require('node-fetch');
const { exec } = require('child_process');
const fs = require('fs');
const main = require('../main');

jest.mock('node-fetch');

const flushPromises = () => new Promise((resolve) => setImmediate(resolve));
function runProtoc(name, includeImports) {
  return new Promise((res) => {
    exec(
      `protoc --descriptor_set_out=./test/data/${name}.desc ${
        includeImports ? '--include_imports' : ''
      } ./test/data/*.proto`,
      res
    );
  });
}

let dataWithNoImports;
let dataWithImports;

const originalErr = console.error;

beforeAll(async () => {
  console.error = jest.fn();
  await runProtoc('filenoinclude', false);
  await runProtoc('file', true);
  dataWithNoImports = fs.readFileSync('./test/data/filenoinclude.desc');
  dataWithImports = fs.readFileSync('./test/data/file.desc');
});

afterAll(() => {
  console.error = originalErr;
  fs.unlinkSync('./test/data/filenoinclude.desc');
  fs.unlinkSync('./test/data/file.desc');
});

describe('Stencil', () => {
  test('should return error if downloadfails', async () => {
    fetchMock.mockResolvedValue({ ok: false });
    expect(main.Stencil.getInstance('http://exampleurl', {})).rejects.toThrow(
      'Unable to download descriptor file'
    );
  });
  test('should return error if downloaded descriptor not valid', async () => {
    const fn = jest.fn();
    fn.mockResolvedValue(Buffer.from('invalid', 'utf-8'));
    fetchMock.mockResolvedValue({ ok: true, buffer: fn });
    expect(main.Stencil.getInstance('http://exampleurl', {})).rejects.toThrow();
  });
  test('should return error if downloaded descriptor not fully contained file', async () => {
    const fn = jest.fn();
    fn.mockResolvedValue(dataWithNoImports);
    fetchMock.mockResolvedValue({ ok: true, buffer: fn });
    expect(main.Stencil.getInstance('http://exampleurl', {})).rejects.toThrow(
      "no such Type or Enum '.google.protobuf.Timestamp' in Type .test.Two.Three"
    );
  });
  test('should return client successfully', async () => {
    const fn = jest.fn();
    fn.mockResolvedValue(dataWithImports);
    fetchMock.mockResolvedValue({ ok: true, buffer: fn });
    expect(
      main.Stencil.getInstance('http://exampleurl', {})
    ).resolves.toBeDefined();
  });

  describe('client', () => {
    let fn;
    beforeEach(() => {
      fn = jest.fn();
      fn.mockResolvedValue(dataWithImports);
      fetchMock.mockResolvedValue({ ok: true, buffer: fn });
      jest.useFakeTimers();
    });
    afterEach(() => {
      jest.clearAllTimers();
    });

    test('should get specified data type', async () => {
      const client = await main.Stencil.getInstance('http://exampleurl', {});
      const type = client.getType('test.Two.Three');
      expect(type.fields.data).toMatchObject({ type: 'string', id: 1 });
      expect(type.fields.timestamp).toMatchObject({
        type: '.google.protobuf.Timestamp',
        id: 3
      });
    });

    test('should throw error if specified data type not found', async () => {
      const client = await main.Stencil.getInstance('http://exampleurl', {});
      expect(() => client.getType('test.Two.Three.Five')).toThrow(
        'no such type: test.Two.Three.Five'
      );
    });

    test('should refresh descriptors on specified refresh interval', async () => {
      const client = await main.Stencil.getInstance('http://exampleurl', {
        shouldRefresh: true,
        refreshInterval: 1
      });
      expect(client).toBeDefined();
      jest.advanceTimersByTime(2000);

      await flushPromises();
      expect(fn).toHaveBeenCalledTimes(3);
    });

    test('client should work if subsequent download fails', async () => {
      fetchMock
        .mockResolvedValueOnce({ ok: true, buffer: fn })
        .mockResolvedValueOnce({ ok: false });

      const client = await main.Stencil.getInstance('http://exampleurl', {
        shouldRefresh: true,
        refreshInterval: 1
      });
      expect(client).toBeDefined();

      jest.advanceTimersByTime(1000);

      await flushPromises();

      expect(client.getType('test.One')).toBeDefined();
    });

    test('should clear timer if close is called', async () => {
      const client = await main.Stencil.getInstance('http://exampleurl', {
        shouldRefresh: true,
        refreshInterval: 1
      });
      expect(client).toBeDefined();
      jest.advanceTimersByTime(2000);
      await flushPromises();
      expect(fn).toHaveBeenCalledTimes(3);
      client.close();
      jest.advanceTimersByTime(2000);
      await flushPromises();
      expect(fn).toHaveBeenCalledTimes(3);
    });
  });
});

describe('MultiURLStencil', () => {
  test('should return error if downloadfails', async () => {
    const fn = jest.fn();
    fn.mockResolvedValue(dataWithImports);
    fetchMock
      .mockResolvedValueOnce({ ok: true, buffer: fn })
      .mockResolvedValue({ ok: false });
    expect(
      main.MultiURLStencil.getInstance(['http://exampleurl', 'anotherurl'], {})
    ).rejects.toThrow('Unable to download descriptor file');
  });
  test('should return error if downloaded descriptor not valid', async () => {
    const fn = jest.fn();
    fn.mockResolvedValue(Buffer.from('invalid', 'utf-8'));
    fetchMock.mockResolvedValue({ ok: true, buffer: fn });
    expect(
      main.MultiURLStencil.getInstance(['http://exampleurl'], {})
    ).rejects.toThrow();
  });
  test('should return error if downloaded descriptor not fully contained file', async () => {
    const fn = jest.fn();
    fn.mockResolvedValue(dataWithNoImports);
    fetchMock.mockResolvedValue({ ok: true, buffer: fn });
    expect(
      main.MultiURLStencil.getInstance(['http://exampleurl'], {})
    ).rejects.toThrow(
      "no such Type or Enum '.google.protobuf.Timestamp' in Type .test.Two.Three"
    );
  });
  test('should return client successfully', async () => {
    const fn = jest.fn();
    fn.mockResolvedValue(dataWithImports);
    fetchMock.mockResolvedValue({ ok: true, buffer: fn });
    expect(
      main.MultiURLStencil.getInstance(['http://exampleurl'], {})
    ).resolves.toBeDefined();
  });

  describe('client', () => {
    let fn;
    beforeEach(() => {
      fn = jest.fn();
      fn.mockResolvedValue(dataWithImports);
      fetchMock.mockResolvedValue({ ok: true, buffer: fn });
      jest.useFakeTimers();
    });
    afterEach(() => {
      jest.clearAllTimers();
    });

    test('should get specified data type', async () => {
      const client = await main.MultiURLStencil.getInstance(
        ['http://exampleurl'],
        {}
      );
      const type = client.getType('test.Two.Three');
      expect(type.fields.data).toMatchObject({ type: 'string', id: 1 });
      expect(type.fields.timestamp).toMatchObject({
        type: '.google.protobuf.Timestamp',
        id: 3
      });
    });

    test('should throw error if specified data type not found', async () => {
      const client = await main.MultiURLStencil.getInstance(
        ['http://exampleurl'],
        {}
      );
      expect(() => client.getType('test.Two.Three.Five')).toThrow(
        'no such type: test.Two.Three.Five'
      );
    });

    test('should refresh descriptors on specified refresh interval', async () => {
      const client = await main.MultiURLStencil.getInstance(
        ['http://exampleurl'],
        { shouldRefresh: true, refreshInterval: 1 }
      );
      expect(client).toBeDefined();
      jest.advanceTimersByTime(2000);

      await flushPromises();
      expect(fn).toHaveBeenCalledTimes(3);
    });

    test('should clear timer if close is called', async () => {
      const client = await main.MultiURLStencil.getInstance(
        ['http://exampleurl'],
        { shouldRefresh: true, refreshInterval: 1 }
      );
      expect(client).toBeDefined();
      jest.advanceTimersByTime(2000);
      await flushPromises();
      expect(fn).toHaveBeenCalledTimes(3);
      client.close();
      jest.advanceTimersByTime(2000);
      await flushPromises();
      expect(fn).toHaveBeenCalledTimes(3);
    });
  });
});
