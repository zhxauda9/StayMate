document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("api/mail");

    form.addEventListener("submit", function (event) {
        event.preventDefault(); // Остановить стандартное поведение формы

        const recipientId = document.getElementById("email-recipients").value;
        const message = document.getElementById("message").value;
        const file = document.getElementById("file").files[0];

        // Добавь логику отправки email, например, через AJAX или fetch
        console.log(`Sending email to user ID: ${recipientId}, message: ${message}, file: ${file?.name}`);

        // Пример запроса на сервер (это зависит от того, как ты реализуешь backend)
        fetch("/send-email-submit", {
            method: "POST",
            body: new FormData(form)
        })
        .then(response => response.json())
        .then(data => {
            console.log("Email sent successfully:", data);
            alert("Email sent successfully!");
        })
        .catch(error => {
            console.error("Error sending email:", error);
            alert("Error sending email.");
        });
    });
});
