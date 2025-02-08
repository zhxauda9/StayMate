let Email;
let ChatUUID;
let socket;

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

        await initChat();
    } catch (error) {
        console.error('Error:', error);
    }
}

async function initChat() {
    const existingChatUUID = getCookie('admin_chat_uuid');
    if (existingChatUUID) {
        ChatUUID = existingChatUUID;
        await fetchChatHistory(existingChatUUID);
    } else {
        await startNewChat();
    }
    connectWebSocket();
}

async function fetchChatHistory(chatUUID) {
    try {
        const response = await fetch(`/api/chat/history/${chatUUID}`, {
            method: 'GET',
            credentials: 'include',
        });

        if (!response.ok) throw new Error('Failed to fetch chat history');

        const chatData = await response.json();
        ChatUUID = chatData.chat_uuid;
        renderChatHistory(chatData.messages);
    } catch (error) {
        console.error('Error:', error);
    }
}

async function startNewChat() {
    try {
        const response = await fetch('/api/chat/start', {
            method: 'POST',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email: Email }),
        });

        if (!response.ok) throw new Error('Failed to start a new chat');

        const chatData = await response.json();
        ChatUUID = chatData.chat_uuid;

        document.cookie = `admin_chat_uuid=${ChatUUID}; path=/; Secure; SameSite=Lax`;
    } catch (error) {
        console.error('Error:', error);
    }
}

function renderChatHistory(messages) {
    const chatMessages = document.getElementById('chat-messages');
    chatMessages.innerHTML = '';

    messages.forEach(msg => {
        const messageElement = document.createElement('div');
        messageElement.textContent = `${msg.sender}: ${msg.message}`;
        messageElement.classList.add(msg.sender === 'user' ? 'user-message' : 'admin-message');
        chatMessages.appendChild(messageElement);
    });

    chatMessages.scrollTop = chatMessages.scrollHeight;
}

function connectWebSocket() {
    if (!ChatUUID) {
        console.error("Chat UUID not found");
        return;
    }

    if (socket && socket.readyState === WebSocket.OPEN) {
        console.log("WebSocket already connected.");
        return;
    }

    socket = new WebSocket(`ws://localhost:8080/ws/user?userID=${encodeURIComponent(ChatUUID)}`);
    socket.onopen = () => console.log("Connected to WebSocket");

    socket.onmessage = event => {
        const chatMessages = document.getElementById('chat-messages');
        const messageData = event.data;
        const message = messageData.split(':')[1];

        const adminMessage = document.createElement('div');
        adminMessage.textContent = `Admin: ${message}`;
        adminMessage.classList.add('admin-message');
        chatMessages.appendChild(adminMessage);

        chatMessages.scrollTop = chatMessages.scrollHeight;
    };

    socket.onclose = () => console.log("Disconnected from WebSocket");
}

function sendMessage(event) {
    if (event && event.key && event.key !== "Enter") return;

    const messageInput = document.getElementById('message');
    const message = messageInput.value.trim();

    if (!message) return;

    const chatMessages = document.getElementById('chat-messages');
    const userMessage = document.createElement('div');
    userMessage.textContent = `${message}`;
    userMessage.classList.add('user-message');
    chatMessages.appendChild(userMessage);
    messageInput.value = "";

    if (socket && socket.readyState === WebSocket.OPEN) {
        socket.send(message);
    } else {
        console.error("WebSocket is not connected.");
    }

    chatMessages.scrollTop = chatMessages.scrollHeight;
}

function getCookie(name) {
    const match = document.cookie.match(new RegExp(`(^| )${name}=([^;]+)`));
    return match ? match[2] : null;
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
