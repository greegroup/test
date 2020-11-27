## Installing MySQL
Install MySQL by the following command (for dnf-enabled systems, replace `yum` in the command with `dnf`):
```
shell> sudo yum install mysql-community-server
```
This installs the package for MySQL server (`mysql-community-server`) and also packages for the components required to run the server, including packages for the client (`mysql-community-client`), the common error messages and character sets for client and server (`mysql-community-common`), and the shared client libraries (`mysql-community-libs`).

## Starting the MySQL Server
Start the MySQL server with the following command:
```
shell> systemctl start mysqld
```
You can check the status of the MySQL server with the following command:
```
shell> systemctl status mysqld
```
If the operating system is systemd enabled, standard `systemctl` (or alternatively, service with the arguments reversed) commands such as `stop`, `start`, `status`, and `restart` should be used to manage the MySQL server service. The mysqld service is enabled by default, and it starts at system reboot. 

At the initial start up of the server, the following happens, given that the data directory of the server is empty:

- The server is initialized.

- SSL certificate and key files are generated in the data directory.

- validate_password is installed and enabled.

- A superuser account 'root'@'localhost is created. A password for the superuser is set and stored in the error log file. To reveal it, use the following command:
```
shell> sudo grep 'temporary password' /var/log/mysqld.log
```
Change the root password as soon as possible by logging in with the generated, temporary password and set a custom password for the superuser account:
```
shell> mysql -uroot -p
```
```
mysql> ALTER USER 'root'@'localhost' IDENTIFIED BY 'MyNewPass4!';
```

> **Note**
validate_password is installed by default. The default password policy implemented by validate_password requires that passwords contain at least one uppercase letter, one lowercase letter, one digit, and one special character, and that the total password length is at least 8 characters.



> **Note**
Compatibility Information for EL7-based platforms: The following RPM packages from the native software repositories of the platforms are incompatible with the package from the MySQL Yum repository that installs the MySQL server. Once you have installed MySQL using the MySQL Yum repository, you cannot install these packages (and vice versa).
- akonadi-mysql


## Testing the Server
After the data directory is initialized and you have started the server, perform some simple tests to make sure that it works satisfactorily. This section assumes that your current location is the MySQL installation directory and that it has a bin subdirectory containing the MySQL programs used here. If that is not true, adjust the command path names accordingly.

Alternatively, add the bin directory to your PATH environment variable setting. That enables your shell (command interpreter) to find MySQL programs properly, so that you can run a program by typing only its name, not its path name. See Section 4.2.9, “Setting Environment Variables”.

Use mysqladmin to verify that the server is running. The following commands provide simple tests to check whether the server is up and responding to connections:
```
shell> bin/mysqladmin version
shell> bin/mysqladmin variables
```
If you cannot connect to the server, specify a -u root option to connect as root. If you have assigned a password for the root account already, you'll also need to specify -p on the command line and enter the password when prompted. For example:
```
shell> bin/mysqladmin -u root -p version
Enter password: (enter root password here)
```
The output from mysqladmin version varies slightly depending on your platform and version of MySQL, but should be similar to that shown here:
```
shell> bin/mysqladmin version
mysqladmin  Ver 14.12 Distrib 8.0.24, for pc-linux-gnu on i686
...

Server version          8.0.24
Protocol version        10
Connection              Localhost via UNIX socket
UNIX socket             /var/lib/mysql/mysql.sock
Uptime:                 14 days 5 hours 5 min 21 sec

Threads: 1  Questions: 366  Slow queries: 0
Opens: 0  Flush tables: 1  Open tables: 19
Queries per second avg: 0.000
```
To see what else you can do with mysqladmin, invoke it with the --help option.

Verify that you can shut down the server (include a -p option if the root account has a password already):
```
shell> bin/mysqladmin -u root shutdown
```
Verify that you can start the server again. Do this by using mysqld_safe or by invoking mysqld directly. For example:
```
shell> bin/mysqld_safe --user=mysql &
```
If mysqld_safe fails, see Section 2.10.2.1, “Troubleshooting Problems Starting the MySQL Server”.

Run some simple tests to verify that you can retrieve information from the server. The output should be similar to that shown here.

Use mysqlshow to see what databases exist:
```
shell> bin/mysqlshow
+--------------------+
|     Databases      |
+--------------------+
| information_schema |
| mysql              |
| performance_schema |
| sys                |
+--------------------+
```
The list of installed databases may vary, but always includes at least mysql and information_schema.

If you specify a database name, mysqlshow displays a list of the tables within the database:
```
shell> bin/mysqlshow mysql
Database: mysql
+---------------------------+
|          Tables           |
+---------------------------+
| columns_priv              |
| component                 |
| db                        |
| default_roles             |
| engine_cost               |
| func                      |
| general_log               |
| global_grants             |
| gtid_executed             |
| help_category             |
| help_keyword              |
| help_relation             |
| help_topic                |
| innodb_index_stats        |
| innodb_table_stats        |
| ndb_binlog_index          |
| password_history          |
| plugin                    |
| procs_priv                |
| proxies_priv              |
| role_edges                |
| server_cost               |
| servers                   |
| slave_master_info         |
| slave_relay_log_info      |
| slave_worker_info         |
| slow_log                  |
| tables_priv               |
| time_zone                 |
| time_zone_leap_second     |
| time_zone_name            |
| time_zone_transition      |
| time_zone_transition_type |
| user                      |
+---------------------------+
```
Use the mysql program to select information from a table in the mysql schema:
```
shell> bin/mysql -e "SELECT User, Host, plugin FROM mysql.user" mysql
+------+-----------+-----------------------+
| User | Host      | plugin                |
+------+-----------+-----------------------+
| root | localhost | caching_sha2_password |
+------+-----------+-----------------------+
```
At this point, your server is running and you can access it.



## Securing the Initial MySQL Account
The MySQL installation process involves initializing the data directory, including the grant tables in the mysql system schema that define MySQL accounts. For details, see Section 2.10.1, “Initializing the Data Directory”.

This section describes how to assign a password to the initial root account created during the MySQL installation procedure, if you have not already done so.

> Note
Alternative means for performing the process described in this section:
- On Windows, you can perform the process during installation with MySQL Installer (see Section 2.3.3, “MySQL Installer for Windows”).
- On all platforms, the MySQL distribution includes mysql_secure_installation, a command-line utility that automates much of the process of securing a MySQL installation.
- On all platforms, MySQL Workbench is available and offers the ability to manage user accounts (see Chapter 31, MySQL Workbench ).

A password may already be assigned to the initial account under these circumstances:

- On Windows, installations performed using MySQL Installer give you the option of assigning a password.

- Installation using the macOS installer generates an initial random password, which the installer displays to the user in a dialog box.

- Installation using RPM packages generates an initial random password, which is written to the server error log.

- Installations using Debian packages give you the option of assigning a password.

- For data directory initialization performed manually using mysqld --initialize, mysqld generates an initial random password, marks it expired, and writes it to the server error log. See Section 2.10.1, “Initializing the Data Directory”.

The mysql.user grant table defines the initial MySQL user account and its access privileges. Installation of MySQL creates only a 'root'@'localhost' superuser account that has all privileges and can do anything. If the root account has an empty password, your MySQL installation is unprotected: Anyone can connect to the MySQL server as root without a password and be granted all privileges.

The 'root'@'localhost' account also has a row in the mysql.proxies_priv table that enables granting the PROXY privilege for ''@'', that is, for all users and all hosts. This enables root to set up proxy users, as well as to delegate to other accounts the authority to set up proxy users. See Section 6.2.18, “Proxy Users”.

To assign a password for the initial MySQL root account, use the following procedure. Replace root-password in the examples with the password that you want to use.

Start the server if it is not running. For instructions, see Section 2.10.2, “Starting the Server”.

The initial root account may or may not have a password. Choose whichever of the following procedures applies:

- If the root account exists with an initial random password that has been expired, connect to the server as root using that password, then choose a new password. This is the case if the data directory was initialized using mysqld --initialize, either manually or using an installer that does not give you the option of specifying a password during the install operation. Because the password exists, you must use it to connect to the server. But because the password is expired, you cannot use the account for any purpose other than to choose a new password, until you do choose one.

If you do not know the initial random password, look in the server error log.

Connect to the server as root using the password:
```
shell> mysql -u root -p
Enter password: (enter the random root password here)
```
Choose a new password to replace the random password:
```
mysql> ALTER USER 'root'@'localhost' IDENTIFIED BY 'root-password';
```
If the root account exists but has no password, connect to the server as root using no password, then assign a password. This is the case if you initialized the data directory using mysqld --initialize-insecure.

Connect to the server as root using no password:
```
shell> mysql -u root --skip-password
```
Assign a password:
```
mysql> ALTER USER 'root'@'localhost' IDENTIFIED BY 'root-password';
```
After assigning the root account a password, you must supply that password whenever you connect to the server using the account. For example, to connect to the server using the mysql client, use this command:
```
shell> mysql -u root -p
Enter password: (enter root password here)
```
To shut down the server with mysqladmin, use this command:
```
shell> mysqladmin -u root -p shutdown
Enter password: (enter root password here)
```
> Note
For additional information about setting passwords, see Section 6.2.14, “Assigning Account Passwords”. If you forget your root password after setting it, see Section B.3.3.2, “How to Reset the Root Password”.
To set up additional accounts, see Section 6.2.8, “Adding Accounts, Assigning Privileges, and Dropping Accounts”.