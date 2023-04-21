#!/bin/bash

# This script will remove all networks that are not in use.
ps aux | grep "geth" | grep -v grep | awk '{print $2}' | xargs kill -9
