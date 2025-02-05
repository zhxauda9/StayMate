async function fetchProfile() {
    try {
        const response = await fetch('/auth/profile', {
            method: 'GET',
            credentials: 'include',
        });

        if (!response.ok) {
            throw new Error('Failed to fetch profile data');
        }

        const userData = await response.json();
        document.getElementById('profilePhoto').src = userData.photo || 'default-photo.jpg';
        document.getElementById('profileName').textContent = userData.name || 'Your Name';
        document.getElementById('profileEmail').textContent = `Email: ${userData.email || 'user@example.com'}`;
        document.getElementById('profileStatus').textContent = userData.status || 'Our honoured guest';
        document.getElementById('editName').value = userData.name;
        document.getElementById('editEmail').value = userData.email;
    } catch (error) {
        console.error('Error:', error);
    }
}

async function saveProfile() {
    const name = document.getElementById('editName').value;
    const photo = document.getElementById('editPhoto').files[0];

    const formData = new FormData();
    formData.append('name', name);
    if (photo) {
        formData.append('photo', photo);
    }

    try {
        const response = await fetch('/auth/profile', {
            method: 'PUT',
            credentials: 'include',
            body: formData,
        });

        if (!response.ok) {
            throw new Error('Failed to save profile');
        }

        fetchProfile();
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

// chat 
function toggleChat() {
    const chatBox = document.getElementById('chat-box');
    chatBox.style.display = chatBox.style.display === 'block' ? 'none' : 'block';
}

function sendMessage(event) {
    // Отправка сообщения при нажатии на Enter
    if (event.key && event.key !== "Enter") return;

    const messageInput = document.getElementById('message');
    const message = messageInput.value.trim();

    if (message) {
        const chatMessages = document.getElementById('chat-messages');

        // Добавляем сообщение пользователя в чат
        const userMessage = document.createElement('div');
        userMessage.textContent = `You: ${message}`;
        userMessage.style.marginBottom = "10px";
        chatMessages.appendChild(userMessage);

        messageInput.value = "";

        // Здесь можно добавить отправку сообщения на сервер
        // Пример:
        // const socket = new WebSocket("ws://localhost:8080/ws");
        // socket.send(JSON.stringify({ sender: "user", content: message }));
    }
}

window.onload = fetchProfile;