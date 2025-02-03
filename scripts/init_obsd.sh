#!/bin/sh

RCFILE="appd"

# Check if the RC file exists
if [ ! -f "$RCFILE" ]; then
    echo "Error: RC file $RCFILE not found"
    exit 1
fi

# Copy the RC file to /etc/rc.d/
cp "$RCFILE" /etc/rc.d/

# Make it executable
chmod 755 "/etc/rc.d/$RCFILE"

# Enable the service in rc.conf.local
rcctl enable "$RCFILE"

"Service $RCFILE has been installed and enabled"
