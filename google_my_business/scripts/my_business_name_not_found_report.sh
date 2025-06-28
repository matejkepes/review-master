#!/bin/bash
# This should be run using cron probably twice a week with a line in /etc/contab example:
# 00 9	* * 1,4	root	/home/ubuntu/Documents/code/golang/google_my_business/my_business_name_not_found_report.sh
cd /home/ubuntu/Documents/code/golang/google_my_business
./my_business -reportnameorpostalcodenotfound true
