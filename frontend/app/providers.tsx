'use client'

import { QueryClientProvider } from '@tanstack/react-query'
import { queryClient } from '@/src/lib/query-client'
import { ThemeProvider } from '@/src/presentation/components/theme-provider'
import { NextIntlClientProvider } from 'next-intl'
import { useLocaleStore } from '@/src/infrastructure/stores/locale-store'
import { ReactNode, useEffect, useState } from 'react'

export function Providers({ children }: { children: ReactNode }) {
    const { locale, hasHydrated } = useLocaleStore();
    const [messages, setMessages] = useState<any>(null);

    useEffect(() => {
        if (!hasHydrated) {
            useLocaleStore.persist.rehydrate();
        }
    }, [hasHydrated]);

    useEffect(() => {
        // Load messages for current locale
        import(`../locales/${hasHydrated ? locale : 'en'}.json`)
            .then((msgs) => setMessages(msgs.default))
            .catch(() => setMessages({}));
    }, [locale, hasHydrated]);

    // Detect theme for loading screen (before ThemeProvider renders)
    const [isDark, setIsDark] = useState(true); // Default to dark

    useEffect(() => {
        // Check localStorage for saved theme
        const savedTheme = localStorage.getItem('theme');
        if (savedTheme) {
            setIsDark(savedTheme === 'dark');
        } else {
            // Fallback to system preference
            setIsDark(window.matchMedia('(prefers-color-scheme: dark)').matches);
        }
    }, []);

    if (!messages) {
        return (
            <div className={`fixed inset-0 z-50 flex items-center justify-center ${isDark ? 'bg-[#0a0a0a]' : 'bg-white'}`}>
                {/* Gradient Background */}
                <div className="absolute inset-0 overflow-hidden pointer-events-none">
                    <div className={`absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 ${isDark ? 'bg-emerald-500/20' : 'bg-emerald-500/30'} rounded-full blur-3xl animate-pulse`} />
                    <div className={`absolute top-1/3 left-1/3 w-64 h-64 ${isDark ? 'bg-teal-500/15' : 'bg-teal-500/25'} rounded-full blur-3xl animate-pulse`} style={{ animationDelay: '0.5s' }} />
                    <div className={`absolute bottom-1/3 right-1/3 w-72 h-72 ${isDark ? 'bg-cyan-500/15' : 'bg-cyan-500/25'} rounded-full blur-3xl animate-pulse`} style={{ animationDelay: '1s' }} />
                </div>

                {/* Loading Content */}
                <div className="relative flex flex-col items-center gap-6">
                    {/* Animated Logo */}
                    <div className="relative">
                        <div className="absolute inset-0 animate-ping opacity-20">
                            <div className="h-16 w-16 text-emerald-500">
                                <svg fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
                                    <path clipRule="evenodd" d="M24 4L44 24L24 44L4 24L24 4ZM24 14V34L34 24L24 14Z" fill="currentColor" fillRule="evenodd" />
                                </svg>
                            </div>
                        </div>
                        <div className="h-16 w-16 text-emerald-500 animate-pulse">
                            <svg fill="none" viewBox="0 0 48 48" xmlns="http://www.w3.org/2000/svg">
                                <path clipRule="evenodd" d="M24 4L44 24L24 44L4 24L24 4ZM24 14V34L34 24L24 14Z" fill="currentColor" fillRule="evenodd" />
                            </svg>
                        </div>
                    </div>

                    {/* Brand Name */}
                    <div className="flex flex-col items-center gap-2">
                        <span className={`text-2xl font-bold tracking-tight ${isDark ? 'text-white' : 'text-gray-900'}`}>Quantro</span>
                        <span className={`text-sm ${isDark ? 'text-gray-400' : 'text-gray-500'}`}>Loading your experience...</span>
                    </div>

                    {/* Loading Bar */}
                    <div className={`w-48 h-1 ${isDark ? 'bg-gray-800' : 'bg-gray-200'} rounded-full overflow-hidden`}>
                        <div className="h-full bg-gradient-to-r from-emerald-500 to-teal-500 rounded-full animate-loading-bar" />
                    </div>
                </div>

                {/* Animation Styles */}
                <style jsx>{`
                    @keyframes loading-bar {
                        0% { width: 0%; transform: translateX(0); }
                        50% { width: 70%; }
                        100% { width: 100%; transform: translateX(0); }
                    }
                    .animate-loading-bar {
                        animation: loading-bar 1.5s ease-in-out infinite;
                    }
                `}</style>
            </div>
        );
    }

    return (
        <NextIntlClientProvider locale={hasHydrated ? locale : 'en'} messages={messages}>
            <ThemeProvider
                attribute="class"
                defaultTheme="dark"
                enableSystem={false}
                disableTransitionOnChange
            >
                <QueryClientProvider client={queryClient}>
                    {children}
                </QueryClientProvider>
            </ThemeProvider>
        </NextIntlClientProvider>
    )
}
