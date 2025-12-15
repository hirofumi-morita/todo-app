'use client';

import { useAuthRedirect } from '@/hooks/useAuth';

export default function Home() {
  useAuthRedirect('/dashboard', '/login');

  return (
    <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
      <p>Loading...</p>
    </div>
  );
}
