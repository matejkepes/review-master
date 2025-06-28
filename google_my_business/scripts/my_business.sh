#!/bin/bash
# This should be run using cron probably once a day with a line in /etc/contab example:
# 30 14	* * *	root	/home/ubuntu/Documents/code/golang/google_my_business/my_business.sh
cd /home/ubuntu/Documents/code/golang/google_my_business
./my_business
