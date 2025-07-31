<?php

namespace App\Services;

use App\Components\RabbitMQComponent;

class RabbitMQService
{
    private RabbitMQComponent $component;
    public function __construct(
        RabbitMQComponent $component
    )
    {
        $this->component = $component;
    }
    public function publish($queueNames, $appName,  $method, $table, $attributes, $searchAttributes = NULL){
        foreach ($queueNames as $queueName) {
            $data = json_encode([
                'appName' => $appName,
                'method' => $method,
                'data' => [
                    'table' => $table,
                    'attributes' => $attributes,
                    'searchAttributes' => $searchAttributes,
                ]
            ]);
            $this->component->publish($queueName, $data);
        }
    }
    public function consume($queueName){
        $data = [];
        $this->component->consume($queueName, function ($message) use (&$data) {
            $data[] = json_decode($message);
            return $message;
        }, true);
        return $data;
    }
}
