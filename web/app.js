// Check if config exists first
fetch('/config-exists')
    .then(response => response.json())
    .then(data => {
        if (data.exists) {
            // Populate fields from existing config
            fetch('/get-config')
                .then(response => response.json())
                .then(config => populateForm(config));
        }
    });

// Function to populate form
function populateForm(config) {
    document.getElementById('url').value = config.homeassistant.url;
    document.getElementById('port').value = config.homeassistant.port;
    // Populate other fields
}

// Function to submit form
function submitConfig() {
    const formData = {
        homeassistant: {
            url: document.getElementById('url').value,
            port: document.getElementById('port').value,
            // ...other fields
        },
        // ...other sections
    };
    fetch('/save-config', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(formData),
    })
    .then(response => response.json())
    .then(data => alert(data.message));
}
