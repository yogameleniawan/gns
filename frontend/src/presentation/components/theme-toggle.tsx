'use client'

import * as React from 'react'
import { Moon, Sun } from 'lucide-react'
import { useTheme } from 'next-themes'
import { Button } from '@/src/presentation/components/ui/button'

export function ThemeToggle() {
    const { setTheme, theme } = useTheme()
    const [mounted, setMounted] = React.useState(false)

    // Prevent hydration mismatch by only rendering after mount
    React.useEffect(() => {
        setMounted(true)
    }, [])

    if (!mounted) {
        return (
            <Button variant="ghost" size="icon" aria-label="Toggle theme">
                <Sun className="h-[1.2rem] w-[1.2rem]" />
                <span className="sr-only">Toggle theme</span>
            </Button>
        )
    }

    return (
        <Button
            variant="ghost"
            size="icon"
            onClick={() => setTheme(theme === 'light' ? 'dark' : 'light')}
            aria-label="Toggle theme"
        >
            <Sun className={`h-[1.2rem] w-[1.2rem] ${theme === 'dark' ? 'hidden' : 'block'}`} />
            <Moon className={`h-[1.2rem] w-[1.2rem] ${theme === 'dark' ? 'block' : 'hidden'}`} />
            <span className="sr-only">Toggle theme</span>
        </Button>
    )
}
