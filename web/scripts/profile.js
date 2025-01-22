async function fetchProfile() {
    try {
        const response = await fetch('/api/profile', {
            method: 'GET',
            credentials: 'include',
        });

        if (!response.ok) {
            throw new Error('Failed to fetch profile data');
        }

        const userData = await response.json();
        document.getElementById('profilePhoto').src = userData.photo || 'default-photo.jpg';
        document.getElementById('profileName').textContent = userData.name;
        document.getElementById('profileEmail').textContent = `Email: ${userData.email}`;
        document.getElementById('profileStatus').textContent = userData.status || 'Our honoured guest';

        document.getElementById('editName').value = userData.name;
        document.getElementById('editEmail').value = userData.email;
    } catch (error) {
        console.error('Error:', error);
        document.body.innerHTML = `<p>Error loading profile</p>`;
    }
}

async function saveProfile() {
    const name = document.getElementById('editName').value;
    const photo = document.getElementById('editPhoto').files[0];

    const formData = new FormData();
    formData.append('name', name);
    if (photo) {
        formData.append('photo', photo);
    }

    try {
        const response = await fetch('/api/profile', {
            method: 'PUT',
            credentials: 'include',
            body: formData,
        });

        if (!response.ok) {
            throw new Error('Failed to save profile');
        }

        fetchProfile();
        alert('Profile updated successfully');
        $('#editProfileModal').modal('hide');
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to update profile');
    }
}

function logout() {
    document.cookie = 'Authorization=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/; Secure; SameSite=Lax';

    localStorage.clear();

    if ('caches' in window) {
        caches.keys().then(cacheNames => {
            cacheNames.forEach(cacheName => caches.delete(cacheName));
        });
    }

    fetch('/api/logout', { method: 'POST', credentials: 'include' })
        .then(() => {
            window.location.href = '/login';
        })
        .catch(err => {
            console.error('Logout failed:', err);
            alert('Failed to logout');
        });
}


window.onload = fetchProfile;