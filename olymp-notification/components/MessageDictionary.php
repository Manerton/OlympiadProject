<?php

namespace components;

use yii\base\Component;

class MessageDictionary extends BaseDictionary
{
    const CODE_MESSAGE = 1;
    const TEXT_MESSAGE = 2;

    public function __construct()
    {
        parent::__construct();
        $this->list = [
            self::CODE_MESSAGE => 'Сообщение с кодом',
            self::TEXT_MESSAGE => 'Текстовое сообщение',
        ];
    }

    public function customSort()
    {
        return [
            $this->list[self::CODE_MESSAGE],
            $this->list[self::TEXT_MESSAGE],
        ];
    }
}