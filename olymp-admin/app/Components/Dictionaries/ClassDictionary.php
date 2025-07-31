<?php

namespace App\Components\Dictionaries;

class ClassDictionary
{
    public const CLASS_9 = 9;
    public const CLASS_10 = 10;
    public const CLASS_11 = 11;
    public static function getList(){
        return [
            self::CLASS_9 => '9 класс',
            self::CLASS_10 => '10 класс',
            self::CLASS_11 => '11 класс',
        ];
    }
}
