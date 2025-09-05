<?php

namespace app\repositories;

use app\models\Notification;

class NotificationRepository
{

    public function get($id)
    {
        return Notification::findOne($id);
    }
    public function getAll()
    {
        return Notification::find()->all();
    }
    public function getByUserId($userId)
    {
        return Notification::find()->where(['user_id' => $userId])->all();
    }
    public function save(Notification $notification){
        return $notification->save();
    }
    public function delete(Notification $notification){
        return $notification->delete();
    }
    public function create(Notification $notification, $userId, $message, $status){
        $notification->user_id = $userId;
        $notification->message = $message;
        $notification->status = $status;
        return $this->save($notification);
    }
}