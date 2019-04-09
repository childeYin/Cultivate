# PHP

1. 什么是php-fpm， 以及重启php-fpm 是如何处理得。
2. PHP-FPM，PHP CLI ，FastCGI
3. 命名空间，__autoload，spl_autoload_register, register_showdown_function,pcntl_fork
4. 读取文件,
    - 按行读取
5. php底层机制
    - Application
    - SAPI
    - PHP, PHP API, Extensions(mysql, standard library etc)
    - Zend API, Zend Extension API
    - Zend Engine
6. php垃圾回收机制
    - refcount
    - is_ref
7. php如何调试代码
8. www.conf(request_slowlog_timeout, slowlog, listen.backlog，max_execution_time，request_terminate_timeout)
9. php.ini(error_log，log_errors )
10. 如果nginx的超时时间是10s，但是php脚本执行需要大于10s，那么log里分别有怎样的记录？
    1. 页面报错是什么？
    2. nginx error_log， php error_log, php slowlog？
    3. nginx进程终止后，php的进程是否继续进行？ 如果php继续执行和哪个参数有关（request_terminate_timeout）
11. php backlog 这个的作用是什么？以及如何调整？
12. null false 0 判断请用 ===  谨记
13. php 索引数组需要确认打印出来是不是数组，可以使用的方法如 打印json串


