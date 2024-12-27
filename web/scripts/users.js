const BASE_URL = 'http://localhost:8080'; // Backend server URL

async function loadUsers() {
    try {
        const response = await fetch(`${BASE_URL}/users`);
        const users = await response.json();
        const table = document.getElementById('users-table');
        table.innerHTML = '';
        users.forEach(user => {
            const row = document.createElement('tr');
            row.innerHTML = `
        <td>${user.id}</td>
        <td>${user.name}</td>
        <td>${user.email}</td>
        <td>
            <button class="btn btn-warning btn-sm" onclick="updateUser(${user.id})">Update</button>
            <button class="btn btn-danger btn-sm" onclick="deleteUser(${user.id})">Delete</button>
        </td>
    `;
            table.appendChild(row);
        });
    } catch (error) {
        console.error(error);
        alert("Ошибка загрузки бронирований.");
    }
}
loadUsers();


document.getElementById('create-user-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const user = {
        name: document.getElementById('user-name').value,
        email: document.getElementById('user-email').value,
    };

    try {
        const response = await fetch(`${BASE_URL}/users`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(user),
        });

        if (!response.ok) {
            throw new Error("Failed to create user.");
        }

        alert("User created successfully!");
        loadUsers();
    } catch (error) {
        console.error(error);
        alert("Error: Failed to create user.");
    }
});


async function deleteUser(id) {
    if (confirm("Are you sure you want to delete the user?")) {
        await fetch(`${BASE_URL}/users/${id}`, { method: 'DELETE' });
        alert("User deleted.")
        loadUsers();
    } else {
        alert("Deletion cancelled.");
    }

}

async function updateUser(id) {
    const name = prompt('Enter new name:');
    const email = prompt('Enter new email:');
    if (!name || !email) return;
    await fetch(`${BASE_URL}/users/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, email }),
    });
    alert("User updated!");
    loadUsers();
}
document.getElementById('search-user-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const userId = document.getElementById('search-user-id').value;
    const resultDiv = document.getElementById('search-result');
    resultDiv.innerHTML = '';

    try {
        const response = await fetch(`${BASE_URL}/users/${userId}`);
        if (!response.ok) throw new Error('User not found.');
        const user = await response.json();

        resultDiv.innerHTML = `
        <div class="alert alert-info">
            <p><strong>ID:</strong> ${user.id}</p>
            <p><strong>Name:</strong> ${user.name}</p>
            <p><strong>Email:</strong> ${user.email}</p>
        </div>
    `;
    } catch (error) {
        console.error(error);
        resultDiv.innerHTML = `
        <div class="alert alert-danger">Error: User not found.</div>
    `;
    }
});