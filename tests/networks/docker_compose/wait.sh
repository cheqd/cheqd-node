#!/bin/bash

for i in 1 2 3; do $($cmd) && break || sleep 15; done