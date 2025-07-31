<?php

namespace App\Components\Dictionaries;

class GenderDictionary
{
    public const MALE = 1;
    public const FEMALE = 2;
    public static function getList(){
        return [
            self::MALE => 'Мужской',
            self::FEMALE => 'Женский'
        ];
    }
}
