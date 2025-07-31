<?php

namespace App\Services;

use App\Builder\EventBuilder;
use App\Builder\EventJuryBuilder;
use App\Builder\UserBuilder;
use App\Repositories\EventJuryRepository;
use App\Repositories\EventRepository;
use App\Repositories\UserRepository;

class EventJuryService
{
    private EventJuryRepository $eventJuryRepository;

    private EventJuryBuilder $eventJuryBuilder;
    private UserRepository $userRepository;
    private UserBuilder $userBuilder;
    private EventRepository $eventRepository;
    private EventBuilder $eventBuilder;
    public function __construct(
        EventJuryRepository $eventJuryRepository,
        EventJuryBuilder $eventJuryBuilder,
        UserRepository $userRepository,
        UserBuilder $userBuilder,
        EventRepository $eventRepository,
        EventBuilder $eventBuilder
    ){
        $this->eventJuryRepository = $eventJuryRepository;
        $this->eventJuryBuilder = $eventJuryBuilder;
        $this->userRepository = $userRepository;
        $this->userBuilder = $userBuilder;
        $this->eventRepository = $eventRepository;
        $this->eventBuilder = $eventBuilder;
    }
    public function findByEventId($eventId)
    {
        $juries = [];
        $data = $this->eventJuryRepository->getByEventId($eventId);
        foreach ($data as $item) {
            $eventJury = $this->eventJuryBuilder->build($item);
            $event = $this->eventBuilder->build($this->eventRepository->getByApiId($eventJury->event_id));
            $user = $this->userBuilder->build($this->userRepository->getByApiId($eventJury->user_id));
            $this->eventJuryBuilder->buildEvent($eventJury, $event);
            $this->eventJuryBuilder->buildUser($eventJury, $user);
            $juries[] = $eventJury;
        }
        return $juries;
    }
}
