'use client';

import * as React from 'react';
import { useState, useEffect } from 'react';

import Logo from '~/svg/Logo.svg';
import ButtonLink from '@/components/links/ButtonLink';

export default function HomePage() {
  const [data, setData] = useState<any[]>([]);
  const [sortedData, setSortedData] = useState<any[]>([]);
  const [isSorting, setIsSorting] = useState(false);

  useEffect(() => {
    async function fetchData() {
      try {
        const response = await fetch('http://localhost:8080/assets');
        const result = await response.json();

        setData(result);
        sortData(result); // Automatically sort after fetching
      } catch (err) {
        console.error('Error fetching data:', err);
      }
    }

    fetchData();
  }, []);

  const sortData = (fetchedData: any[]) => {
    setIsSorting(true);
    setTimeout(() => {
      const newSortedData = [...fetchedData].sort((a, b) => a.Host.localeCompare(b.Host));
      setSortedData(newSortedData);
      setIsSorting(false);
    }, 1000); // Artificial delay to slow down sorting
  };

  const renderedData = sortedData.map((item, index) => {
    return (
        <div key={index} className="p-4">
          <p>{`ID: ${item.ID}`}</p>
          <p>{`Host: ${item.Host}`}</p>
          <p>{`Comment: ${item.Comment}`}</p>
          <p>{`Owner: ${item.Owner}`}</p>
          <p>{`IPs: ${(item.IPs || []).map((ip: any) => ip.Address).join(', ')}`}</p>
          <p>{`Ports: ${(item.Ports || []).map((port: any) => port.Port).join(', ')}`}</p>
        </div>
    );
  });

  return (
      <main>
        <section className='bg-white'>
          <div className='layout relative flex min-h-screen flex-col items-center justify-center py-12 text-center'>
            <Logo className='w-16' />
            <h1 className='mt-4'>Code challenge</h1>

            <p className='mt-2 text-sm text-gray-800'>
              You have complete freedom to present the data here.
            </p>

            <ButtonLink className='mt-6' href='/components' variant='light'>
              See all included components
            </ButtonLink>

            <div className="mt-8 w-full max-w-2xl mx-auto bg-gray-100 p-4">
              {sortedData.length === 0 ? (
                  <p>{isSorting ? 'Sorting...' : 'Loading...'}</p>
              ) : (
                  <div>{renderedData}</div>
              )}
            </div>
          </div>
        </section>
      </main>
  );
}
