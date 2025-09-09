<?php

namespace Database\Seeders;

use App\Models\EventScore;
use Illuminate\Database\Seeder;

class EventScoreSeeder extends Seeder
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
