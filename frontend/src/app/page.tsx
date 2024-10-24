'use client';

import { Timer, X } from 'lucide-react';
import { debounce } from 'next/dist/server/utils';
import * as React from 'react';
import { useEffect, useRef, useState } from 'react';

import { client } from '@/lib/client';
import { Params, Result, SortOrder } from '@/lib/types';

import IconButton from '@/components/buttons/IconButton';
import { Paginator } from '@/components/Paginator';
import { Table } from '@/components/Table';
import { Tooltip } from '@/components/Tooltip';

import Logo from '~/svg/Logo.svg';

export default function HomePage() {
  const [params, setParams] = useState<Params>({
    limit: 10,
    page: 1,
    query: '',
    order: 'host',
    desc: false,
  });
  const [result, setResult] = useState<Result>({
    assets: [],
    totalCount: 0,
  });
  const [isLoading, setIsLoading] = useState(true);
  const [loadTime, setLoadTime] = useState<number>(0);
  const mountedRef = useRef(true);

  useEffect(() => {
    async function fetchData() {
      const start = performance.now();
      setIsLoading(true);
      const result = await client.request(params);
      // check if a newer request is already in-flight, and if so discard the current result
      if (!mountedRef.current) return null;
      setResult(result);
      setLoadTime(performance.now() - start);
      setIsLoading(false);
    }
    void fetchData();
  }, [params]);

  const setPage = (page: number) => {
    setParams({ ...params, page });
  };

  // these setter methods will redirect to page 1 when changing the query, limit, or order
  const setQuery = debounce((query: string) => {
    setParams({ ...params, page: 1, query });
  }, 300);

  const setLimit = (limit: number) => {
    setParams({ ...params, page: 1, limit });
  };

  const setOrder = (order: SortOrder, desc: boolean) => {
    setParams({ ...params, page: 1, order, desc });
  };

  const searchInputRef = useRef<HTMLInputElement>(null);

  const paginator = (
    <Paginator
      page={params.page ?? 1}
      limit={params.limit ?? 10}
      totalCount={result.totalCount}
      setPage={setPage}
      setLimit={setLimit}
    />
  );

  return (
    <main>
      <section className='bg-white'>
        <div className='layout relative flex min-h-screen flex-col py-4'>
          <header className='text-left border-gray-400 border-b-2 pb-4 mb-2'>
            <Logo className='w-8 float-left mr-2' />
            <h1 className='float-left text-2xl'>Code challenge</h1>
            {!isLoading ? (
              <span className='float-right text-gray-600'>
                <span className='inline-block mt-2 italic'>
                  {Intl.NumberFormat('en-us').format(loadTime)}ms
                </span>
                <Timer className='inline-block ml-1 -mt-1' />
              </span>
            ) : (
              <span></span>
            )}
          </header>

          <div className='h-10'>
            <input
              type='text'
              placeholder='Search'
              className='form-input py-1 mt-1 float-left'
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                // remove any irrelevant characters from the search query
                const value = e.target.value.replace(/[^a-z0-9.]/gi, '');
                if (value !== e.target.value) {
                  e.target.value = value;
                }
                setQuery(value);
              }}
              ref={searchInputRef}
            />
            {searchInputRef?.current?.value ? (
              <Tooltip
                className='inline-block'
                tooltipChildren={<span>Clear</span>}
                tooltipClassName='left-0 top-full mt-3'
              >
                <IconButton
                  icon={X}
                  className='ml-1 mt-1 p-0 float-left border-red-300 text-red-500'
                  variant='ghost'
                  onClick={() => {
                    if (searchInputRef?.current?.value) {
                      searchInputRef.current.value = '';
                    }
                    setQuery('');
                  }}
                />
              </Tooltip>
            ) : (
              <span></span>
            )}
            {paginator}
          </div>

          <div className='mt-2 w-full mx-auto'>
            <div className='w-full border-gray-400 border-t-2'>
              <div className='h-1 w-full bg-gray-200 overflow-hidden'>
                {isLoading ? (
                  <div className='animate-progress w-full h-full bg-blue-500 origin-left-right'></div>
                ) : (
                  <div></div>
                )}
              </div>
            </div>
            {result.assets.length > 0 ? (
              <div>
                <Table
                  assets={result.assets}
                  params={params}
                  setOrder={setOrder}
                />
                {paginator}
              </div>
            ) : !isLoading ? (
              <div className='p-4 italic'>No results</div>
            ) : (
              <div></div>
            )}
          </div>
        </div>
      </section>
    </main>
  );
}
