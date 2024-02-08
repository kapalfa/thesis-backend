#!/bin/bash

SECRET_KEY=$(openssl rand -base64 32)

echo "Ssecret key: $SECRET_KEY"