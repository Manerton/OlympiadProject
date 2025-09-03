<?php

namespace app\controllers;

use Yii;
use yii\web\Controller;

class NotificationController extends Controller
{
    public function actionSend()
    {
        if (Yii::$app->request->isPost) {
            $message = Yii::$app->request->post('message');
            Yii::$app->websocket->send($message);
        }
        return Yii::$app->response->data = json_encode([
            'status' => 200
        ]);
    }
    public function actionSendTo($id)
    {
        $message = 'Hello';
        Yii::$app->websocket->sendTo($message, $id);
        return Yii::$app->response->data = json_encode([
            'status' => 200
        ]);
    }
}