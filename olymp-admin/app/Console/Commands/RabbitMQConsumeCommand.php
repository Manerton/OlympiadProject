<?php

namespace App\Console\Commands;

use App\Components\RabbitMQComponent;
use App\Components\RabbitMQHelper;
use Illuminate\Console\Command;

class RabbitMQConsumeCommand extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'rabbitmq-consume';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Test data consume';
    private RabbitMQComponent $rabbitMQComponent;

    /**
     * Create a new command instance.
     *
     * @return void
     */
    public function __construct(
        RabbitMQComponent $rabbitMQComponent
    )
    {
        $this->rabbitMQComponent = $rabbitMQComponent;
        parent::__construct();
    }

    /**
     * Execute the console command.
     *
     * @return int
     */
    public function handle()
    {
        $data = [];
        $this->rabbitMQComponent->consume(RabbitMQHelper::ADMIN_QUEUE_NAME, function ($message) use (&$data) {
            $data[] = json_decode($message);
            return $message;
        });
        return 0;
    }
}
