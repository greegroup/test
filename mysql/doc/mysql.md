## Installing MySQL
通过以下命令安装MySQL（对于已启用dnf的系统，请使用替换`yum`命令`dnf`）：
```
shell> sudo yum install mysql-community-server
```
这将安装MySQL服务器的软件包（`mysql-community-server`）以及运行服务器所需组件的软件包，包括客户端的软件包（`mysql-community-client`），客户端和服务器的常见错误消息和字符集（`mysql-community-common`）以及共享的客户端库（`mysql-community-libs`）。 

## Starting the MySQL Server
使用以下命令启动MySQL服务器：
```
shell> systemctl start mysqld
```
您可以使用以下命令检查MySQL服务器的状态：
```
shell> systemctl status mysqld
```
如果操作系统已启用systemd, 标准 `systemctl` (or alternatively, service with the arguments reversed) 命令，如 `stop`, `start`, `status` 和 `restart`, 可用于管理MySQL服务器服务。mysqld服务默认情况下处于启用状态，并在系统重新启动时启动。

假设服务器的数据目录为空，则在服务器首次启动时，会发生以下情况：

- 服务器已初始化。

- SSL证书和密钥文件在数据目录中生成。

- validate_password已安装并启用。

- 创建root用户帐户 'root'@'localhost . 设置root用户的密码并将其存储在错误日志文件中。要显示它，请使用以下命令：
```
shell> sudo grep 'temporary password' /var/log/mysqld.log
```
通过使用生成的临时密码登录并尽快更改超级用户帐户的root密码，以更改root密码：
```
shell> mysql -uroot -p
```
```
mysql> ALTER USER 'root'@'localhost' IDENTIFIED BY 'MyNewPass4!';
```

> **Note** 
默认情况下会安装validate_password。validate_password实现的默认密码策略要求密码至少包含一个大写字母，一个小写字母，一位数字和一个特殊字符，并且密码总长度至少为8个字符。



> **Note**
基于EL7的平台的兼容性信息：平台的本机软件存储库中的以下RPM软件包与安装MySQL服务器的MySQL Yum存储库中的软件包不兼容。使用MySQL Yum存储库安装MySQL后，就无法安装这些软件包（反之亦然）。



## Testing the Server
初始化数据目录并启动服务器后，请执行一些简单的测试以确保其正常运行。本节假定您的当前位置是MySQL安装目录，并且它具有一个bin子目录，其中包含此处使用的MySQL程序。如果不正确，请相应地调整命令路径名。

或者，将bin目录添加到PATH环境变量设置中。这使您的外壳程序（命令解释器）能够正确找到MySQL程序，因此您可以通过仅键入其名称而不是其路径名来运行程序。

使用mysqladmin验证服务器是否正在运行。以下命令提供了简单的测试，以检查服务器是否已启动并响应连接：
```
shell> bin/mysqladmin version
shell> bin/mysqladmin variables
```
如果无法连接到服务器，请指定`-u root`选项以root身份连接。如果已经为根帐户分配了密码，则还需要在命令行上指定`-p`并在出现提示时输入密码。例如
```
shell> bin/mysqladmin -u root -p version
Enter password: (enter root password here)
```
mysqladmin版本的输出会根据您的平台和MySQL版本而略有不同，但应与此处显示的类似：
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
要查看您还可以使用mysqladmin做什么，请使用`--help`选项调用它。

验证是否可以关闭服务器（如果根帐户已经有密码，则包括`-p`选项）：
```
shell> bin/mysqladmin -u root shutdown
```
确认您可以再次启动服务器。通过使用mysqld_safe或直接调用mysqld来执行此操作。例如：
```
shell> bin/mysqld_safe --user=mysql &
```

运行一些简单的测试以验证您可以从服务器检索信息。输出应类似于此处所示。

使用`mysqlshow`查看存在的数据库：
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
安装的数据库列表可能会有所不同，但始终至少包括`mysql`和`information_schema`。

如果指定数据库名称，则`mysqlshow`将显示数据库中表的列表：
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
使用mysql程序从mysql模式中的表中选择信息：
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
MySQL安装过程涉及初始化数据目录，包括定义MySQL帐户的mysql系统架构中的授权表。

本节介绍了如何为在MySQL安装过程中创建的初始根帐户分配密码（如果尚未这样做的话）。

> Note
Alternative means for performing the process described in this section:
- 在Windows上，可以使用MySQL Installer在安装过程中执行该过程
- 在所有平台上，MySQL发行版均包含mysql_secure_installation，这是一个命令行实用程序，可自动执行许多保护MySQL安装的过程。
- 在所有平台上，都可以使用MySQL Workbench，并且可以管理用户帐户

在以下情况下，可能已经为初始帐户分配了密码：

- 在Windows上，使用MySQL Installer执行的安装使您可以选择分配密码。

- 使用macOS安装程序进行安装会生成一个初始随机密码，该密码会在对话框中显示给用户。

- 使用RPM软件包进行安装会生成一个初始随机密码，该密码将写入服务器错误日志中。

- 使用Debian软件包进行安装时，可以选择分配密码。

- 对于使用`mysqld --initialize`手动执行的数据目录初始化，mysqld会生成一个初始随机密码，将其标记为过期，然后将其写入服务器错误日志。

mysql.user授予表定义了初始的MySQL用户帐户及其访问权限。安装MySQL只会创建一个`'root'@'localhost'`超级用户帐户，该帐户具有所有特权并且可以执行任何操作。如果root帐户的密码为空，则您的MySQL安装不受保护：任何人都可以以root用户身份连接到MySQL服务器，而无需输入密码，并被授予所有特权。

'root'@'localhost'帐户在mysql.proxies_priv表中还具有一行，可为“ @”（即所有用户和所有主机）授予PROXY特权。这使root用户可以设置代理用户，以及将设置代理用户的权限委托给其他帐户。

要为初始MySQL root帐户分配密码，请使用以下过程。将示例中的root-password替换为您要使用的密码。

如果服务器未运行，请启动它。

初始的root帐户可能没有密码。选择以下任一过程：

- 如果root用户帐户的初始随机密码已过期，请使用该密码以root用户身份连接到服务器，然后选择一个新密码。如果数据目录是使用mysqld --initialize手动或使用安装程序初始化的，则该目录没有在安装操作期间为您提供指定密码的选项，则是这种情况。因为密码存在，所以必须使用它来连接服务器。但是由于密码已过期，因此除非选择一个新密码，否则您不能将帐户用于选择其他目的，而只能选择一个新密码。

    a. If you do not know the initial random password, look in the server error log.

    b. Connect to the server as root using the password:
    ```
    shell> mysql -u root -p
    Enter password: (enter the random root password here)
    ```
    c. Choose a new password to replace the random password:
    ```
    mysql> ALTER USER 'root'@'localhost' IDENTIFIED BY 'root-password';
    ```
- 如果root帐户存在但没有密码，请使用root用户不使用密码连接到服务器，然后分配密码。如果使用`mysqld --initialize-insecure`初始化了数据目录，就是这种情况。

    a. 不使用密码以root用户身份连接到服务器：
    ```
    shell> mysql -u root --skip-password
    ```
    b. Assign a password:
    ```
    mysql> ALTER USER 'root'@'localhost' IDENTIFIED BY 'root-password';
    ```

   

在为根帐户分配密码后，每当使用该帐户连接到服务器时，都必须提供该密码。例如，要使用mysql客户端连接到服务器，请使用以下命令：
```
shell> mysql -u root -p
Enter password: (enter root password here)
```

要使用mysqladmin关闭服务器，请使用以下命令：
```
shell> mysqladmin -u root -p shutdown
Enter password: (enter root password here)
```
