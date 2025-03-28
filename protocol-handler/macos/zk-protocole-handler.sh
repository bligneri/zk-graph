#!/bin/bash
# Extract the file path from the zk:// URL scheme
url="${1//zk:\/\//}"

# Escape single quotes in the URL (if needed)
escaped_url="${url//\'/\\'}"

# Open a new terminal window using Kitty and run the zk edit command
osascript <<EOF
tell application "kitty"
    activate
    do script "zk edit '$escaped_url'"
end tell
EOF

