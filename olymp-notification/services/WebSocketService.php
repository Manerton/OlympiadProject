<?php

namespace app\services;

use WebSocket\Client;

class WebSocketService
{
    public function send($data)
    {
        $client = new Client("ws://127.0.0.1:8090/notify");
        $client->send(json_encode([
            'event' => 'entity_created',
            'data' => [
                'data' => $data
            ]
        ]));
    }
}