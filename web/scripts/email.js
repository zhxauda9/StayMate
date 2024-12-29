document.getElementById('send-email-form').addEventListener('submit', function (event) {
    event.preventDefault(); // Prevent form submission

    // Get form values
    // Get the selected recipient's email by selecting the option element
    const recipientSelect = document.getElementById('email-recipients');
    const selectedOption = recipientSelect.options[recipientSelect.selectedIndex];
    const email = selectedOption.text.split('(')[1].split(')')[0]; // Extract email from the option text
    const subject = document.getElementById('subject').value;
    const message = document.getElementById('message').value;

    // Create the POST request payload
    const formData = new FormData();
    formData.append('email', email);
    formData.append('subject', subject);
    formData.append('message', message);

    // Send POST request to the server
    fetch('/api/mail', {
        method: 'POST',
        body: formData
    })
    .then(response => response.text()) // Assuming the server responds with plain text
    .then(data => {
        alert(data); // Alert with the server response (e.g., "Email sended successfully")
    })
    .catch(error => {
        console.error('Error sending email:', error);
        alert('Failed to send email. Please try again.');
    });
});
