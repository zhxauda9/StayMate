document.getElementById('login-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const user = {
        email: document.getElementById('email').value,
        password: document.getElementById('password').value,
    };

    try {
        const response = await fetch('/api/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(user),
        });

        if (!response.ok) {
            const errorResponse = await response.json();
            throw new Error(errorResponse.message || 'Failed to log in.');
        }

        console.log('Logged in successfully!');
        window.location.href = "/profile";
    } catch (error) {
        console.error(error);
        alert(`Error: ${error.message}`);
    }
});
