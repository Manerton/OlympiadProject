// test-websocket.php
<?php
require __DIR__ . '/vendor/autoload.php';

use WebSocket\Client;

try {
    echo "Connecting to ws://notification.olymp.local:8095/notify...\n";

    $client = new Client("ws://notification.olymp.local:8095/notify", [
        'timeout' => 10
    ]);

    echo "Connected! Sending message...\n";

    $client->send(json_encode([
        'event' => 'olymp_admin',
        'data' => 'test message from PHP'
    ]));

    echo "Message sent. Waiting for response...\n";

    // Попробуйте получить ответ
    $response = $client->receive();
    echo "Response: " . $response . "\n";

    $client->close();
    echo "Connection closed.\n";

} catch (Exception $e) {
    echo "Error: " . $e->getMessage() . "\n";
}