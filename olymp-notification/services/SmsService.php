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


        $url = 'https://zvonok.com/manager/cabapi_external/api/v1/phones/tellcode/?' . http_build_query([
            'public_key'   => Yii::$app->params['sms_zvonok_public_key'],
            'phone'        => $phone,
            'campaign_id'  => Yii::$app->params['sms_zvonok_campaign_id'],
            'pincode'      => $code,
        ]);

        return $this->apiService->post($url);
        // return true;
                

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