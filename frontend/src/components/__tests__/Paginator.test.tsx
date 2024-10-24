import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import { Paginator } from '@/components/Paginator';

describe('Paginator component', () => {
  const setup = (page = 1, limit = 10, totalCount = 100) => {
    const setPage = jest.fn();
    const setLimit = jest.fn();
    const utils = render(
      <Paginator
        page={page}
        limit={limit}
        totalCount={totalCount}
        setPage={setPage}
        setLimit={setLimit}
      />,
    );
    return {
      ...utils,
      setPage,
      setLimit,
    };
  };

  it('calls setPage with 1 when the first page button is clicked', () => {
    const { getByTestId, setPage } = setup(2);
    fireEvent.click(getByTestId('paginator-first'));
    expect(setPage).toHaveBeenCalledWith(1);
  });

  it('calls setPage with the previous page number when the previous page button is clicked', () => {
    const { getByTestId, setPage } = setup(2);
    fireEvent.click(getByTestId('paginator-previous'));
    expect(setPage).toHaveBeenCalledWith(1);
  });

  it('calls setPage with the next page number when the next page button is clicked', () => {
    const { getByTestId, setPage } = setup(1);
    fireEvent.click(getByTestId('paginator-next'));
    expect(setPage).toHaveBeenCalledWith(2);
  });

  it('calls setPage with the last page number when the last page button is clicked', () => {
    const { getByTestId, setPage } = setup(1, 10, 100);
    fireEvent.click(getByTestId('paginator-last'));
    expect(setPage).toHaveBeenCalledWith(10);
  });

  it('calls setLimit with the selected value when the limit select element is changed', () => {
    const { getByDisplayValue, setLimit } = setup();
    fireEvent.change(getByDisplayValue('10'), { target: { value: '20' } });
    expect(setLimit).toHaveBeenCalledWith(20);
  });
});
