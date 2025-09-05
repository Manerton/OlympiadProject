<?php

namespace App\Console\Commands;

use App\Components\RabbitMQComponent;
use App\Components\RabbitMQHelper;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;

class RabbitMQHandlerCommand extends Command
{
    /**
     * The name and signature of the console command.
     *
     * @var string
     */
    protected $signature = 'rabbitmq-handler';

    /**
     * The console command description.
     *
     * @var string
     */
    protected $description = 'Command description';

    /**
     * Create a new command instance.
     *
     * @return void
     */
    private RabbitMQComponent $rabbitMQComponent;
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
        $this->rabbitMQComponent->consume(RabbitMQHelper::ADMIN_QUEUE_NAME, function ($message) use (&$data) {
            $this->message([json_decode($message)]);
            return $message;
        });
        return 0;
    }
    public function message($data){
        foreach ($data as $item){
            switch ($item->method) {
                case RabbitMQHelper::CREATE:
                    $this->create($item);
                    break;
                case RabbitMQHelper::UPDATE:
                    $this->update($item);
                    break;
                case RabbitMQHelper::DELETE:
                    $this->delete($item);
                    break;
            }
        }
    }
    public function create($item)
    {
        DB::table($item->data->table)->insert((array)$item->data->attributes);
    }

    public function update($item)
    {
        DB::table($item->data->table)
            ->where((array)$item->data->searchAttributes)
            ->update((array)$item->data->attributes);
    }

    public function delete($item)
    {
        DB::table($item->data->table)
            ->where((array)$item->data->searchAttributes)
            ->delete();
    }
}
