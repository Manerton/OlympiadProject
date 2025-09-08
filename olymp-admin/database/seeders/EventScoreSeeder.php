<?php

namespace Database\Seeders;

use App\Models\EventScore;
use Illuminate\Database\Seeder;

class EventScoreSeeder extends Seeder
{
    private const EVENTS = ['1', '2', '3' ,'4' ,'5' ,'6' ,'7' ,'8' ,'9', '10', '11', '12'];
    /**
     * Run the database seeds.
     *
     * @return void
     */
    public function run()
    {
        foreach (self::EVENTS as $event) {
            $score = rand(1, 100);
            $eventScore = new EventScore();
            $eventScore->event_id = $event;
            $eventScore->prize_score = $score;
            $eventScore->winner_score = 2 * $score;
            $eventScore->save();
        }
    }
}
