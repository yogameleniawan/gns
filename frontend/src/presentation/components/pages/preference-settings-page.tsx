'use client'

import { useState, useEffect } from 'react'
import { Globe, Moon, Sun, Check } from 'lucide-react'
import { useTheme } from 'next-themes'
import { useLocale, useTranslations } from '@/src/presentation/hooks/use-locale'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/src/presentation/components/ui/card'
import { Label } from '@/src/presentation/components/ui/label'
import { RadioGroup, RadioGroupItem } from '@/src/presentation/components/ui/radio-group'
import { Button } from '@/src/presentation/components/ui/button'
import { Separator } from '@/src/presentation/components/ui/separator'
import { toast } from 'sonner'
import { MainLayout } from '../layout/main-layout'
import { ProtectedRoute } from '../layout/protected-route'

export function PreferenceSettingsPage() {
    const t = useTranslations('settings.preferences')
    const tCommon = useTranslations('common')
    const { theme, setTheme } = useTheme()
    const { locale, changeLocale, isReady } = useLocale()

    const [selectedLocale, setSelectedLocale] = useState<'en' | 'id'>(locale)
    const [selectedTheme, setSelectedTheme] = useState(theme || 'dark')
    const [mounted, setMounted] = useState(false)

    useEffect(() => {
        setMounted(true)
    }, [])

    useEffect(() => {
        if (isReady) {
            setSelectedLocale(locale)
        }
    }, [locale, isReady])

    useEffect(() => {
        if (theme) {
            setSelectedTheme(theme)
        }
    }, [theme])

    const handleSave = () => {
        // Apply theme
        if (selectedTheme !== theme) {
            setTheme(selectedTheme)
        }

        // Apply locale (this will reload the page)
        if (selectedLocale !== locale) {
            changeLocale(selectedLocale)
        } else {
            // If no locale change, just show success toast
            toast.success(t('changesSaved'))
        }
    }

    if (!mounted || !isReady) {
        return (
            <div className="flex h-screen items-center justify-center">
                <div className="text-muted-foreground">{tCommon('loading')}</div>
            </div>
        )
    }

    return (
        <ProtectedRoute>
            <MainLayout>
                <div className="space-y-6 pb-8">
                    <div>
                        <h1 className="text-3xl font-bold tracking-tight">{t('title')}</h1>
                        <p className="text-muted-foreground mt-2">
                            Manage your application preferences and settings
                        </p>
                    </div>

                    <Separator />

                    {/* Language Settings */}
                    <Card>
                        <CardHeader>
                            <div className="flex items-center gap-3">
                                <div className="h-10 w-10 rounded-lg bg-gradient-to-br from-emerald-500 to-teal-600 flex items-center justify-center">
                                    <Globe className="h-5 w-5 text-white" />
                                </div>
                                <div>
                                    <CardTitle>{t('language.title')}</CardTitle>
                                    <CardDescription>{t('language.description')}</CardDescription>
                                </div>
                            </div>
                        </CardHeader>
                        <CardContent>
                            <RadioGroup value={selectedLocale} onValueChange={(value) => setSelectedLocale(value as 'en' | 'id')}>
                                <div className="flex flex-col gap-3">
                                    <div className={`flex items-center space-x-3 rounded-lg border-2 p-4 cursor-pointer transition-all ${selectedLocale === 'en'
                                        ? 'border-emerald-500 bg-emerald-500/5'
                                        : 'border-border hover:border-foreground/20 hover:bg-muted/50'
                                        }`} onClick={() => setSelectedLocale('en')}>
                                        <RadioGroupItem value="en" id="en" />
                                        <Label htmlFor="en" className="flex-1 cursor-pointer font-medium">
                                            {t('language.english')}
                                        </Label>
                                        {selectedLocale === 'en' && (
                                            <Check className="h-5 w-5 text-emerald-500" />
                                        )}
                                    </div>

                                    <div className={`flex items-center space-x-3 rounded-lg border-2 p-4 cursor-pointer transition-all ${selectedLocale === 'id'
                                        ? 'border-emerald-500 bg-emerald-500/5'
                                        : 'border-border hover:border-foreground/20 hover:bg-muted/50'
                                        }`} onClick={() => setSelectedLocale('id')}>
                                        <RadioGroupItem value="id" id="id" />
                                        <Label htmlFor="id" className="flex-1 cursor-pointer font-medium">
                                            {t('language.indonesian')}
                                        </Label>
                                        {selectedLocale === 'id' && (
                                            <Check className="h-5 w-5 text-emerald-500" />
                                        )}
                                    </div>
                                </div>
                            </RadioGroup>
                        </CardContent>
                    </Card>

                    {/* Theme Settings */}
                    <Card>
                        <CardHeader>
                            <div className="flex items-center gap-3">
                                <div className="h-10 w-10 rounded-lg bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center">
                                    {selectedTheme === 'dark' ? (
                                        <Moon className="h-5 w-5 text-white" />
                                    ) : (
                                        <Sun className="h-5 w-5 text-white" />
                                    )}
                                </div>
                                <div>
                                    <CardTitle>{t('theme.title')}</CardTitle>
                                    <CardDescription>{t('theme.description')}</CardDescription>
                                </div>
                            </div>
                        </CardHeader>
                        <CardContent>
                            <RadioGroup value={selectedTheme} onValueChange={setSelectedTheme}>
                                <div className="flex flex-col gap-3">
                                    <div className={`flex items-center space-x-3 rounded-lg border-2 p-4 cursor-pointer transition-all ${selectedTheme === 'light'
                                        ? 'border-emerald-500 bg-emerald-500/5'
                                        : 'border-border hover:border-foreground/20 hover:bg-muted/50'
                                        }`} onClick={() => setSelectedTheme('light')}>
                                        <RadioGroupItem value="light" id="light" />
                                        <Label htmlFor="light" className="flex-1 cursor-pointer font-medium flex items-center gap-2">
                                            <Sun className="h-4 w-4" />
                                            {t('theme.light')}
                                        </Label>
                                        {selectedTheme === 'light' && (
                                            <Check className="h-5 w-5 text-emerald-500" />
                                        )}
                                    </div>

                                    <div className={`flex items-center space-x-3 rounded-lg border-2 p-4 cursor-pointer transition-all ${selectedTheme === 'dark'
                                        ? 'border-emerald-500 bg-emerald-500/5'
                                        : 'border-border hover:border-foreground/20 hover:bg-muted/50'
                                        }`} onClick={() => setSelectedTheme('dark')}>
                                        <RadioGroupItem value="dark" id="dark" />
                                        <Label htmlFor="dark" className="flex-1 cursor-pointer font-medium flex items-center gap-2">
                                            <Moon className="h-4 w-4" />
                                            {t('theme.dark')}
                                        </Label>
                                        {selectedTheme === 'dark' && (
                                            <Check className="h-5 w-5 text-emerald-500" />
                                        )}
                                    </div>
                                </div>
                            </RadioGroup>
                        </CardContent>
                    </Card>

                    {/* Save Button */}
                    <div className="flex justify-end gap-3">
                        <Button
                            variant="outline"
                            onClick={() => {
                                setSelectedLocale(locale)
                                setSelectedTheme(theme || 'dark')
                            }}
                        >
                            {tCommon('cancel')}
                        </Button>
                        <Button
                            onClick={handleSave}
                        >
                            {t('saveChanges')}
                        </Button>
                    </div>
                </div>
            </MainLayout>
        </ProtectedRoute>
    )
}
