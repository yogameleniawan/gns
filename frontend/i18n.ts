import { getRequestConfig} from 'next-intl/server';
import { cookies } from 'next/headers';

export const locales = ['en', 'id'] as const;
export type Locale = (typeof locales)[number];

export const defaultLocale: Locale = 'en';

export default getRequestConfig(async () => {
  // Get locale from cookie or use default
  const cookieStore = await cookies();
  const locale = (cookieStore.get('NEXT_LOCALE')?.value as Locale) || defaultLocale;

  return {
    locale,
    messages: (await import(`../locales/${locale}.json`)).default
  };
});
