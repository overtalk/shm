#!/bin/bash
out=$1
gcc -g -v  shmthex.c -lpthread -o ${out}

