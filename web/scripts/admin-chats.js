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
                    <td id="status-${chat.chat_uuid}">${chat.status}</td>
                    <td>${new Date(chat.created_at).toLocaleString()}</td>
                    <td>
                        <a href="/admin/chats/${chat.chat_uuid}" class="btn btn-primary btn-sm">Open Chat</a>
                        <button class="btn btn-danger btn-sm" onclick="deactivateChat('${chat.chat_uuid}')">Deactivate</button>
                    </td>
                `;
                chatList.appendChild(row);
            });
        })
        .catch(error => console.error("Error fetching chats:", error));
});

function deactivateChat(chatUUID) {
    fetch(`/api/chat/close/${chatUUID}`, {
        method: "PUT",
        credentials: "include",
        headers: { "Content-Type": "application/json" }
    })
    .then(response => {
        if (!response.ok) throw new Error("Failed to deactivate chat");
        document.getElementById(`status-${chatUUID}`).textContent = "inactive";
    })
    .catch(error => console.error("Error deactivating chat:", error));
}
