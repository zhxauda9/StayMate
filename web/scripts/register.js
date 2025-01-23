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

    // Password validation
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

    const user = { name, email, password};

    try {
        const response = await fetch('/api/register', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(user),
        });

        if (!response.ok) {
            const errorResponse = await response.json();
            throw new Error(errorResponse.message || 'Failed to create user.');
        }

        window.location.href = "/login";
    } catch (error) {
        console.error(error);
        alert(`Error: ${error.message}`);
    }
});

