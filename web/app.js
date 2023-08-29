// web/app.js
document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('configForm');

    form.addEventListener('submit', function(event) {
        event.preventDefault();
        
        const username = document.getElementById('username').value;
        const enableFeatureX = document.getElementById('enableFeatureX').checked;

        const configData = {
            username,
            enableFeatureX,
            // Add more fields here
        };

        // Send the data to the API
        fetch('/api/config', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(configData)
        }).then(response => response.json()).then(data => {
            console.log("Config saved:", data);
        });
    });
});
