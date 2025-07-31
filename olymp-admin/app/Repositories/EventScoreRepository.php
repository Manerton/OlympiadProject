<?php

namespace App\Repositories;

use App\Models\EventScore;

class EventScoreRepository
{
    public function get($id)
    {
        return EventScore::find($id);
    }
    public function getByEventId($eventId){
        return EventScore::where(['event_id' => $eventId])->first();
    }
    public function create($eventId, $prizeScore, $winnerScore)
    {
        $eventScore = new EventScore();
        $eventScore->event_id = $eventId;
        $eventScore->prize_score = $prizeScore;
        $eventScore->winner_score = $winnerScore;
        return $this->save($eventScore);
    }
    public function save(EventScore $eventScore){
        $eventScore->save();
        return $eventScore->id;
    }
    public function delete(EventScore $eventScore){
        return $eventScore->delete();
    }

}
