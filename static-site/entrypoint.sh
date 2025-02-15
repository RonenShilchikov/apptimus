#!/bin/sh
# Run the static site generator
/app/static-generator

# Start nginx
nginx -g "daemon off;"