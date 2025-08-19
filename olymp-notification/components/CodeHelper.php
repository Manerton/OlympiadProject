<?php

namespace app\components;

class CodeHelper
{
    /**
     * Возвращает случайный 4-значный код в виде строки (возможно с ведущими нулями)
     *
     * @return string
     */
    public static function generateCode(): string
    {
        return str_pad((string)random_int(0, 9999), 4, '0', STR_PAD_LEFT);
    }
}