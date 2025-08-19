<?php

namespace app\components;

use Yii;
use yii\base\Component;
use yii\httpclient\Client;

class MessageComponent extends Component
{
    //рассылка по Telegram не работает, в связи с политикой приватности мессенджера
    public function sendTelegramMessage($chatId, $message)
    {
        $token = Yii::$app->params['telegramBotToken'];
        $url = "https://api.telegram.org/bot{$token}/sendMessage";

        $client = new Client();
        return $client->createRequest()
            ->setMethod('POST')
            ->setUrl($url)
            ->setData([
                'chat_id' => $chatId,
                'text' => $message,
            ])
            ->send();
    }
    //аналогично Telegram
    public function sendWhatsAppMessage($phoneNumber, $message)
    {
        $apiUrl = 'https://api.chat-api.com/instanceXXXX/message?token=XXXX';

        $client = new Client();
        return $client->createRequest()
            ->setMethod('POST')
            ->setUrl($apiUrl)
            ->setData([
                'phone' => $phoneNumber,
                'body' => $message,
            ])
            ->send();
    }
}