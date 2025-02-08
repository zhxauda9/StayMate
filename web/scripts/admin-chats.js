document.addEventListener("DOMContentLoaded", () => {
    fetch("/api/admin/chats")
        .then(response => response.json())
        .then(data => {
            const chatList = document.getElementById("chat-list");
            chatList.innerHTML = ""; 

            data.forEach((chat, index) => {
                const row = document.createElement("tr");
                row.innerHTML = `
                    <td>${index + 1}</td>
                    <td>${chat.chat_uuid}</td>
                    <td>${chat.user_id}</td>
                    <td>${chat.status}</td>
                    <td>${new Date(chat.created_at).toLocaleString()}</td>
                    <td>
                        <a href="/admin/chats/${chat.chat_uuid}" class="btn btn-primary btn-sm">Открыть чат</a>
                    </td>
                `;
                chatList.appendChild(row);
            });
        })
        .catch(error => console.error("Error fetching chats:", error));
});
