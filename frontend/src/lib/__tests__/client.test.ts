import { ClientImpl } from '@/lib/client';
import { Params, Result } from '@/lib/types';

describe('client', () => {
  const mockFetch = jest.fn();
  global.fetch = mockFetch;
  const client = new ClientImpl('http://localhost:8080/assets');

  beforeEach(() => {
    jest.resetAllMocks();
  });

  it.each(
    Array<[Params, Result, string]>(
      [
        {},
        {
          assets: [
            {
              id: 1,
              host: 'foo.com',
              owner: 'Mrs. Test',
              signature: 'test',
              comment: 'test',
              ips: [],
              ports: [],
            },
          ],
          totalCount: 1,
        },
        'http://localhost:8080/assets?limit=10',
      ],
      [
        {
          page: 5,
          limit: 20,
          query: 'foo',
          order: 'host',
          desc: true,
        },
        {
          assets: [
            {
              id: 1,
              host: 'foo.com',
              owner: 'Mrs. Test',
              signature: 'test',
              comment: 'test',
              ips: [],
              ports: [],
            },
          ],
          totalCount: 1,
        },
        'http://localhost:8080/assets?limit=20&offset=80&q=foo&order=host&desc=true',
      ],
    ),
  )('should fetch a list of assets', async (params, result, url) => {
    mockFetch.mockResolvedValueOnce({ json: () => result });
    await expect(client.request(params)).resolves.toEqual(result);
    expect(mockFetch).toHaveBeenCalledWith(url);
  });
});
