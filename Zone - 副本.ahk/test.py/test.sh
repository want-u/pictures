#! /bin/bash
########################################
# S-Data Monitor System Package Script #
# Create Time: 2021-11                 #
# Powered by S-Data                    #
# web: http://www.s-data.cn            #
########################################

# create sdata user
echo "Create User and Group: sdata"
getent group sdata &> /dev/null || groupadd --system -g 10050 sdata
getent passwd sdata &> /dev/null || useradd --system -u 10050 -g sdata -d /usr/lib/sdata -s /sbin/nologin -c "S-Data Monitoring System" sdata

# create lib directory
test -d /usr/lib/sdata || mkdir -m u=rwx,g=rwx,o= -p /usr/lib/sdata
chown sdata:sdata /usr/lib/sdata -R

#! /bin/bash
########################################
# S-Data Monitor System Package Script #
# Create Time: 2021-11                 #
# Powered by S-Data                    #
# web: http://www.s-data.cn            #
########################################

test -d /var/run/sdata || mkdir /var/run/sdata
test -d /var/log/sdata || mkdir /var/log/sdata
test -f /var/log/sdata/redis.log || touch /var/log/sdata/redis.log

/bin/touch /var/log/sdata/mysqld.log
chown -R sdata:sdata /usr/lib/sdata /usr/local/sdata /var/log/sdata /var/run/sdata/
test -f /etc/sysctl.conf || touch /etc/sysctl.conf
sysctl -p -q
ldconfig &> /dev/null
/usr/local/sdata/mysql/bin/mysql_install_db --basedir=/usr/local/sdata/mysql --datadir=/usr/local/sdata/mysql/data --defaults-file=/usr/local/sdata/mysql/conf/my.cnf --user=sdata &>/dev/null

systemctl enable mariadb &>> /var/log/sdata/system.log
systemctl start mariadb &>> /var/log/sdata/system.log

test -d /var/lib/mysql || mkdir /var/lib/mysql
chown -R sdata:sdata /var/lib/mysql
ln -sf /var/run/sdata/mysql.sock /var/lib/mysql/mysql.sock

#/usr/local/sdata/mysql/bin/mysqladmin -u root -h 127.0.0.1 password 'sData#888'
/usr/local/sdata/mysql/bin/mysqladmin -u root password 'sData#888'
/usr/local/sdata/mysql/bin/mysql -uroot -p'sData#888' -e "ALTER USER 'root'@'localhost' IDENTIFIED BY 'sData#888';"
/usr/local/sdata/mysql/bin/mysql -uroot -p'sData#888' -e "create database zabbix;"
#zcat /usr/local/sdata/mysql/share/create.sql.gz | /usr/local/sdata/mysql/bin/mysql -uroot -p'sData#888' -D zabbix >> /var/log/sdata/system.log
/usr/local/sdata/mysql/bin/mysql -uroot -h127.0.0.1 -p'sData#888' -e "set password for 'sdata'@'localhost' = password('sData#888');" >> /var/log/sdata/system.log
/usr/local/sdata/mysql/bin/mysql -uroot -h127.0.0.1 -p'sData#888' -e "grant all privileges on *.* to 'sdata'@'%' identified by 'sData#888';" >> /var/log/sdata/system.log
/usr/local/sdata/mysql/bin/mysql -uroot -h127.0.0.1 -p'sData#888' -D zabbix < /usr/local/sdata/mysql/share/zabbix_create_5.0.17.sql >> /var/log/sdata/system.log
/usr/local/sdata/mysql/bin/mysql -uroot -p'sData#888' -e "drop user ''@'localhost';"
#/usr/local/sdata/mysql/bin/mysql -uroot -p'sData#888' -e "ALTER TABLE zabbix.dbversion ADD id INT(16) NOT NULL PRIMARY KEY AUTO_INCREMENT FIRST,ADD num INT(16) DEFAULT 0;" >> /var/log/sdata/system.log
/usr/local/sdata/mysql/bin/mysql -uroot -p'sData#888' -e "flush privileges;" >> /var/log/sdata/system.log
/usr/local/sdata/mysql/bin/mysql -uroot -p'sData#888' zabbix < /usr/local/sdata/mysql/share/partition_all.sql >> /var/log/sdata/system.log
/usr/local/sdata/mysql/bin/mysql -usdata -p'sData#888' zabbix -e "CALL partition_maintenance_all('"'zabbix'"');" >> /var/log/sdata/system.log
#echo '1 0 * * * /usr/local/sdata/mysql/bin/mysql -usdata -p'sData#888' zabbix -e "CALL partition_maintenance_all('"'zabbix'"');"' >> /var/spool/cron/root
echo "Done."