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
        card.style = 'width: 18rem; border-radius: 15px; overflow: hidden; box-shadow: 0 5px 15px rgba(0, 0, 0, 0.11);';
        card.setAttribute('data-room-id', room.id);  // Add the room ID as a data attribute

        card.innerHTML = `
            <img src="${room.photo}" class="card-img-top" alt="Room Photo" style="height: 200px; object-fit: cover;">
            <div class="card-body text-center">
                <h3 class="card-title">Room ${room.number}</h3>
                <p class="card-text font-weight-bolder" style="color: crimson;">${room.price.toLocaleString()} $</p>
                <p class="card-text">Status: ${room.status}</p>
                <p class="card-text">Room class - ${room.class}</p>
                <p class="card-text">${room.description}</p>
                <button class="btn text-dark rounded book-now-btn" style="background-color: rgb(228, 213, 130)">Book Now</button>
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

document.addEventListener("click", function (event) {
  if (event.target.classList.contains("book-now-btn")) {
    $('#bookingModal').modal('show'); // Используем Bootstrap для показа модального окна
  }
});

document.addEventListener("click", async function (event) {
  if (event.target.classList.contains("book-now-btn")) {
    try {
      // Fetch profile data
      const profileResponse = await fetch('/auth/profile');
      if (!profileResponse.ok) {
        throw new Error('Failed to fetch profile data');
      }
      const profileData = await profileResponse.json();

      if (!profileData.email || !profileData.id) {
        alert('Failed to retrieve user information. Please log in again.');
        return;
      }

      // Populate hidden form fields
      document.getElementById('userID').value = profileData.id;
      document.getElementById('email').value = profileData.email;

      // Get roomID and price from the card
      const roomCard = event.target.closest('.card');
      const roomID = roomCard.getAttribute('data-room-id');  // Use the data-room-id attribute
      const price = parseFloat(roomCard.querySelector('.card-text.font-weight-bolder').innerText.replace('$', '').replace(',', ''));

      // Populate hidden form fields for roomID and price
      document.getElementById('roomID').value = roomID;
      document.getElementById('price').value = price;

      // Show the modal
      $('#bookingModal').modal('show');
    } catch (error) {
      console.error('Error fetching profile data:', error);
      alert('An error occurred while fetching profile data.');
    }
  }
});

document.getElementById("bookingForm").addEventListener("submit", async function (event) {
  event.preventDefault();

  // Собираем данные из формы
  const userID = parseInt(document.getElementById("userID").value);
  const email = document.getElementById("email").value;
  const address = document.getElementById("address").value;
  const cardNumber = document.getElementById("cardNumber").value;
  const expirationDate = document.getElementById("expirationDate").value;
  const cvv = document.getElementById("cvv").value;
  const roomID = parseInt(document.getElementById("roomID").value);
  const price = parseFloat(document.getElementById("price").value);
  const checkIn = document.getElementById("checkIn").value;
  const checkOut = document.getElementById("checkOut").value;

  const bookingData = {
    user_id: userID,
    email: email,
    address: address,
    card_number: cardNumber,
    expiration_date: expirationDate,
    cvv: cvv,
    room_id: roomID,
    price: price,
    check_in: checkIn,
    check_out: checkOut,
  };

  try {
    // Отправляем POST-запрос на хэндлер /api/v2/bookings
    const response = await fetch('/api/v2/bookings', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(bookingData),
    });

    const result = await response.json();

    if (response.ok) {
      alert(result.message);  // Успешный ответ
    } else {
      alert(`Error: ${result.message}`); // Ошибка
    }
  } catch (error) {
    console.error('Error processing booking:', error);
    alert('An unexpected error occurred.');
  }

  // Закрываем модальное окно
  $('#bookingModal').modal('hide');
});


