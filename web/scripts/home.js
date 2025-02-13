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


function getCookieValue(name) {
    const matches = document.cookie.match(new RegExp(
        '(?:^|; )' + name.replace(/([.$?*|{}()[]\\\/+^])/g, '\\$1') + '=([^;]*)'
    ));
    return matches ? decodeURIComponent(matches[1]) : undefined;
    }

    window.addEventListener('DOMContentLoaded', () => {
    const token = getCookieValue('Authorization');
    const loginBtn = document.getElementById('loginBtn');
    const registerBtn = document.getElementById('registerBtn');
    const profileBtn = document.getElementById('profileBtn');

    if (token) {
        // Hide Login & Register, show Profile
        loginBtn.style.display = 'none';
        registerBtn.style.display = 'none';
        profileBtn.style.display = 'inline-block';
    }
});

