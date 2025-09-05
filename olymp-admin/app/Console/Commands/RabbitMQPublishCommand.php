<?php

namespace App\Console\Commands;

use App\Components\RabbitMQComponent;
use App\Components\RabbitMQHelper;
use App\Services\RabbitMQService;
use Illuminate\Console\Command;
use Yii;

class RabbitMQPublishCommand extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'rabbitmq-publish';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Test data publish';
    private RabbitMQService $rabbitMQService;

    /**
     * Create a new command instance.
     *
     * @return void
     */
    public function __construct(
        RabbitMQService $rabbitMQService
    )
    {
        $this->rabbitMQService = $rabbitMQService;
        parent::__construct();
    }

    /**
     * Execute the console command.
     *
     * @return int
     */
    public function handle()
    {
        $this->rabbitMQService->publish(
            [RabbitMQHelper::NOTIFICATION_QUEUE_NAME],
            RabbitMQHelper::ADMIN_QUEUE_NAME,
            RabbitMQHelper::CREATE,
            'notification',
            [
                'message' => 'RERERERE',
                'user_id' => 'frwfrewger',
                'status' => 12
            ]
        );
        return 0;
    }
}
