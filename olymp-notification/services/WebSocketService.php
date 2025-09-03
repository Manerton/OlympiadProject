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
                'data' => $data,
                'to' => 'ALL'
            ]));
        }
        catch (\Exception $e) {
            var_dump($e->getMessage());
        }
    }
    public function sendTo($data, $id)
    {
        try {
            $client = new Client("ws://olymp_websocket:8095/notify");
            $client->send(json_encode([
                'event' => 'entity_created',
                'data' => $data,
                'to' => $id
            ]));
        }
        catch (\Exception $e) {
            var_dump($e->getMessage());
        }
    }
}