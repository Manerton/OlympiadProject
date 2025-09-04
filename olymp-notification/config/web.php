<?php

$params = require __DIR__ . '/params.php';
$db = require __DIR__ . '/db.php';

$config = [
    'id' => 'basic',
    'basePath' => dirname(__DIR__),
    'bootstrap' => ['log'],
    'aliases' => [
        '@bower' => '@vendor/bower-asset',
        '@npm'   => '@vendor/npm-asset',
    ],
    'components' => [
        'request' => [
            'parsers' => [
                'application/json' => 'yii\web\JsonParser'
            ],
            'cookieValidationKey' => 'pP-Q0KSwTo_kMItrFg4uOGuNM3Lnih7_',
        ],
        'cache' => [
            'class' => 'yii\caching\FileCache',
        ],
        'user' => [
            'identityClass' => 'app\models\User',
            'enableAutoLogin' => true,
        ],
        'errorHandler' => [
            'errorAction' => 'site/error',
        ],
        'websocket' => [
            'class' => \app\services\WebSocketService::class
        ],
        'mailService' => [
            'class' => 'app\services\MailService',
            'fromEmail' => 'rkuzurgaliev@schooltech.ru', // этот адрес должен совпадать с SMTP-логином
        ],
        'redis' => [
            'class' => 'yii\redis\Connection',
            'hostname' => 'redis',
            'port' => 6379,
            'database' => 0,
        ],
        'mailer' => [
            'class' => \app\services\MailSymfonyService::class,
            'dsn' => 'smtp://rkuzurgaliev%40schooltech.ru:Pass123$45@smtp.yandex.ru:465?encryption=ssl',
        ],

        'apiService' => [
            'class' => 'app\services\ApiService',
            'baseUrl' => 'https://172.16.0.94/olymp_notification/web/',
            'timeout' => 60, // опционально
        ],
        'log' => [
            'traceLevel' => YII_DEBUG ? 3 : 0,
            'targets' => [
                [
                    'class' => 'yii\log\FileTarget',
                    'levels' => ['error', 'warning'],
                ],
            ],
        ],
        'db' => $db,
        'queue' => [
            'class' => \yii\queue\amqp_interop\Queue::class,
            'as log' => \yii\queue\LogBehavior::class,
            'driver' => 'enqueue/amqp-lib',
            'dsn' => 'amqp://notification-service:notification-password@rabbitmq:5672/%2f',
            'queueName' => 'olymp_notification',
        ],
        /*
        'urlManager' => [
            'enablePrettyUrl' => true,
            'showScriptName' => false,
            'rules' => [
            ],
        ],
        */
    ],
    'params' => $params,
];

if (YII_ENV_DEV) {
    // configuration adjustments for 'dev' environment
    $config['bootstrap'][] = 'debug';
    $config['modules']['debug'] = [
        'class' => 'yii\debug\Module',
        // uncomment the following to add your IP if you are not connecting from localhost.
        'allowedIPs' => ['*'],
    ];

    $config['bootstrap'][] = 'gii';
    $config['modules']['gii'] = [
        'class' => 'yii\gii\Module',
        // uncomment the following to add your IP if you are not connecting from localhost.
        //'allowedIPs' => ['127.0.0.1', '::1'],
    ];
}

return $config;
