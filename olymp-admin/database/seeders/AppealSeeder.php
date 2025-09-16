<?php

namespace Database\Seeders;

use App\Components\Dictionaries\ReasonParticipantDictionary;
use App\Components\Dictionaries\StatusDictionary;
use App\Models\Appeal;
use App\Models\Task;
use Illuminate\Database\Seeder;

class AppealSeeder extends Seeder
{
    private const USERS = [
        '33333333-3333-3333-3333-333333333331',
        '33333333-3333-3333-3333-333333333332',
        '33333333-3333-3333-3333-333333333333',
        '33333333-3333-3333-3333-333333333334',
        '33333333-3333-3333-3333-333333333335',
        '33333333-3333-3333-3333-333333333336',
        '33333333-3333-3333-3333-333333333337',
        '33333333-3333-3333-3333-333333333338',
        '33333333-3333-3333-3333-333333333339',
        '33333333-3333-3333-3333-33333333333a',
    ];
    /**
     * Run the database seeds.
     *
     * @return void
     */
    public function run()
    {
        foreach(self::USERS as $user){
            foreach(Task::all() as $task){
                $model = new Appeal();
                $model->user_id = $user;
                $model->task_id = $task->id;
                $model->reason = ReasonParticipantDictionary::LAST_YEAR;
                $model->status = rand(StatusDictionary::AWAITING, StatusDictionary::REJECTED);
                $model->save();
            }
        }
    }
}
