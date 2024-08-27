#!/bin/sh
whoami_txt=$(who am i)
export LOGIN_IP=$(echo $whoami_txt | awk -F"[()]" '{print $2}')
export LOGIN_USER=$(echo $whoami_txt | awk '{print $1}')
export LOGIN_TTY=$(echo $whoami_txt | awk '{print $2}')

export PROMPT_COMMAND='{ 
  cmd=$(history 1 | { read x cmd; echo "$cmd"; });
  log_entry=$(date "+%Y-%m-%d %T ## `hostname` ## $LOGIN_USER ## $LOGIN_TTY ## $LOGIN_IP ## USER=$(whoami) PWD=$(pwd) CMD=$cmd");
  echo "$log_entry" | /usr/local/bin/history_log;
}'