#!/bin/bash
# This should be run using cron probably once a month with a line in /etc/contab example:
# 30 10	1 * *	root	/home/ubuntu/Documents/code/golang/google_my_business/my_business_report.sh
cd /home/ubuntu/Documents/code/golang/google_my_business
./my_business -reportmonthback 1
