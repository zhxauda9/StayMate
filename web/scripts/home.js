let currentPage = 1;
const limit = 10;
const prevBtn = document.getElementById("prevBtn");
const nextBtn = document.getElementById("nextBtn");
const currentPageSpan = document.getElementById("currentPage");


function updateButtons() {
    prevBtn.disabled = currentPage === 1;
}

async function loadRooms(filterStart = '', filterEnd = '', sort = '') {
    try {
        let url = `/rooms?filterStart=${filterStart}&filterEnd=${filterEnd}&limit=${limit}&page=${currentPage}&sort=${sort}`;
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error('Failed to load rooms');
        }

        const rooms = await response.json();
        const cardsContainer = document.getElementById('rooms-cards');
        cardsContainer.innerHTML = '';

        rooms.forEach(room => {
            const card = document.createElement('div');
            card.className = 'card m-4 grow';
            card.style = 'width: 18rem; border-radius: 15px; overflow: hidden;box-shadow: 0 5px 15px rgba(0, 0, 0, 0.11);';

            card.innerHTML = `
                <img src="${room.photo}" class="card-img-top" alt="Room Photo" style="height: 200px; object-fit: cover;">
                <div class="card-body text-center">
                    <h3 class="card-title">Room ${room.number}</h3>
                    <p class="card-text font-weight-bolder" style="color: crimson;">${room.price.toLocaleString()} ₸</p>
                    <p class="card-text">Status: ${room.status}</p>
                    <p class="card-text">Room class - ${room.class}</p>
                    <p class="card-text">${room.description}</p>
                    <button class="btn text-dark rounded" style="background-color:rgb(228, 213, 130)">Book
            Now</button>
                </div>
            `;

            cardsContainer.appendChild(card);
        });
        currentPageSpan.textContent = `Page ${currentPage}`;
    } catch (error) {
        console.error(error);
        alert("Failed to load rooms");
    }
}

prevBtn.addEventListener("click", () => {
    if (currentPage > 1) {
        currentPage--;
        loadRooms();
        updateButtons();
    }
});

nextBtn.addEventListener("click", () => {
    currentPage++;
    loadRooms();
    updateButtons();
});

loadRooms();

document.getElementById('filter-sort').addEventListener('submit', (e) => {
    e.preventDefault();
    const filterStart = document.getElementById('filterStart').value.trim();
    const filterEnd = document.getElementById('filterEnd').value.trim();
    const sortSelect = document.getElementById('sort').value.trim();
    loadRooms(filterStart, filterEnd, sortSelect);
});

document.getElementById('contact-email-form').addEventListener('submit', function (event) {
    event.preventDefault();

    const email = document.getElementById('user-email').value;

    const formData = new FormData();
    formData.append('emails', email);
    formData.append('subject',"Hi, how can we help you?");
    formData.append('message',"Thank you for reaching out! How can we assist you?")
    fetch('/api/mail', {
        method: 'POST',
        body: formData
    })
        .then(response => response.text())

        .catch(error => {
            console.error('Error sending email:', error);
        });
});


// chat 
function toggleChat() {
    const chatBox = document.getElementById('chat-box');
    chatBox.style.display = chatBox.style.display === 'block' ? 'none' : 'block';
}

function sendMessage(event) {
    // // Отправка сообщения при нажатии на Enter
    // if (event.key && event.key !== "Enter") return;

    // const messageInput = document.getElementById('message');
    // const message = messageInput.value.trim();

    // if (message) {
    //     const chatMessages = document.getElementById('chat-messages');

    //     // Добавляем сообщение пользователя в чат
    //     const userMessage = document.createElement('div');
    //     userMessage.textContent = `You: ${message}`;
    //     userMessage.style.marginBottom = "10px";
    //     chatMessages.appendChild(userMessage);

    //     messageInput.value = "";

    //     // Здесь можно добавить отправку сообщения на сервер
    //     // Пример:
    //     // const socket = new WebSocket("ws://localhost:8080/ws");
    //     // socket.send(JSON.stringify({ sender: "user", content: message }));
    // }
}


