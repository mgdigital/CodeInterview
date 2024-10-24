import { Client, Params, Result } from '@/lib/types';

export class ClientImpl implements Client {
  constructor(private baseUrl: string) {}

  async request(params: Params): Promise<Result> {
    const limit = params.limit ?? 10;
    const page = params.page ?? 1;
    const offset = page === 1 ? undefined : (page - 1) * limit;
    const q = params.query || undefined;
    const order = params.order;
    const desc = params.desc ? true : undefined;
    const url =
      this.baseUrl +
      '?' +
      Object.entries({
        limit,
        offset,
        q,
        order,
        desc,
      })
        .flatMap(([k, v]) =>
          v === undefined ? [] : [[k, encodeURIComponent(v)].join('=')],
        )
        .join('&');
    return fetch(url)
      .then((response): Promise<Result> => response.json())
      .catch((err) => {
        // eslint-disable-next-line no-console
        console.error('Error fetching data:', err);
        throw err;
      });
  }
}

export const client: Client = new ClientImpl('http://localhost:8080/assets');
