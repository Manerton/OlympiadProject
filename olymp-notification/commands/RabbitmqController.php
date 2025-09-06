<?php

namespace app\commands;

use app\components\RabbitMQHelper;
use app\jobs\SendNotificationJob;
use Yii;
use yii\console\Controller;

class RabbitmqController extends Controller
{
    public function actionListen(){
        Yii::$app->rabbitmq->consume(RabbitMQHelper::NOTIFICATION_QUEUE_NAME, function ($message) use (&$data) {
            $this->message([json_decode($message)]);
            Yii::$app->queue->push(new SendNotificationJob($message, json_decode($message)['user_id']));
            return $message;
        });

        return 0;
    }
    public function message($data){
        foreach ($data as $item){
            switch ($item->method) {
                case RabbitMQHelper::CREATE:
                    $this->create($item);
                    break;
                case RabbitMQHelper::UPDATE:
                    $this->update($item);
                    break;
                case RabbitMQHelper::DELETE:
                    $this->delete($item);
                    break;
            }
        }
    }
    public function create($item){
        $command = Yii::$app->db->createCommand();
        $command->insert($item->data->table, (array)$item->data->attributes);
        $command->execute();
    }
    public function update($item){
        $command = Yii::$app->db->createCommand();
        $command->update($item->data->table, (array)$item->data->attributes, (array)$item->data->searchAttributes);
        $command->execute();
    }
    public function delete($item){
        $command = Yii::$app->db->createCommand();
        $command->delete($item->data->table, (array)$item->data->searchAttributes);
        $command->execute();
    }
}