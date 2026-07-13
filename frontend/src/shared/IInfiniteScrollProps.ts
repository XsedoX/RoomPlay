export interface IInfiniteScrollProps {
  side: 'end' | 'start' | 'both';
  done: (status: 'error' | 'loading' | 'empty' | 'ok') => void;
}
