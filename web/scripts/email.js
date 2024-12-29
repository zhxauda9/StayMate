document.getElementById('send-email-form').addEventListener('submit', function (event) {
    event.preventDefault(); // Prevent form submission

    // Get form values
    const recipientSelect = document.getElementById('email-recipients');
    const selectedOption = recipientSelect.options[recipientSelect.selectedIndex];
    const email = selectedOption.text.split('(')[1].split(')')[0]; // Extract email from the option text
    const subject = document.getElementById('subject').value;
    const message = document.getElementById('message').value;

    // Get the file input element
    const fileInput = document.getElementById('file');
    const file = fileInput.files[0]; // Get the first selected file

    // Create the POST request payload
    const formData = new FormData();
    formData.append('emails', email);
    formData.append('subject', subject);
    formData.append('message', message);

    // If a file is attached, add it to the form data
    if (file) {
        formData.append('file', file);
    }

    // Choose the correct endpoint based on whether a file is attached
    const endpoint = file ? '/api/mailFile' : '/api/mail';

    // Send POST request to the correct endpoint
    fetch(endpoint, {
        method: 'POST',
        body: formData
    })
    .then(response => response.text()) // Assuming the server responds with plain text
    .then(data => {
        alert(data); 
    })
    .catch(error => {
        console.error('Error sending email:', error);
        alert('Failed to send email.');
    });
});
