<?php

namespace App\Components\Dictionaries;

class DisabilityDictionary
{
    public const HEALTHY = 1;
    public const DISABILITY = 2;
    public static function getList(){
        return [
            self::HEALTHY => 'Нет ОВЗ',
            self::DISABILITY => 'Есть ОВЗ'
        ];
    }
}
