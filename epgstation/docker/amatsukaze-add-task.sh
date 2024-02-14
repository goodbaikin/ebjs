#!/bin/bash
/usr/local/bin/amatsukaze-add-task -i "$1" -c $AMATSUKAZE_HOST -e Z: -r Z: -p "デフォルト"
sleep 10

# check encoding started
while true
do
  if [ -d "/app/recorded/succeeded" ]; then
    break
  else
    sleep 1
  fi
done

# wait for succeeded folder deletion
while true
do
  if [ -d "/app/recorded/succeeded" ]; then
    sleep 1
  else
    break
  fi
done

