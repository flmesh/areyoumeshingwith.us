#!/usr/bin/env bash

npm ci
if [ -z $CODESPACE_NAME ]; then
    npm run dev
else
    npm run csdev
fi
