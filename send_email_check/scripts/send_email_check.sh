#!/bin/bash
# This should be run using cron every day with a line in /etc/contab example:
# 0 8	* * *	root	/home/ubuntu/Documents/code/golang/send_email_check/send_email_check.sh
cd /home/ubuntu/Documents/code/golang/send_email_check
./send_email_check
