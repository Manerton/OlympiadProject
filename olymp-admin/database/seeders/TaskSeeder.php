<?php

namespace Database\Seeders;

use App\Components\Dictionaries\TaskTypeDictionary;
use App\Models\Task;
use Illuminate\Database\Seeder;

class TaskSeeder extends Seeder
{
    private const EVENTS = [
        '55555555-5555-5555-5555-555555555571',
        '55555555-5555-5555-5555-555555555572',
        '55555555-5555-5555-5555-555555555573',
        '55555555-5555-5555-5555-555555555574'
    ];
    /**
     * Run the database seeds.
     *
     * @return void
     */
    public function run()
    {
        foreach(self::EVENTS as $eventId) {
            for($i = 1; $i <= 10; $i++) {
                $task = new Task();
                $task->event_id = $eventId;
                $task->number = $i;
                $task->type = rand(TaskTypeDictionary::READING, TaskTypeDictionary::TEST);
                $task->max_points = rand(0, 20);
                $task->save();
            }
        }
    }
}
