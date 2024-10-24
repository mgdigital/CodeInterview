import {
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
} from 'lucide-react';

import IconButton from '@/components/buttons/IconButton';
import { Tooltip } from '@/components/Tooltip';

const formatNumber = Intl.NumberFormat('en-us').format;

export const Paginator = ({
  page,
  limit,
  totalCount,
  setPage,
  setLimit,
}: {
  page: number;
  limit: number;
  totalCount: number;
  setPage: (page: number) => void;
  setLimit: (limit: number) => void;
}) => {
  const endIndex = Math.min(page * limit, totalCount);
  const startIndex = endIndex === 0 ? 0 : (page - 1) * limit + 1;
  const pageCount = Math.ceil(totalCount / limit);
  const hasPreviousPage = page > 1;
  const hasNextPage = page < pageCount;

  return (
    <div className='float-right'>
      <span>
        {formatNumber(startIndex)} – {formatNumber(endIndex)} of{' '}
        {formatNumber(totalCount)}
      </span>
      <select
        className='ml-6'
        value={limit}
        onChange={(event) => setLimit(Number(event.target.value))}
      >
        {[10, 20, 50, 100, 500, 1000].map((value) => (
          <option key={value} value={value}>
            {value}
          </option>
        ))}
      </select>
      <span className='ml-2 mr-4'>per page</span>
      <PaginatorTooltip text='First page'>
        <IconButton
          data-testid='paginator-first'
          variant={hasPreviousPage ? 'outline' : 'ghost'}
          disabled={!hasPreviousPage}
          onClick={() => setPage(1)}
          icon={ChevronsLeft}
        />
      </PaginatorTooltip>
      <PaginatorTooltip text='Previous page'>
        <IconButton
          data-testid='paginator-previous'
          variant={hasPreviousPage ? 'outline' : 'ghost'}
          disabled={!hasPreviousPage}
          onClick={() => setPage(page - 1)}
          icon={ChevronLeft}
        />
      </PaginatorTooltip>
      <PaginatorTooltip text='Next page'>
        <IconButton
          data-testid='paginator-next'
          variant={hasNextPage ? 'outline' : 'ghost'}
          disabled={!hasNextPage}
          onClick={() => setPage(page + 1)}
          icon={ChevronRight}
        />
      </PaginatorTooltip>
      <PaginatorTooltip text='Last page'>
        <IconButton
          data-testid='paginator-last'
          variant={hasNextPage ? 'outline' : 'ghost'}
          disabled={!hasNextPage}
          onClick={() => setPage(pageCount)}
          icon={ChevronsRight}
        />
      </PaginatorTooltip>
    </div>
  );
};

const PaginatorTooltip = ({
  children,
  text,
}: {
  children: React.ReactNode;
  text: string;
}) => (
  <Tooltip
    className='inline-block ml-2'
    tooltipChildren={<span className='text-sm'>{text}</span>}
    tooltipClassName='-left-4 top-full mt-3 text-nowrap'
  >
    {children}
  </Tooltip>
);
