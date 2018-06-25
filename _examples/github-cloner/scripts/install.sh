#!/bin/bash

set -e

DIR="$( cd "$( dirname "$0" )" && pwd )"

if [ "$(uname -s)" == "Darwin" ]; then
  if [ "$(whoami)" == "root" ]; then
    TARGET_DIR="/Library/Google/Chrome/NativeMessagingHosts"
  else
    TARGET_DIR="$HOME/Library/Application Support/Google/Chrome/NativeMessagingHosts"
  fi
else
  if [ "$(whoami)" == "root" ]; then
    TARGET_DIR="/etc/opt/chrome/native-messaging-hosts"
  else
    TARGET_DIR="$HOME/.config/google-chrome/NativeMessagingHosts"
  fi
fi

HOST_NAME=com.sniperkit.snk.dev-assistant
PROG_NAME=snk-assistant

# Create directory to store native messaging host.
mkdir -p "$TARGET_DIR"

# Copy native messaging host manifest.
cp "$DIR/$HOST_NAME.json" "$TARGET_DIR"

# Update host path in the manifest.
if [ "$(whoami)" == "root" ]; then
  HOST_PATH=/usr/local/bin/${PROG_NAME}
else
  mkdir -p $HOME/.local/bin
  HOST_PATH=$HOME/.local/bin/snk-assistant
fi

if [ "$HGSANDBOX_ENV" == "development" ]; then
  HOST_PATH=`pwd`/${PROG_NAME}
fi

if hash go > /dev/null 2>&1 ; then
  echo "Build ${PROG_NAME} host app from source"
  if go build -o $HOST_PATH . ; then
		echo 'Build success'
  else
		echo 'Build failed.'
		exit 1
	fi
else
  echo "Copying pre-built ${PROG_NAME} binary (install golang compiler to build from source)"
  cp ./bin/${PROG_NAME} $HOST_PATH
fi

ESCAPED_HOST_PATH=${HOST_PATH////\\/}
sed -i -e "s/HOST_PATH_PLACEHOLDER/$ESCAPED_HOST_PATH/" "$TARGET_DIR/$HOST_NAME.json"

# set app id
[[ $EXTENSION_ID == '' ]] && echo 'ERROR: EXTENSION_ID environment variable is not set.' && exit 1
sed -i -e "s/EXTENSION_ID_PLACEHOLDER/$EXTENSION_ID/" "$TARGET_DIR/$HOST_NAME.json"

# Set permissions for the manifest so that all users can read it.
chmod o+r "$TARGET_DIR/$HOST_NAME.json"
echo "Native messaging host $HOST_NAME has been installed."

