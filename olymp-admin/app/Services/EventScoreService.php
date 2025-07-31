<?php

namespace App\Services;

use App\Repositories\EventScoreRepository;

class EventScoreService
{
    private EventScoreRepository $eventScoreRepository;
    public function __construct(
        EventScoreRepository $eventScoreRepository
    )
    {
        $this->eventScoreRepository = $eventScoreRepository;
    }

    public function create($eventId){

        if ($this->eventScoreRepository->getByEventId($eventId)){
            $this->eventScoreRepository->create($eventId, 0 , 0);
        }
        return $this->eventScoreRepository->getByEventId($eventId);
    }
    public function fill(
        $eventId,
        $data
    ){
        $eventScore = $this->eventScoreRepository->getByEventId($eventId);
        $eventScore->setScore($data['winner_score'], $data['prize_score']);
        $this->eventScoreRepository->save($eventScore);
    }
}
