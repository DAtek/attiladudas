#!/bin/bash

API_PUBLIC_KEY=$(cat {{ public_key }}) \
API_PRIVATE_KEY=$(cat {{ private_key }}) \
{{ exec_start }}
