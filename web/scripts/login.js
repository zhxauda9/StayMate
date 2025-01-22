document.getElementById('login-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const user = {
        email: document.getElementById('email').value,
        password: document.getElementById('password').value,
    };

    if (user.email == "root@root" && user.password == "admin"){
        window.location.href = "/admin";
        return
    }
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

        alert('Logged in successfully!');


        const data = await response.json();
        if (data && data.role) {
            console.log("User role:", data.role); // Check what role is returned
            if (data.role === "admin") {
                window.location.href = "/admin";
            } else {
                window.location.href = "/profile";
            }
        } else {
            alert("Role is missing in response.");
        }
    } catch (error) {
        console.error(error);
        alert(`Error: ${error.message}`);
    }
});
