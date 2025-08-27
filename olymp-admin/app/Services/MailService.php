<?php

namespace App\Services;

use App\Components\ApiHelper;

class MailService
{
    private ApiService $apiService;
    public function __construct(
        ApiService $apiService
    )
    {
        $this->apiService = $apiService;
    }

    public function sendMessage($email, $message){
       $this->apiService->post(ApiHelper::SEND_MESSAGE_URL_API,
            [
                'email' => $email,
                'message' => $message,
                'requestToken' => ApiHelper::OLYMP_NOTIFICATION_TOKEN,
            ],
            [
                'Host' => 'notification.olymp.local',
            ]
       );
    }
}
