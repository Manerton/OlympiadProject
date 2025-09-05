<?php

namespace app\jobs;

use Yii;
use yii\base\BaseObject;

class SendNotificationJob extends BaseObject implements \yii\queue\JobInterface
{
    private $message;
    private $id;
    public function __construct(
        $message,
        $id = null,
        $config = []
    )
    {
        $this->message = $message;
        $this->id = $id;
        parent::__construct($config);
    }

    public function execute($queue){
        if ($this->id){
            Yii::$app->websocket->sendTo($this->message, $this->id);
        }
        else {
            Yii::$app->websocket->send($this->message);
        }

    }
}