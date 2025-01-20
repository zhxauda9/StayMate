let currentPage = 1;
const limit = 10;
const prevBtn = document.getElementById("prevBtn");
const nextBtn = document.getElementById("nextBtn");
const currentPageSpan = document.getElementById("currentPage");

function updateButtons() {
    prevBtn.disabled = currentPage === 1;
}

async function loadBookings(filterStart='',filterEnd='',sort='') {
    try {
        let url = `/bookings?filterStart=${filterStart}&filterEnd=${filterEnd}&limit=${limit}&page=${currentPage}&sort=${sort}`;
        const response = await fetch(url);
        if (!response.ok) {
            throw new Error('Не удалось загрузить бронирования.');
        }
        const bookings = await response.json();
        const table = document.getElementById('bookings-table');
        table.innerHTML = '';
        bookings.forEach(booking => {
            const row = document.createElement('tr');
            row.innerHTML = `
            <td>${booking.id}</td>
            <td>${booking.user_id}</td>
            <td>${booking.room_id}</td>
            <td>${booking.check_in}</td>
            <td>${booking.check_out}</td>
            <td>
                <button class="btn btn-warning btn-sm" onclick="updateBooking(${booking.id})">Update</button>
                <button class="btn btn-danger btn-sm" onclick="deleteBooking(${booking.id})">Delete</button>
            </td>
        `;
            table.appendChild(row);
        });
        currentPageSpan.textContent = `Page ${currentPage}`;
    } catch (error) {
        console.error(error);
        alert("Ошибка загрузки бронирований.");
    }
}

prevBtn.addEventListener("click", () => {
    if (currentPage > 1) {
        currentPage--;
        loadBookings();
        updateButtons();
    }
});

nextBtn.addEventListener("click", () => {
    currentPage++;
    loadBookings();
    updateButtons();
});

loadBookings();

document.getElementById('filter-form').addEventListener('submit', (e) => {
    e.preventDefault();
    const filterStart = document.getElementById('filterStart').value.trim();
    const filterEnd=document.getElementById('filterEnd').value.trim();
    loadBookings(filterStart,filterEnd,'');
});

document.getElementById('sort-form').addEventListener('submit', (e) => {
    e.preventDefault();
    const sortSelect = document.getElementById('sort').value.trim();
    loadBookings('', '', sortSelect);
});


function isValidDate(date) {
    return !isNaN(Date.parse(date));
}

document.getElementById('create-booking-form').addEventListener('submit', async (e) => {
    e.preventDefault();

    const userId = document.getElementById('user-id').value;
    const roomId = document.getElementById('room-id').value;
    const checkIn = document.getElementById('check-in').value;
    const checkOut = document.getElementById('check-out').value;

    if (!userId || !roomId || !isValidDate(checkIn) || !isValidDate(checkOut)) {
        alert("Please provide valid input data.");
        return;
    }

    const booking = {
        user_id: parseInt(userId),
        room_id: parseInt(roomId),
        check_in: new Date(checkIn).toISOString(),
        check_out: new Date(checkOut).toISOString(),
    };

    try {
        const response = await fetch(`/bookings`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(booking),
        });

        if (!response.ok) {
            throw new Error("Failed to create booking.");
        }

        alert("Booking created successfully!");
        loadBookings();
    } catch (error) {
        console.error(error);
        alert("Error: Failed to create booking.");
    }
});

async function deleteBooking(id) {
    if (confirm("Are you sure you want to delete the booking?")) {
        try {
            const response = await fetch(`/bookings/${id}`, { method: 'DELETE' });

            if (!response.ok) {
                throw new Error('Failed to delete booking.');
            }

            alert("Booking deleted.");
            loadBookings();
        } catch (error) {
            console.error(error);
            alert("Error deleting booking.");
        }
    } else {
        alert("Deletion cancelled.");
    }
}

async function updateBooking(id) {
    const userID = prompt('Enter new user ID:');
    const roomID = prompt('Enter new room ID:');
    const checkIn = prompt('Enter new check-in date (YYYY-MM-DD):');
    const checkOut = prompt('Enter new check-out date (YYYY-MM-DD):');

    if (!isValidDate(checkIn) || !isValidDate(checkOut)) {
        alert("Invalid date format.");
        return;
    }

    const booking = {
        user_id: parseInt(userID),
        room_id: parseInt(roomID),
        check_in: new Date(checkIn).toISOString(),
        check_out: new Date(checkOut).toISOString(),
    };

    try {
        const response = await fetch(`/bookings/${id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(booking),
        });

        if (!response.ok) {
            throw new Error("Failed to update booking.");
        }

        alert("Booking updated successfully!");
        loadBookings();
    } catch (error) {
        console.error(error);
        alert("Error: Failed to update booking.");
    }
}