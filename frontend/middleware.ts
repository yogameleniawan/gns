import { NextRequest, NextResponse } from 'next/server';

export function middleware(request: NextRequest) {
  // Get locale from cookie or default to 'en'
  const locale = request.cookies.get('NEXT_LOCALE')?.value || 'en';
  
  // Add locale to response headers for server components
  const response = NextResponse.next();
  response.headers.set('x-locale', locale);
  
  return response;
}

export const config = {
  // Skip middleware for static files and API routes
  matcher: ['/((?!api|_next/static|_next/image|favicon.ico).*)'],
};
