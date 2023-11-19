document.addEventListener("DOMContentLoaded", function() {
   

    const form = document.getElementById('urlForm'); // selecting the form

    form.addEventListener('submit', function(event) {
        event.preventDefault(); // preventing the form from submitting the default way

        const urlInput = document.getElementById('url-to-shorten');
        const url = urlInput.value; // getting the value of the input field

        // Sending the URL to the server
        fetch('/shorten', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/x-www-form-urlencoded',
            },
            body: `url=${encodeURIComponent(url)}`, // sending the url data in the body
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.text(); // this should only receive the new shortened URL
        })
        .then((data) => {
            // Handle the response data here. ex: display the shortened URL to the user.
            console.log('Shortened URL:', data);
            // Display the shortened URL in the "result" paragraph
            document.getElementById('result').textContent = `Shortened URL: ${data}`;
        })
        .catch((error) => {
            console.error('Error:', error);
            // Handle errors here
        });
        
        
    });
});
