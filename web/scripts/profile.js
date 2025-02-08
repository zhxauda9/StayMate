let Email;

async function fetchProfile() {
    try {
        const response = await fetch('/auth/profile', {
            method: 'GET',
            credentials: 'include',
        });

        if (!response.ok) throw new Error('Failed to fetch profile data');

        const userData = await response.json();
        Email = userData.email;

        document.getElementById('profilePhoto').src = userData.photo || 'default-photo.jpg';
        document.getElementById('profileName').textContent = userData.name || 'Your Name';
        document.getElementById('profileEmail').textContent = `Email: ${userData.email || 'user@example.com'}`;
        document.getElementById('profileStatus').textContent = userData.status || 'Our honoured guest';
        document.getElementById('editName').value = userData.name;
        document.getElementById('editEmail').value = userData.email;

        connectWebSocket(); // Подключаем WebSocket после загрузки профиля
    } catch (error) {
        console.error('Error:', error);
    }
}

async function saveProfile() {
    const name = document.getElementById('editName').value;
    const photo = document.getElementById('editPhoto').files[0];

    const formData = new FormData();
    formData.append('name', name);
    if (photo) formData.append('photo', photo);

    try {
        const response = await fetch('/auth/profile', {
            method: 'PUT',
            credentials: 'include',
            body: formData,
        });

        if (!response.ok) throw new Error('Failed to save profile');

        await fetchProfile();
        alert('Profile updated successfully');
        $('#editProfileModal').modal('hide');
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to update profile');
    }
}

function logout() {
    document.cookie = 'Authorization=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/; Secure; SameSite=Lax';
    localStorage.clear();

    if ('caches' in window) {
        caches.keys().then(cacheNames => cacheNames.forEach(cacheName => caches.delete(cacheName)));
    }

    fetch('/auth/logout', { method: 'POST', credentials: 'include' })
        .then(() => window.location.href = '/login')
        .catch(err => {
            console.error('Logout failed:', err);
            alert('Failed to logout');
        });
}

function toggleChat() {
    const chatBox = document.getElementById('chat-box');
    chatBox.style.display = chatBox.style.display === 'block' ? 'none' : 'block';

    if (chatBox.style.display === 'block') {
        connectWebSocket();
    } else {
        disconnectWebSocket();
    }
}

let socket;

function connectWebSocket() {
    if (!Email) {
        console.error("User email not found");
        return;
    }

    if (socket && socket.readyState === WebSocket.OPEN) {
        console.log("WebSocket already connected.");
        return;
    }

    socket = new WebSocket(`ws://localhost:8080/ws/user?userID=${encodeURIComponent(Email)}`);

    socket.onopen = () => console.log("Connected to WebSocket");

    socket.onmessage = event => {
        const chatMessages = document.getElementById('chat-messages');
        const messageData = event.data;
        const message = messageData.split(':')[1];

        const adminMessage = document.createElement('div');
        adminMessage.textContent = `Admin: ${message}`;
        adminMessage.classList.add('admin-message');
        chatMessages.appendChild(adminMessage);

        chatMessages.scrollTop = chatMessages.scrollHeight; // Авто-прокрутка вниз
    };

    socket.onclose = () => console.log("Disconnected from WebSocket");
}

function disconnectWebSocket() {
    if (socket) {
        socket.close();
        socket = null;
    }
}

function sendMessage(event) {
    if (event && event.key && event.key !== "Enter") return;

    const messageInput = document.getElementById('message');
    const message = messageInput.value.trim();

    if (!message) return;

    const chatMessages = document.getElementById('chat-messages');

    const userMessage = document.createElement('div');
    userMessage.textContent = `You: ${message}`;
    userMessage.classList.add('user-message');
    chatMessages.appendChild(userMessage);

    messageInput.value = "";

    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(message);
    } else {
        console.error("WebSocket is not connected.");
    }

    chatMessages.scrollTop = chatMessages.scrollHeight; // Авто-прокрутка вниз
}

document.getElementById('message').addEventListener('keypress', sendMessage);
document.getElementById('send-button').addEventListener('click', sendMessage);
document.getElementById('chat-icon').addEventListener('click', function () {
    document.getElementById('chat-box').style.display = 'block';
});

document.getElementById('close-chat').addEventListener('click', function () {
    document.getElementById('chat-box').style.display = 'none';
});

window.onload = fetchProfile;
