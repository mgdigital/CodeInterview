import natsort from 'natsort';
import Highlighter from 'react-highlight-words';

import { Asset, Params, SortOrder } from '@/lib/types';

import Button from '@/components/buttons/Button';
import { Tooltip } from '@/components/Tooltip';

export const Table = ({
  assets,
  params,
  setOrder,
}: {
  assets: Asset[];
  params: Params;
  setOrder: (order: SortOrder, desc: boolean) => void;
}) => (
  <table className='table-auto text-left border-collapse w-full mb-2 border-gray-400 border-b-2'>
    <thead>
      <tr className='bg-gray-200 border-gray-400 border-b-2'>
        <Th>
          <OrderControl
            label='ID'
            column='id'
            params={params}
            setOrder={setOrder}
          />
        </Th>
        <Th>
          <OrderControl
            label='Host'
            column='host'
            params={params}
            setOrder={setOrder}
          />
        </Th>
        <Th>IPs</Th>
        <Th>Ports</Th>
        <Th>
          <OrderControl
            label='Owner'
            column='owner'
            params={params}
            setOrder={setOrder}
          />
        </Th>
        <Th>Comment</Th>
      </tr>
    </thead>
    <tbody>
      {assets.map((a, i) => (
        <tr key={a.id} className={i % 2 ? 'bg-gray-100' : 'bg-white'}>
          <Td>{a.id}</Td>
          <Td className='font-bold group relative'>
            <Tooltip
              tooltipChildren={
                <span>
                  Signature: <span className='font-normal'>{a.signature}</span>
                </span>
              }
              tooltipClassName='left-0 top-1 translate-y-full'
            >
              <Highlighter
                searchWords={params.query ? [params.query] : []}
                autoEscape={true}
                textToHighlight={a.host}
                className='cursor-default'
              />
            </Tooltip>
          </Td>
          <Td>
            {[...a.ips]
              .sort((a, b) => natsort()(a.address, b.address))
              .map((ip, i) => (
                <span key={ip.address}>
                  <Tooltip
                    className='inline-block cursor-default'
                    tooltipChildren={
                      <span>
                        <span className='font-bold'>Signature:</span>
                        <br />
                        {ip.signature}
                      </span>
                    }
                    tooltipClassName='left-0 top-7'
                  >
                    <span className='hover:underline'>{ip.address}</span>
                  </Tooltip>
                  {i === a.ips.length - 1 ? '' : <br />}
                </span>
              ))}
          </Td>
          <Td>
            {(a.ports || []).map((port, i) => (
              <span key={port.port}>
                <Tooltip
                  className='inline-block cursor-default'
                  tooltipChildren={
                    <span>
                      <span className='font-bold'>Signature:</span>
                      <br />
                      {port.signature}
                    </span>
                  }
                  tooltipClassName='left-0 top-7'
                >
                  <span className='hover:underline'>{port.port}</span>
                </Tooltip>
                {i === a.ports.length - 1 ? '' : ', '}
              </span>
            ))}
          </Td>
          <Td>{a.owner}</Td>
          <Td>{a.comment}</Td>
        </tr>
      ))}
    </tbody>
  </table>
);

const Th = ({ children }: { children: React.ReactNode }) => (
  <th className='px-2 pb-2 pt-1 whitespace-nowrap'>{children}</th>
);

const Td = ({
  children,
  className,
}: {
  children: React.ReactNode;
  className?: string;
}) => <td className={'p-2 ' + className}>{children}</td>;

const OrderControl = ({
  label,
  column,
  params,
  setOrder,
}: {
  label: string;
  column: SortOrder;
  params: Pick<Params, 'order' | 'desc'>;
  setOrder: (order: SortOrder, desc: boolean) => void;
}) => {
  const active = params.order === column;
  const desc = active && params.desc;
  return (
    <span
      className='cursor-pointer group'
      onClick={() => setOrder(column, active && !desc)}
    >
      {label}
      <Button
        size='sm'
        variant='ghost'
        className={
          'ml-2 px-1 py-0 group-hover:bg-white ' +
          (active
            ? 'border-blue-500 text-blue-500'
            : 'border-gray-500 text-gray-500')
        }
      >
        {desc ? '▼' : '▲'}
      </Button>
    </span>
  );
};
