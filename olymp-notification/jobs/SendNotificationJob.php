<?php

namespace app\jobs;

use Yii;
use yii\base\BaseObject;

class SendNotificationJob extends BaseObject implements \yii\queue\JobInterface
{
    private $message;
    public function __construct(
        $message,
        $config = []
    )
    {
        $this->message = $message;
        parent::__construct($config);
    }

    public function execute($queue){
        Yii::$app->websocket->send($this->message);
    }
}