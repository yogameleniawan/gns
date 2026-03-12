'use client'

import { ReactNode, useState } from 'react'
import { Navbar } from './navbar'
import { Sidebar } from './sidebar'

interface MainLayoutProps {
    children: ReactNode
}

export function MainLayout({ children }: MainLayoutProps) {
    const [isSidebarOpen, setIsSidebarOpen] = useState(false)

    return (
        <div className="min-h-screen bg-background">
            <Sidebar
                isOpen={isSidebarOpen}
                onClose={() => setIsSidebarOpen(false)}
            />

            <div className="lg:pl-[280px]">
                <Navbar onMenuToggle={() => setIsSidebarOpen(!isSidebarOpen)} />
                <main className="container mx-auto px-4 py-6 sm:px-6 lg:px-8">
                    {children}
                </main>
            </div>
        </div>
    )
}
