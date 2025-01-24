document.getElementById('verify-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const verificationCode = document.getElementById('verificationCode').value.trim();
    try {
        const response = await fetch('/api/verify', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ verificationCode }),
        });

        if (!response.ok) {
            throw new Error('Verification failed. Please try again.');
        }

        alert('Verification successful!');
        window.location.href = '/login';
    } catch (error) {
        alert(error.message);
    }
});
