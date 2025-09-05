<?php

namespace app\controllers;

use app\components\MessageDictionary;
use app\jobs\SendNotificationJob;
use app\services\NotificationService;
use Yii;
use yii\web\Controller;

class NotificationController extends Controller
{
    private NotificationService $notificationService;
    public function __construct(
        $id,
        $module,
        NotificationService $notificationService,
        $config = []
    )
    {
        $this->notificationService = $notificationService;
        parent::__construct($id, $module, $config);
    }

    public function beforeAction($action)
    {
        if ($action->id === 'send') {
            $this->enableCsrfValidation = false;
        }
        return parent::beforeAction($action);
    }
    public function actionSend()
    {
        if (Yii::$app->request->isPost) {
            $message = Yii::$app->request->post('message') ?: 'MESSAGE';
            Yii::$app->queue->push(new SendNotificationJob(
                $message
            ));
            $this->notificationService->create('ALL', $message, MessageDictionary::ONLINE_NOTIFICATION);
        }
        return $this->render('index');
    }
    public function actionSendTo($id)
    {
        $message = Yii::$app->request->post('message') ?: 'MESSAGE';
        Yii::$app->queue->push(new SendNotificationJob(
            $message,
            $id
        ));
        $this->notificationService->create('ALL', $message, MessageDictionary::ONLINE_NOTIFICATION);
        return $this->render('index');
    }
}