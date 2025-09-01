<?php

namespace app\services;

use WebSocket\Client;

class WebSocketService
{
    public function send($data)
    {
        try {
            $client = new Client("ws://olymp_websocket:8095/notify");
            $client->send(json_encode([
                'event' => 'entity_created',
                'data' => $data
            ]));
        }
        catch (\Exception $e) {
            var_dump($e->getMessage());
        }
    }
}