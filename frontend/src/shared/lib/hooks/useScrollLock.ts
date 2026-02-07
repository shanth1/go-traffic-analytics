import { useLayoutEffect } from 'react';

export const useScrollLock = (lock: boolean) => {
  useLayoutEffect(() => {
    if (!lock) return;

    // Lock
    document.body.classList.add('scroll-locked');

    return () => {
      // Unlock
      document.body.classList.remove('scroll-locked');
    };
  }, [lock]);
};
