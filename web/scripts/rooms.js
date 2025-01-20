let currentPage = 1;
const limit = 10;
const prevBtn = document.getElementById("prevBtn");
const nextBtn = document.getElementById("nextBtn");
const currentPageSpan = document.getElementById("currentPage");


function updateButtons() {
    prevBtn.disabled = currentPage === 1;
}

async function loadRooms(filterStart='',filterEnd='',sort='') {
    try {
        let url = `/rooms?filterStart=${filterStart}&filterEnd=${filterEnd}&limit=${limit}&page=${currentPage}&sort=${sort}`;
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error('Failed to load rooms');
        }
        const rooms = await response.json();
        const table = document.getElementById('rooms-table');
        table.innerHTML = '';
        rooms.forEach(room => {
            const row = document.createElement('tr');
            row.innerHTML = `
            <td>${room.id}</td>
            <td>${room.number}</td>
            <td>${room.class}</td>
            <td>${room.price}</td>
            <td>${room.status}</td>
            <td>
                <button class="btn btn-warning btn-sm" onclick="updateRoom(${room.id})">Update</button>
                <button class="btn btn-danger btn-sm" onclick="deleteRoom(${room.id})">Delete</button>
            </td>
        `;
            table.appendChild(row);
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

document.getElementById('filter-form').addEventListener('submit', (e) => {
    e.preventDefault();
    const filterStart = document.getElementById('filterStart').value.trim();
    const filterEnd=document.getElementById('filterEnd').value.trim();
    loadRooms(filterStart,filterEnd,'');
});

document.getElementById('sort-form').addEventListener('submit', (e) => {
    e.preventDefault();
    const sortSelect = document.getElementById('sort').value.trim();
    loadRooms('','',sortSelect);
});

document.getElementById('create-room-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const number = document.getElementById('number').value;
    const classRoom = document.getElementById('class').value;
    const price = document.getElementById('price').value;
    const status=document.getElementById('status').value;

    if (!number || !classRoom || !price) {
        alert("Please provide valid input data.");
        return;
    }

    const room = {
        number: parseInt(number),
        class: classRoom,
        price: parseFloat(price),
        status:status,
    };

    try {
        const response = await fetch(`/rooms`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(room),
        });

        if (!response.ok) {
            throw new Error("Failed to create room.");
        }

        alert("Room created successfully!");
        loadRooms();
    } catch (error) {
        console.error(error);
        alert("Error: Failed to create room.");
    }
});

async function deleteRoom(id) {
    if (confirm("Are you sure you want to delete the room?")) {
        try {
            const response = await fetch(`/rooms/${id}`, { method: 'DELETE' });

            if (!response.ok) {
                throw new Error('Failed to delete room.');
            }
            alert("Room deleted.");
            loadRooms();
        } catch (error) {
            console.error(error);
            alert("Error deleting room.");
        }
    } else {
        alert("Deletion cancelled.");
    }
}

async function updateRoom(id) {
    const number = prompt('Enter new number:');
    const classRoom = prompt('Enter new class:');
    const price = prompt('Enter new price numeric(10,2):');
    const status=prompt('Enter new status:')

    if (!number || !classRoom || !price) {
        alert("Please provide valid input data.");
        return;
    }

    const room = {
        number: parseInt(number),
        class: classRoom,
        price: parseFloat(price),
        status:status,
    };

    try {
        const response = await fetch(`/rooms/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(room),
        });
        if (!response.ok) {
            throw new Error("Failed to update room.");
        }
        alert("Room updated successfully!");
        loadRooms();
    } catch (error) {
        console.error(error);
        alert("Error: Failed to update room.");
    }
}