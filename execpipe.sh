#!/bin/bash 

while true; do 
  printf "==== Listening for the job ====\n"
  sh -c "printf '==== Job started ====\n'; $(cat ./mypipe); printf '==== Job done ====\n'"
  printf "\n"
done
