async function loadUsers() {
    try {
        const response = await fetch(`/users`);
        const users = await response.json();
        const table = document.getElementById('users-table');
        table.innerHTML = '';
        users.forEach(user => {
            const row = document.createElement('tr');
            row.innerHTML = `
        <td>${user.id}</td>
        <td>${user.name}</td>
        <td>${user.email}</td>
        <td>${user.role}</td>
        <td>${user.status ? user.status : 'No status'}</td>
        <td><img src="${user.photo}" alt="Room Photo" style="width: 100px; height: auto;"></td>
        <td>
            <button class="btn btn-warning btn-sm" onclick="updateUser(${user.id})">Update</button>
            <button class="btn btn-danger btn-sm" onclick="deleteUser(${user.id})">Delete</button>
            <button class="btn btn-info btn-sm" onclick="sendEmail(${user.id})">Send Email</button>
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
    const name= document.getElementById('user-name').value;
    const email= document.getElementById('user-email').value;
    const role=document.getElementById('user-role').value;
    const status= document.getElementById('user-status').value;
    const photoInput = document.getElementById('photo');
    const photoFile = photoInput.files[0];

    if(!name || !email || !role || !status){
        alert("Please provid valid input data");
        return
    }
    const formData = new FormData();
    formData.append('name', name);
    formData.append('email', email);
    formData.append('role', role);
    formData.append('status', status);
    formData.append('photo', photoFile);

    try {
        const response = await fetch(`/users`, {
            method: 'POST',
            body: formData,
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
        await fetch(`/users/${id}`, { method: 'DELETE' });
        alert("User deleted.")
        loadUsers();
    } else {
        alert("Deletion cancelled.");
    }

}

async function updateUser(id) {
    const name = prompt('Enter new name:');
    const email = prompt('Enter new email:');
    const role=prompt('Enter new role');
    const status=prompt('Enter new status:');
    const photoInput=prompt('Upload a new photo:');

    const formData = new FormData();
    formData.append('name', name);
    formData.append('email', email);
    formData.append('role', role);
    formData.append('status', status);
    formData.append('photo', photoInput);

    if (!name || !email) return;
    await fetch(`/users/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: formData,
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
        const response = await fetch(`/users/${userId}`);
        if (!response.ok) throw new Error('User not found.');
        const user = await response.json();

        resultDiv.innerHTML = `
        <div class="alert alert-info">
            <p><strong>ID:</strong> ${user.id}</p>
            <p><strong>Name:</strong> ${user.name}</p>
            <p><strong>Email:</strong> ${user.email}</p>
            <p><strong>Status:</strong> ${user.status}</p>
        </div>
    `;
    } catch (error) {
        console.error(error);
        resultDiv.innerHTML = `
        <div class="alert alert-danger">Error: User not found.</div>
    `;
    }
});

function sendEmail(userId) {
    window.location.href = `/mail?userId=${userId}`;
}
