#!/bin/sh
cd /var/www/olymp-admin
composer install --no-interaction --optimize-autoloader
php-fpm -F -R &
exec php artisan rabbitmq-handler
