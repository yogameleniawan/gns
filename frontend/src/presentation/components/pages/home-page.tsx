'use client'

import Link from 'next/link'
import {
    ArrowRight, Briefcase, LineChart, Calculator,
    DollarSign, TrendingUp, Activity, Receipt, Repeat, Shield, Target,
    PieChart, GitBranch, Scissors, ArrowUpCircle
} from 'lucide-react'
import { Card } from '@/src/presentation/components/ui/card'
import { Separator } from '@/src/presentation/components/ui/separator'

interface FeatureCard {
    title: string
    description: string
    href: string
    icon: React.ElementType
    badge?: string
}

export function HomePage() {
    const mainFeatures: FeatureCard[] = [
        {
            title: 'Market Overview',
            description: 'Real-time IHSG chart and market screener for all Indonesian stocks. Track top gainers, losers, and most active stocks with live data from TradingView.',
            href: '/market',
            icon: LineChart,
            badge: 'Live Data',
        },
        {
            title: 'Broker Flow',
            description: 'Monitor broker trading activities and institutional money flow. Analyze stock holdings and track big players movements in real-time.',
            href: '/broker-flow',
            icon: Briefcase,
            badge: 'Real-time',
        },
    ]

    const basicCalculators: FeatureCard[] = [
        {
            title: 'Profit Calculator',
            description: 'Calculate potential profits/losses from stock trading with buy and sell price analysis',
            href: '/calculator/profit',
            icon: DollarSign,
        },
        {
            title: 'Dividend Calculator',
            description: 'Estimate dividend returns based on stock price, yield, and investment amount',
            href: '/calculator/dividend',
            icon: TrendingUp,
        },
        {
            title: 'ARA/ARB Calculator',
            description: 'Calculate auto rejection ranges for your stock positions based on BEI rules',
            href: '/calculator/ara-arb',
            icon: Activity,
        },
        {
            title: 'Lot & Fee Calculator',
            description: 'Calculate complete transaction costs including broker fees and taxes',
            href: '/calculator/lot-fee',
            icon: Receipt,
        },
    ]

    const tradingCalculators: FeatureCard[] = [
        {
            title: 'Average Down/Up',
            description: 'Calculate new average price after buying additional stocks',
            href: '/calculator/average',
            icon: Repeat,
        },
        {
            title: 'Position Sizing',
            description: 'Determine optimal position size based on risk management rules',
            href: '/calculator/position-sizing',
            icon: Shield,
        },
        {
            title: 'Risk Reward',
            description: 'Calculate risk/reward ratio and evaluate trade quality',
            href: '/calculator/risk-reward',
            icon: Target,
        },
        {
            title: 'DCA Calculator',
            description: 'Simulate dollar cost averaging investment strategy',
            href: '/calculator/dca',
            icon: Repeat,
        },
        {
            title: 'Breakeven',
            description: 'Calculate required profit to recover from losses',
            href: '/calculator/breakeven',
            icon: ArrowUpCircle,
        },
        {
            title: 'Investment Growth',
            description: 'Compare growth between saving vs investing with compound returns',
            href: '/calculator/investment-growth',
            icon: TrendingUp,
        },
    ]

    const advancedCalculators: FeatureCard[] = [
        {
            title: 'Portfolio Diversification',
            description: 'Analyze portfolio allocation and get rebalancing recommendations',
            href: '/calculator/portfolio',
            icon: PieChart,
        },
        {
            title: 'Rights Issue (HMETD)',
            description: 'Calculate rights allocation and dilution effects',
            href: '/calculator/rights-issue',
            icon: GitBranch,
        },
        {
            title: 'Stock Split',
            description: 'Calculate effects of stock splits on your portfolio',
            href: '/calculator/stock-split',
            icon: Scissors,
        },
        {
            title: 'Moving Average Analyzer',
            description: 'Perform technical analysis using moving averages',
            href: '/calculator/moving-average',
            icon: TrendingUp,
        },
    ]

    return (
        <div className="space-y-12 pb-8">
            {/* Hero Section */}
            <div className="text-center space-y-4 pt-8">
                <h1 className="text-4xl md:text-6xl font-bold tracking-tight">
                    Welcome to <span className="text-foreground">Quantro</span>
                </h1>
                <p className="text-lg md:text-xl text-muted-foreground max-w-3xl mx-auto">
                    Your comprehensive platform for Indonesian stock market analysis, trading tools, and investment calculators.
                </p>
            </div>

            {/* Main Features */}
            <section className="space-y-6">
                <div>
                    <h2 className="text-2xl font-bold mb-2">Market Analysis</h2>
                    <p className="text-muted-foreground">Real-time market data and broker activity tracking</p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    {mainFeatures.map((feature) => {
                        const Icon = feature.icon
                        return (
                            <Link
                                key={feature.href}
                                href={feature.href}
                                className="group relative overflow-hidden rounded-xl border-2 bg-card p-6 transition-all duration-300 hover:border-foreground hover:shadow-xl hover:scale-[1.02]"
                            >
                                <div className="absolute inset-0 bg-foreground/0 transition-all duration-300 group-hover:bg-foreground/5" />

                                <div className="relative flex items-start gap-4">
                                    <div className="flex h-14 w-14 shrink-0 items-center justify-center rounded-xl border-2 border-foreground/20 bg-background transition-all duration-300 group-hover:border-foreground group-hover:bg-foreground group-hover:scale-110">
                                        <Icon className="h-7 w-7 text-foreground transition-colors duration-300 group-hover:text-background" />
                                    </div>

                                    <div className="flex-1 space-y-2">
                                        <div className="flex items-center gap-2">
                                            <h3 className="text-xl font-semibold">{feature.title}</h3>
                                            {feature.badge && (
                                                <span className="text-xs px-2 py-1 rounded-full border bg-background font-medium">
                                                    {feature.badge}
                                                </span>
                                            )}
                                            <ArrowRight className="h-5 w-5 ml-auto transition-transform duration-300 group-hover:translate-x-1" />
                                        </div>
                                        <p className="text-sm text-muted-foreground leading-relaxed">
                                            {feature.description}
                                        </p>
                                    </div>
                                </div>
                            </Link>
                        )
                    })}
                </div>
            </section>

            <Separator />

            {/* Calculators Section */}
            <section className="space-y-8">
                <div>
                    <div className="flex items-center gap-3 mb-2">
                        <Calculator className="h-6 w-6" />
                        <h2 className="text-2xl font-bold">Stock Calculators</h2>
                    </div>
                    <p className="text-muted-foreground">14 powerful calculators for stock trading and investment analysis</p>
                </div>

                {/* Basic Calculators */}
                <div className="space-y-4">
                    <h3 className="text-lg font-semibold text-muted-foreground">Basic Calculators</h3>
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                        {basicCalculators.map((calc) => {
                            const Icon = calc.icon
                            return (
                                <Link
                                    key={calc.href}
                                    href={calc.href}
                                    className="group relative overflow-hidden rounded-lg border bg-card p-4 transition-all duration-200 hover:border-foreground hover:shadow-md"
                                >
                                    <div className="absolute inset-0 bg-foreground/0 transition-all duration-200 group-hover:bg-foreground/5" />

                                    <div className="relative space-y-3">
                                        <div className="flex items-center justify-between">
                                            <div className="flex h-10 w-10 items-center justify-center rounded-lg border border-foreground/20 bg-background transition-all duration-200 group-hover:border-foreground group-hover:bg-foreground">
                                                <Icon className="h-5 w-5 text-foreground transition-colors duration-200 group-hover:text-background" />
                                            </div>
                                            <ArrowRight className="h-4 w-4 transition-transform duration-200 group-hover:translate-x-1" />
                                        </div>
                                        <div>
                                            <h4 className="font-semibold text-sm mb-1">{calc.title}</h4>
                                            <p className="text-xs text-muted-foreground leading-relaxed">
                                                {calc.description}
                                            </p>
                                        </div>
                                    </div>
                                </Link>
                            )
                        })}
                    </div>
                </div>

                {/* Trading Calculators */}
                <div className="space-y-4">
                    <h3 className="text-lg font-semibold text-muted-foreground">Trading & Risk Management</h3>
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                        {tradingCalculators.map((calc) => {
                            const Icon = calc.icon
                            return (
                                <Link
                                    key={calc.href}
                                    href={calc.href}
                                    className="group relative overflow-hidden rounded-lg border bg-card p-4 transition-all duration-200 hover:border-foreground hover:shadow-md"
                                >
                                    <div className="absolute inset-0 bg-foreground/0 transition-all duration-200 group-hover:bg-foreground/5" />

                                    <div className="relative space-y-3">
                                        <div className="flex items-center justify-between">
                                            <div className="flex h-10 w-10 items-center justify-center rounded-lg border border-foreground/20 bg-background transition-all duration-200 group-hover:border-foreground group-hover:bg-foreground">
                                                <Icon className="h-5 w-5 text-foreground transition-colors duration-200 group-hover:text-background" />
                                            </div>
                                            <ArrowRight className="h-4 w-4 transition-transform duration-200 group-hover:translate-x-1" />
                                        </div>
                                        <div>
                                            <h4 className="font-semibold text-sm mb-1">{calc.title}</h4>
                                            <p className="text-xs text-muted-foreground leading-relaxed">
                                                {calc.description}
                                            </p>
                                        </div>
                                    </div>
                                </Link>
                            )
                        })}
                    </div>
                </div>

                {/* Advanced Calculators */}
                <div className="space-y-4">
                    <h3 className="text-lg font-semibold text-muted-foreground">Corporate Actions & Analysis</h3>
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                        {advancedCalculators.map((calc) => {
                            const Icon = calc.icon
                            return (
                                <Link
                                    key={calc.href}
                                    href={calc.href}
                                    className="group relative overflow-hidden rounded-lg border bg-card p-4 transition-all duration-200 hover:border-foreground hover:shadow-md"
                                >
                                    <div className="absolute inset-0 bg-foreground/0 transition-all duration-200 group-hover:bg-foreground/5" />

                                    <div className="relative space-y-3">
                                        <div className="flex items-center justify-between">
                                            <div className="flex h-10 w-10 items-center justify-center rounded-lg border border-foreground/20 bg-background transition-all duration-200 group-hover:border-foreground group-hover:bg-foreground">
                                                <Icon className="h-5 w-5 text-foreground transition-colors duration-200 group-hover:text-background" />
                                            </div>
                                            <ArrowRight className="h-4 w-4 transition-transform duration-200 group-hover:translate-x-1" />
                                        </div>
                                        <div>
                                            <h4 className="font-semibold text-sm mb-1">{calc.title}</h4>
                                            <p className="text-xs text-muted-foreground leading-relaxed">
                                                {calc.description}
                                            </p>
                                        </div>
                                    </div>
                                </Link>
                            )
                        })}
                    </div>
                </div>
            </section>
        </div>
    )
}
