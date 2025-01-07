const BASE_URL = 'http://localhost:8080';

async function loadRooms() {
    try {
        const response = await fetch(`${BASE_URL}/rooms`);
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
            <td>
                <button class="btn btn-warning btn-sm" onclick="updateRoom(${room.id})">Update</button>
                <button class="btn btn-danger btn-sm" onclick="deleteRoom(${room.id})">Delete</button>
            </td>
        `;
            table.appendChild(row);
        });
    } catch (error) {
        console.error(error);
        alert("Failed to load rooms");
    }
}

loadRooms();


document.getElementById('create-room-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const number = document.getElementById('number').value;
    const classRoom = document.getElementById('class').value;
    const price = document.getElementById('price').value;

    if (!number || !classRoom || !price) {
        alert("Please provide valid input data.");
        return;
    }

    const room = {
        number: parseInt(number),
        class: classRoom,
        price: parseFloat(price),
    };

    try {
        const response = await fetch(`${BASE_URL}/rooms`, {
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
            const response = await fetch(`${BASE_URL}/rooms/${id}`, { method: 'DELETE' });

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

    if (!number || !classRoom || !price) {
        alert("Please provide valid input data.");
        return;
    }

    const room = {
        number: parseInt(number),
        class: classRoom,
        price: parseFloat(price),
    };

    try {
        const response = await fetch(`${BASE_URL}/rooms/${id}`, {
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