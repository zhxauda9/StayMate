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
            const errorText = await response.text();
            throw new Error(`Failed to load rooms: ${errorText}`);
        }

        const rooms = await response.json();
        const cardsContainer = document.getElementById('rooms-cards');
        cardsContainer.innerHTML = '';

        rooms.forEach(room => {
            const card = document.createElement('div');
            card.className = 'card m-4 grow';
            card.style = 'width: 18rem; border-radius: 15px; overflow: hidden;';

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
s
        currentPageSpan.textContent = `Page ${currentPage}`;
    } catch (error) {
        console.error(error);
        alert('Failed to load rooms');
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

