<?php

namespace app\commands;

use app\components\RedisComponent;
use yii\console\Controller;

class RedisController extends Controller
{
    public function actionSet(){
        RedisComponent::set('email' , 'drive16052003@gmail.com');
    }
}