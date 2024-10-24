export type IP = {
  address: string;
  signature: string;
};

export type Port = {
  port: number;
  signature: string;
};

export type Asset = {
  id: number;
  host: string;
  comment: string;
  owner: string;
  ips: IP[];
  ports: Port[];
  signature: string;
};

export type Result = {
  assets: Asset[];
  totalCount: number;
};

export type SortOrder = 'id' | 'host' | 'owner';

export type Params = {
  limit?: number;
  page?: number;
  query?: string;
  order?: SortOrder;
  desc?: boolean;
};

export interface Client {
  request(params: Params): Promise<Result>;
}
