<?php

namespace app\services;

use Yii;

class SmsService
{
    private ApiService $apiService;
    public function __construct(
        ApiService $apiService
    )
    {
        $this->apiService = $apiService;
    }

    public function sendSMS($phone, $code){

        $smsAeroMessage = new \SmsAero\SmsAeroMessage('assbanlucky@gmail.com', 'nPFlpnCJATeiaSE84Rc-YOiAcBp9-j43');
        // Отправка SMS сообщений
        $response = $smsAeroMessage->send(['number' => $phone, 'text' => 'ВСОШ. Ваш код: ' . $code, 'sign' => 'SMS Aero']);
        return $response;
        

        // $this->apiService->post(
        //     'https://api.exolve.ru/messaging/v1/SendSMS',
        //     [
        //         'number' => Yii::$app->params['sms_phone_from'],
        //         'text' => 'ВСОШ. Ваш код: ' . $code ,
        //         'destination' => $phone
        //     ],
        //     [
        //         'Authorization: Bearer '. Yii::$app->params['sms_api_token'],
        //     ]
        // );

    }
}