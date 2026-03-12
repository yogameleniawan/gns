'use client';

import { useLocaleStore, type Locale } from '@/src/infrastructure/stores/locale-store';
import { useTranslations as useNextIntlTranslations } from 'next-intl';
import { useEffect } from 'react';

export function useLocale() {
  const { locale, setLocale, hasHydrated } = useLocaleStore();
  
  // Ensure hydration is complete
  useEffect(() => {
    if (!hasHydrated) {
      useLocaleStore.persist.rehydrate();
    }
  }, [hasHydrated]);

  const changeLocale = (newLocale: Locale) => {
    setLocale(newLocale);
    // Reload to apply new locale
    window.location.reload();
  };

  return {
    locale: hasHydrated ? locale : 'en',
    changeLocale,
    isReady: hasHydrated,
  };
}

// Re-export useTranslations for convenience
export function useTranslations(namespace?: string) {
  return useNextIntlTranslations(namespace);
}
