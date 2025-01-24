document.getElementById('register-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const name = document.getElementById('name').value.trim();
    const email = document.getElementById('email').value.trim();
    const password = document.getElementById('password').value.trim();

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
        alert('Please enter a valid email address.');
        return;
    }

    if (password.length < 8) {
        alert('Password must be at least 8 characters long.');
        return;
    }

    if (!/[A-Z]/.test(password)) {
        alert('Password must contain at least one uppercase letter.');
        return;
    }

    if (!/[a-z]/.test(password)) {
        alert('Password must contain at least one lowercase letter.');
        return;
    }

    if (!/[0-9]/.test(password)) {
        alert('Password must contain at least one number.');
        return;
    }

    if (!/[\W_]/.test(password)) {
        alert('Password must contain at least one special character (e.g., @, #, $, etc.).');
        return;
    }

    const user = { name, email, password };

    try {
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(user),
        });
        if (!response.ok) {
            const contentType = response.headers.get('Content-Type');
            let errorMessage;
            if (contentType && contentType.includes('application/json')) {
                const errorResponse = await response.json();
                errorMessage = errorResponse.message || 'Failed to create user.';
            } else {
                errorMessage = await response.text();
            }
            throw new Error(errorMessage);
        }
        window.location.href = "/verify";
    } catch (error) {
        console.error(error);
        alert(`Error: ${error.message}`);
    }
});
