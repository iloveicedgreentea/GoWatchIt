#!/bin/sh

# Start nginx
nginx -g "daemon off;" &

# Start the main app
/gowatchit