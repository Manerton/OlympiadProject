<?php

namespace app\components;

use yii\base\Component;

class MessageDictionary extends BaseDictionary
{
    const CODE_MESSAGE = 1;
    const TEXT_MESSAGE = 2;
    const ONLINE_NOTIFICATION = 3;

    public function __construct()
    {
        parent::__construct();
        $this->list = [
            self::CODE_MESSAGE => 'Сообщение с кодом',
            self::TEXT_MESSAGE => 'Текстовое сообщение',
            self::ONLINE_NOTIFICATION => 'Онлайн-уведомление'
        ];
    }

    public function customSort()
    {
        return [
            $this->list[self::CODE_MESSAGE],
            $this->list[self::TEXT_MESSAGE],
            $this->list[self::ONLINE_NOTIFICATION]
        ];
    }
}