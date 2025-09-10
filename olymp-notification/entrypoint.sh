#!/bin/sh
cd /var/www/olymp-notification
composer install --no-interaction --optimize-autoloader
php-fpm -F -R &
until nc -z rabbitmq 5672; do
  echo "Жду RabbitMQ..."
  sleep 2
done
php yii queue/listen &
php yii rabbitmq/listen &
tail -f /dev/null