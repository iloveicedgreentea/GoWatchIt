#!/bin/sh

# Kill any existing nginx processes
pkill nginx || true

# Wait a moment for the port to be released
sleep 1

# Start nginx
nginx -g "daemon off;" &

# Start the main app
ls
pwd
stat "/gowatchit"
exec /gowatchit