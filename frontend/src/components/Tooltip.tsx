export const Tooltip = ({
  children,
  className,
  tooltipChildren,
  tooltipClassName,
}: {
  children: React.ReactNode;
  tooltipChildren: React.ReactNode;
  className?: string;
  tooltipClassName?: string;
}) => (
  <div
    className={['group', 'relative', ...(className ? [className] : [])].join(
      ' ',
    )}
  >
    {children}
    <span
      className={[
        'group-hover:opacity-100',
        'pointer-events-none',
        'transition-opacity',
        'bg-gray-500',
        'text-sm',
        'text-gray-100',
        'rounded-md',
        'absolute',
        'z-10',
        'opacity-0',
        'px-2',
        'py-1',
        ...(tooltipClassName ? [tooltipClassName] : []),
      ].join(' ')}
    >
      {tooltipChildren}
    </span>
  </div>
);
