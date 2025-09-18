<?php

namespace App\Components\Dictionaries;

class EventStatusDictionary
{
    public const CONCLUDE_OFF = 1;
    public const CONCLUDE_ON = 2;
    public static function getList(){
        return [
            self::CONCLUDE_OFF => 'Результаты не опубликованы',
            self::CONCLUDE_ON => 'Результаты опубликованы'
        ];
    }
}
