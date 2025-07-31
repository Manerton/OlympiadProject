<?php

namespace App\Components;

use PhpAmqpLib\Connection\AMQPStreamConnection;
use PhpAmqpLib\Message\AMQPMessage;
class RabbitMQComponent
{
    protected $connection;
    protected $channel;

    public function __construct()
    {
        $config = config('queue.connections.rabbitmq');

        $this->connection = new AMQPStreamConnection(
            $config['host'],
            $config['port'],
            $config['user'],
            $config['password'],
            $config['vhost']
        );

        $this->channel = $this->connection->channel();
    }

    public function getConnection()
    {
        return $this->connection;
    }

    public function publish(string $queue, string $message, array $options = [])
    {
        $this->channel->queue_declare($queue, false, true, false, false);

        $msgOptions = array_merge([
            'delivery_mode' => AMQPMessage::DELIVERY_MODE_PERSISTENT
        ], $options);

        $msg = new AMQPMessage($message, $msgOptions);
        $this->channel->basic_publish($msg, '', $queue);
    }

    public function consume(string $queue, callable $callback, bool $processAll = false)
    {
        $this->channel->queue_declare($queue, false, true, false, false);

        if ($processAll) {
            $queueInfo = $this->channel->queue_declare($queue, false, true, false, false, false);
            $messageCount = $queueInfo[1];

            if ($messageCount > 0) {
                $this->channel->basic_consume(
                    $queue,
                    '',
                    false,
                    true,
                    false,
                    false,
                    function ($msg) use ($callback, &$messageCount) {
                        $callback($msg->body);
                        $messageCount--;
                        if ($messageCount <= 0) {
                            $msg->delivery_info['channel']->basic_cancel($msg->delivery_info['consumer_tag']);
                        }
                    }
                );

                while ($messageCount > 0 && $this->channel->is_consuming()) {
                    $this->channel->wait();
                }
            }
        } else {
            $this->channel->basic_consume(
                $queue,
                '',
                false,
                true,
                false,
                false,
                function ($msg) use ($callback) {
                    $callback($msg->body);
                }
            );

            while ($this->channel->is_consuming()) {
                $this->channel->wait();
            }
        }
    }

    public function __destruct()
    {
        if ($this->channel) {
            $this->channel->close();
        }

        if ($this->connection) {
            $this->connection->close();
        }
    }
}
