import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatLots(lots: number): string {
  const absLots = Math.abs(lots)
  if (absLots >= 1000000) {
    return `${(absLots / 1000000).toFixed(1)}M`
  }
  if (absLots >= 1000) {
    return `${(absLots / 1000).toFixed(1)}K`
  }
  return absLots.toFixed(0)
}

export function formatPrice(price: number): string {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(price)
}

export function formatValue(value: number): string {
  const absValue = Math.abs(value)
  if (absValue >= 1000000000) {
    return `${(absValue / 1000000000).toFixed(1)}B`
  }
  if (absValue >= 1000000) {
    return `${(absValue / 1000000).toFixed(1)}M`
  }
  if (absValue >= 1000) {
    return `${(absValue / 1000).toFixed(1)}K`
  }
  return new Intl.NumberFormat('id-ID').format(absValue)
}

import { TimeRange } from '@/src/application/store/broker.store'
import { startOfMonth, endOfMonth, subDays } from 'date-fns'

export function getDateRange(timeRange: TimeRange, customDateRange: { from: Date | null; to: Date | null }) {
  const today = new Date()
  let from = today
  let to = today

  switch (timeRange) {
    case 'today':
      from = today
      to = today
      break
    case '1w':
      from = subDays(today, 7)
      to = today
      break
    case '1m':
      from = startOfMonth(today)
      to = endOfMonth(today)
      break
    case 'range':
      if (customDateRange.from && customDateRange.to) {
        from = customDateRange.from
        to = customDateRange.to
      }
      break
  }

  return { from, to }
}
