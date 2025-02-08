function logout() {
    document.cookie = 'Authorization=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/; Secure; SameSite=Lax';

    localStorage.clear();

    if ('caches' in window) {
        caches.keys().then(cacheNames => {
            cacheNames.forEach(cacheName => caches.delete(cacheName));
        });
    }

    fetch('/auth/logout', { method: 'POST', credentials: 'include' })
        .then(() => {
            window.location.href = '/login';
        })
        .catch(err => {
            console.error('Logout failed:', err);
            alert('Failed to logout');
        });
}