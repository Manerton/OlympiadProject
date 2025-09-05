<?php

namespace app\services;

use app\models\Notification;
use app\repositories\NotificationRepository;

class NotificationService
{
    private NotificationRepository $notificationRepository;
    public function __construct(
        NotificationRepository $notificationRepository
    )
    {
        $this->notificationRepository = $notificationRepository;
    }

    public function create($userId, $message, $status)
    {
        $notification = new Notification();
        $this->notificationRepository->create($notification, $userId, $message, $status);
    }
}