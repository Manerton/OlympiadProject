<?php

namespace App\Console;

use App\Console\Commands\RabbitMQConsumeCommand;
use App\Console\Commands\RabbitMQHandlerCommand;
use App\Console\Commands\RabbitMQPublishCommand;
use Illuminate\Console\Scheduling\Schedule;
use Illuminate\Foundation\Console\Kernel as ConsoleKernel;

class Kernel extends ConsoleKernel
{
    /**
     * Define the application's command schedule.
     *
     * @param  \Illuminate\Console\Scheduling\Schedule  $schedule
     * @return void
     */
    protected $commands = [
        RabbitMQConsumeCommand::class,
        RabbitMQPublishCommand::class,
        RabbitMQHandlerCommand::class,
    ];
    protected function schedule(Schedule $schedule)
    {
        // $schedule->command('inspire')->hourly();
    }

    /**
     * Register the commands for the application.
     *
     * @return void
     */
    protected function commands()
    {
        $this->load(__DIR__.'/Commands');

        require base_path('routes/console.php');
    }
}
