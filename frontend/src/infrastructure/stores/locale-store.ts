import { create } from 'zustand';
import { persist } from 'zustand/middleware';

export type Locale = 'en' | 'id';

interface LocaleState {
  locale: Locale;
  hasHydrated: boolean;
  setLocale: (locale: Locale) => void;
  setHasHydrated: (state: boolean) => void;
}

export const useLocaleStore = create<LocaleState>()(
  persist(
    (set) => ({
      locale: 'en',
      hasHydrated: false,

      setLocale: (locale: Locale) => {
        set({ locale });
        // Set cookie for server-side rendering
        if (typeof document !== 'undefined') {
          document.cookie = `NEXT_LOCALE=${locale};path=/;max-age=31536000`;
        }
      },

      setHasHydrated: (state: boolean) => {
        set({ hasHydrated: state });
      },
    }),
    {
      name: 'locale-storage',
      skipHydration: true,
      partialize: (state) => ({
        locale: state.locale,
      }),
      onRehydrateStorage: () => (state) => {
        state?.setHasHydrated(true);
      },
    }
  )
);
