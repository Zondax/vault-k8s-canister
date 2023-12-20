#!/bin/bash

temp_file=$(mktemp)

echo "Creating tunnel for " $URL

# Start the Cloudflare tunnel in the background and redirect output to a temporary file
cloudflared tunnel --url ${URL} > "$temp_file" 2>&1 &

# Get the PID of the cloudflared process
cloudflared_pid=$!

# Function to kill cloudflared and delete the filename when exiting
cleanup() {
    echo "Stopping cloudflared..."
    kill $cloudflared_pid
    echo "Deleting $TUNNEL_URL_FILENAME..."
    rm $TUNNEL_URL_FILENAME
}

# Trap script exit
trap cleanup EXIT

# Wait for the URL to appear in the output
while : ; do
    if grep -q 'https://[a-zA-Z0-9.-]*\.trycloudflare.com' "$temp_file"; then
        url=$(grep -o 'https://[a-zA-Z0-9.-]*\.trycloudflare.com' "$temp_file")
        echo "$url" > $TUNNEL_URL_FILENAME
        echo "Tunnel URL written to $TUNNEL_URL_FILENAME... url $url"
        break
    fi
    sleep 1
done

# Log everything after the tunnel connection is registered
tail -f "$temp_file"

# Keep the script running to maintain the tunnel
wait $cloudflared_pid