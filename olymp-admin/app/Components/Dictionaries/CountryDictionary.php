<?php

namespace App\Components\Dictionaries;

class CountryDictionary
{
    public const RUSSIA = 1;
    public static function getList(){
        return [
            self::RUSSIA => 'Россия'
        ];
    }
}
