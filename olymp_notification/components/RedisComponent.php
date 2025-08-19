<?php
namespace app\components;
use Yii;

class RedisComponent
{
    public static function set($key, $data){
        Yii::$app->redis->set($key, json_encode($data));
    }
    public static function delete($key)
    {
        Yii::$app->redis->del($key);
    }
    public static function get($key){
        return json_decode(Yii::$app->redis->get($key));
    }
    public static function add($key, $newData)
    {
        $currentData = self::get($key);
        if ($currentData === null) {
            $currentData = [];
        }
        if (is_array($currentData)) {
            if (is_array($newData)) {
                $currentData = array_merge($currentData, $newData);
            } else {
                $currentData[] = $newData;
            }
        }
        elseif (is_numeric($currentData) && is_numeric($newData)) {
            $currentData += $newData;
        }
        self::set($key, $currentData);
        return $currentData;
    }
}