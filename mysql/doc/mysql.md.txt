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

