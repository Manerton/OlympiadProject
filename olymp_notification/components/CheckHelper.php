<?php

namespace app\components;

use Yii;

class CheckHelper
{
    public static function checkAccess($appToken)
    {
        return $appToken == Yii::$app->params['requestToken'];
    }
}