'use client';

import { useCallback, useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { getUser, isAuthenticated as checkIsAuthenticated, logout as clearAuth } from '@/lib/auth';
import { User } from '@/types';

type AuthStatus = 'checking' | 'authenticated' | 'unauthenticated';

interface AuthGuardOptions {
  redirectTo?: string;
  requireAdmin?: boolean;
}

export function useAuthGuard(options: AuthGuardOptions = {}) {
  const { redirectTo = '/login', requireAdmin = false } = options;
  const router = useRouter();
  const [user, setUser] = useState<User | null>(null);
  const [status, setStatus] = useState<AuthStatus>('checking');

  useEffect(() => {
    const authed = checkIsAuthenticated();
    const storedUser = getUser();
    const authorized = authed && (!requireAdmin || storedUser?.role === 'admin');

    if (!authorized) {
      setStatus('unauthenticated');
      router.push(redirectTo);
      return;
    }

    setUser(storedUser);
    setStatus('authenticated');
  }, [router, redirectTo, requireAdmin]);

  const handleLogout = useCallback(() => {
    clearAuth();
    router.push('/login');
  }, [router]);

  return {
    user,
    status,
    isLoading: status === 'checking',
    isAuthenticated: status === 'authenticated',
    isAdmin: user?.role === 'admin',
    logout: handleLogout,
  };
}

export function useAuthRedirect(authenticatedPath: string, unauthenticatedPath: string) {
  const router = useRouter();

  useEffect(() => {
    if (checkIsAuthenticated()) {
      router.push(authenticatedPath);
    } else {
      router.push(unauthenticatedPath);
    }
  }, [router, authenticatedPath, unauthenticatedPath]);
}
