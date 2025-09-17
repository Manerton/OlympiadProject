<?php

namespace app\controllers\api;

use app\jobs\SendNotificationJob;
use Yii;
use yii\web\Controller;

class NotificationController extends Controller
{
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
        }
        return Yii::$app->response->data = json_encode([
            'status' => 200
        ]);
    }
    public function actionSendTo($id)
    {
        if (Yii::$app->request->isPost) {
            $message = Yii::$app->request->post('message') ?: 'MESSAGE';
            Yii::$app->websocket->sendTo($message, $id);
        }
        return Yii::$app->response->data = json_encode([
            'status' => 200
        ]);
    }
}