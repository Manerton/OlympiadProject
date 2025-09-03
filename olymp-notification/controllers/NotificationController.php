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

        return $this->render('index');
    }
    public function actionSendTo($id)
    {
        $message = 'Hello';
        Yii::$app->websocket->sendTo($message, $id);
        return $this->render('index');
    }
}